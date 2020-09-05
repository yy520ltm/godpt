package dispatch

import (
	"gogo/pools"
	"gogo/timerupdate"
)

//全局层
type GlobalLayer struct {
	CDL  *CDLayer              //中心调度层
	TUCL *timerupdate.TUCLayer //定时更新层

}

//初始全局各个层
func (gl *GlobalLayer) Inits(goworkerNum GoWorkerNum) {

	//中心调度层
	gl.CDL = &CDLayer{
		CDs:       []*CenterDispatch{},
		WorkerNum: goworkerNum,
	}

	//定时更新层
	gl.TUCL = &timerupdate.TUCLayer{
		TUs: []*timerupdate.TimerUpdate{},
	}

}

//根据更新器参数，创建对应的调度器-->更新器-->协程池
func (gl *GlobalLayer) BatchCreate(tus ...*timerupdate.TimerUpdate) {

	for _, v := range tus {
		//创建中心调度器
		centerD := gl.CDL.CreateCD(v.TimerName + "_CD",v.Isstatistics)
		gl.CDL.CDs = append(gl.CDL.CDs, centerD)

		//创建更新器层
		gl.TUCL.TUs = append(gl.TUCL.TUs, v)
		//把定时器附加到调度器中
		centerD.TU = v
	}
	//创建协程池
	gl.CDL.PL = pools.CreatePools("就我一个协程池", int(gl.CDL.WorkerNum))

}

//运行主程序调用
func GoRun(goworkerNum GoWorkerNum, tus ...*timerupdate.TimerUpdate) *GlobalLayer {
	gl := &GlobalLayer{}
	//初始化全局层中各个层
	gl.Inits(goworkerNum)
	//根据更新器信息，创建各个层 包括定时器层
	gl.BatchCreate(tus...)
	//开始构建各个层之间的管道对接
	for _, v := range gl.CDL.CDs {

		//开启定时器
		go v.StarTicker()
		//接收定时器中的任务
		go v.ReceiveTask()

	}
	//之所以要有上面两步把任务转换到通道中
	// 是因为后期要加入监听通道以及控制更新器操作
	//把RC通道里面的任务发送给协程池对外通道EC
	go gl.CDL.SendTaskToEC()
	// 当任务池处理了某些任务时，或者是 要反馈某些信号 时
	// 由 调度层 接收协程池的对外发送通道的数据  ETC 通道
	go gl.CDL.ReceiveETChan()

	gl.CDL.PL.Run()
	return gl
}
