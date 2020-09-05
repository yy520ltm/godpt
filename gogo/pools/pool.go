package pools

import (
	"context"
	"fmt"
	"gogo/timerupdate"
	"time"
)

type GoPools struct {
	PoolName             string                 //协程池名称
	GoWorkNum            int                    //协程工作的数量  通过调度层哪边根据任务情况 来安排协程池这边的工作数量
	CalculateChannel     chan int               //计算接收任务数的通道
	ExternalChannel      chan *timerupdate.Task //EC 对外【接收】通道
	JobsChannel          chan *timerupdate.Task //JC 内部工作通道
	ExternalTransChannel chan interface{}       //ETC 对外【发送】通道
	//map['Instruction']=stopticker Instruction 代表指令 具体做什么功能
	//map['Data']=result 返回的具体数据
}

//创建一个协程池
func CreatePools(poolName string, goNum int) *GoPools {
	gp := &GoPools{
		PoolName:             poolName,
		GoWorkNum:            goNum,
		CalculateChannel:     make(chan int),
		ExternalChannel:      make(chan *timerupdate.Task),
		JobsChannel:          make(chan *timerupdate.Task),
		ExternalTransChannel: make(chan interface{}),
	}

	return gp
}

// 独立工作单元
func (gp *GoPools) Worker(cancelCtx context.Context) {

	//对内工作通道
	for worker := range gp.JobsChannel {

		func (){
			cancelCtx1, cfu1 := context.WithTimeout(cancelCtx, 3*time.Second)
			defer func() { cfu1() }()

			result := worker.ExcuTask()
			//检测到工作任务完成之后不为nil 就协定为有向外部DW通道发送 信号给 中心调度器
			if result != nil {
				gp.ExternalTransChannel <- result
			}

			select {
			case <-cancelCtx1.Done():
				fmt.Println("time out!~!!我的guid是：", worker.Guid)

			default:


			}
		}()

	}

}

func (gp *GoPools) Run() {

	rootCtx := context.Background()
	cancelCtx, _ := context.WithCancel(rootCtx)

	//派发多个goroutine 来工作
	for i := 1; i <= gp.GoWorkNum; i++ {

		go gp.Worker(cancelCtx)
	}

	for t := range gp.ExternalChannel {

		gp.JobsChannel <- t
	}
}
