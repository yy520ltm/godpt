package timerupdate

import (
	"fmt"
	"time"
)

//任务结构
type Task struct {
	Guid          string                                         //唯一标识符
	ExcuStartTime string                                         //执行时间
	ExcuendTime   string                                         //执行结束时间
	Level         int                                            //任务等级
	Func          func(args ...interface{}) (result interface{}) //具体任务函数

}


////执行任务
func (t *Task) ExcuTask() interface{} {
	t.ExcuStartTime = time.Now().String()
	result := t.Func(t.Guid)
	t.ExcuendTime = time.Now().String()


	fmt.Printf("开始时间执行时间:%s,结束时间: %s ,我的任务标志GUID是 %s \n", t.ExcuStartTime, t.ExcuendTime,t.Guid)
	return result

}
