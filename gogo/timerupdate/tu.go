package timerupdate

import (
	"fmt"
	"time"
	"tools"
)

//定时更新器
type TimerUpdate struct {
	TimerName     string                                //定时器名字
	TimerWorkType int                                   //定时器工作类型
	STask         chan *Task                            //发送任务通道
	CycleTime     time.Duration                         //定时器周期
	Tt            *time.Ticker                          //存放每一个定时更新器里面的time.Ticker
	Arg_f         func(args ...interface{}) interface{} //具体任务
	Isstatistics  bool                                  //是否统计

}

//创建一个Task
func (tu *TimerUpdate) CreateTask(arg_f func(args ...interface{}) interface{}) *Task {
	t := Task{
		Func:  arg_f,
		Level: 1,
		Guid: tools.GUID(),
	}
	return &t
}

//发送任务
func (tu *TimerUpdate) SendTask() {
	t := tu.CreateTask(tu.Arg_f)
	fmt.Println(tu.TimerName, t.Guid, "发送任务成功")
	//把任务发送给对应的更新器通道
	tu.STask <- t

}

//开启定时器
func (tu *TimerUpdate) StarTicker() {
	//开启一个定时器
	ticker := time.NewTicker(tu.CycleTime)
	tu.Tt = ticker

	for {
		for range ticker.C {
			fmt.Printf("------开启%s定时器--------\n", tu.TimerName)

			//挖掘并创建任务
			tu.SendTask()

		}
	}
}
