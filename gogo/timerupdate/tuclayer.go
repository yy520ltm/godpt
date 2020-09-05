package timerupdate

import (
	"time"
)

//Timing updater component layer
//定时更新器组件层接口
type TUCLayer struct {
	TUs []*TimerUpdate //定时器集合

}

//创建定时器
func CreateExcuTu(timerName string, wt int, cyt time.Duration, Isstat bool, arg_f func(args ...interface{}) interface{}) *TimerUpdate {
	//创建定时器
	tu := &TimerUpdate{
		TimerName:     timerName,
		TimerWorkType: wt,
		STask:         make(chan *Task),
		CycleTime:     cyt,
		Arg_f:         arg_f,
		Isstatistics:  Isstat,
	}

	return tu
}

//查询定时更新器
func (tucl *TUCLayer) FindTu(timeName string) *TimerUpdate {
	for _, v := range tucl.TUs {
		if v.TimerName == timeName {
			return v
		}
	}
	return nil
}
