package dispatch

import (
	"fmt"
	"gogo/timerupdate"
	"time"
	"tools"
)

//中心调度接口
type ICenterDispatch interface {
	ListChannel()                 //监听管道
	SendUpdate(args func())       //发送通知给更新器//发送自定义策略规则给更新器
	ReceiveTask()                 //接收定时器那边传入的任务
	StarTicker()                  //启动定时器
	StopTicker(timerName string)  //stop 定时器
	ResetTicker(timerName string) //重置定时器
}

//调度器
type CenterDispatch struct {
	Name           string                   //调度器名称
	ReceiveChannel chan *timerupdate.Task   //调度器的接收管道 RC
	TU             *timerupdate.TimerUpdate //对应的更新器
	Isstatistics   bool                     //是否统计
}

func (cd *CenterDispatch) ListChannel() {

}
func (cd *CenterDispatch) SendUpdate(args func()) {

}
func (cd *CenterDispatch) ReceiveTask() {

	ty := &Tally{
		StartTime:       time.Now(),
		EndTime:         time.Time{},
		TallyNum:        0,
		LogTimeNum:      5,
		LogTimeDuration: time.Second,
		IsClear:         false,
	}

	for task := range cd.TU.STask {

		//开启统计功能
		if cd.Isstatistics {
			//开始统计某一段时间内发送的任务量
			ty.EndTime = time.Now()

			//获取当前的相差时间值
			beApart := tools.GetTimeSub(ty.StartTime, ty.EndTime)
			//具体单位数值
			apart := tools.GetTimeNum(ty.LogTimeDuration, beApart)

			//如果当前相隔时间小于 记录时间单位
			ty.IsClear = int(apart) <= ty.LogTimeNum
			//2 根据是否是否清空 结构来操作计数
			if ty.IsClear {
				//相当于 还在未来 (time.Duration(ty.LogTimeNum) * ty.LogTimeDuration) 的时间之前
				ty.TallyNum++
			} else {

				//相当于 还在未来5秒钟的时间之后 ，清空当前计数 并且发送通道，
				fmt.Println("-----------------------------------------------")
				fmt.Printf("在调度器:[ %v ]中,\n 从 %v到%v : \n %v秒 内一共发送了%v个任务<--->\n", cd.Name, ty.StartTime, ty.EndTime,ty.LogTimeNum, ty.TallyNum)
				fmt.Println("-----------------------------------------------")
				//NewGoWorkNum := SetGoWorkNum(ty.TallyNum)

				time.Sleep(5 * time.Second)
				//重置统计值
				ty.TallyNum = 0
				ty.IsClear = false
				//然后又把未来之后的时间给初始时间
				ty.StartTime = time.Now()

			}
		}


		cd.ReceiveChannel <- task
	}

}

func (cd *CenterDispatch) StarTicker() {

	cd.TU.StarTicker()
}
func (cd *CenterDispatch) StopTicker(timerName string) {
	cd.TU.Tt.Stop()

	fmt.Println(">>>>>>>>>>>>>>>>>>>>> Stop <<<<<<<<<<<<<<<<<<<<<")
	fmt.Println(timerName, "定时器以停止")

}
func (cd *CenterDispatch) ResetTicker(timerName string) {

}
