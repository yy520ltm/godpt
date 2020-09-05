package main

import (
	"fmt"
	"gogo/dispatch"
	"gogo/timerupdate"
	"strconv"
	"time"
	"tools"
)

func main() {
	test()
}

//最简单的调用 开启一个定时器
func test() {
	args := func(args ...interface{}) interface{} {
		fmt.Println("我是测试函数,我的任务标志为:",args)
		return nil

	}

	//设置一个定时器
	tu1 := timerupdate.CreateExcuTu("testTimer", 1, time.Second, false,args)

	dispatch.GoRun(dispatch.GoWorkerNum(50), tu1)
}
//测试开启统计功能
func testIsStat() {
	args := func(args ...interface{}) interface{} {
		fmt.Println("我是测试函数,我的任务标志为:",args)
		return nil
	}

	//设置一个定时器  isstat true 开启统计
	tu1 := timerupdate.CreateExcuTu("testTimer", 1, 5*time.Nanosecond, true,args)

	dispatch.GoRun(dispatch.GoWorkerNum(50), tu1)
}
//测试 函数执行超时
func testTimeout()  {


	args := func(args ...interface{}) interface{} {
		//这里随机生成延迟时间来测试模拟超时情况
		num:=time.Duration(tools.RandInt64(1, 5, false))
		fmt.Println("测试ing,---",num,",guid:",args)
		time.Sleep(num*time.Second)
		return nil

	}

	//设置一个定时器
	tu1 := timerupdate.CreateExcuTu("testTimer", 1, time.Second,false, args)

	dispatch.GoRun(dispatch.GoWorkerNum(50), tu1)
}
//批量生成多个定时器
func test1() {

	urls := []string{
		"www.xxx1.com",
		"www.xxx2.com",
		"www.xxx3.com",
		"www.xxx4.com",
		"www.xxx5.com",
	}

	ts := []*timerupdate.TimerUpdate{}

	for ii, vv := range urls {

		tuName := "BookTimerUpdate" + strconv.Itoa(ii)
		link := vv
		args := func(args ...interface{}) interface{} {

			fmt.Println(link)

			return nil
		}

		cyt := time.Duration(
			tools.RandInt64(1, 3, true)) * time.Millisecond
		tu1 := timerupdate.CreateExcuTu(tuName, 1, cyt,false, args)
		ts = append(ts, tu1)
	}
	dispatch.GoRun(dispatch.GoWorkerNum(50), ts...)
}
//测试停用定时更新器
func testStopTicker()  {
	args := func(args ...interface{}) interface{} {

		num:=time.Duration(tools.RandInt64(1, 5, false))
		fmt.Println("我是几",num)
		if num>3 {
			//这里是当出现某种情况需要停止定时器，通知调度层做停止该定时器
			insMap := make(map[string]interface{})
			insMap["Instruction"] = "StopTicker"
			//这里告诉调度层，我要关闭那个调度器 ，在创建调度层时 就默认是更新器名称加上 _CD
			insMap["Data"] = "Books"+"_CD"

			return insMap
		}else {
			return nil
		}


	}

	cyt := time.Duration(
		tools.RandInt64(1, 5, false)) * time.Second
	tu1 := timerupdate.CreateExcuTu("Books", 1, cyt,false, args)

	cyt2 := time.Duration(
		tools.RandInt64(2, 7, false)) * time.Second
	tu2 := timerupdate.CreateExcuTu("Chapter", 1, cyt2,false, args)
	ts := []*timerupdate.TimerUpdate{tu1, tu2}

	dispatch.GoRun(dispatch.GoWorkerNum(500), ts...)
}
//自定义开启多个定时器
func test2() {

	args := func(args ...interface{}) interface{} {

		//这里是通知该函数完了，通知调度层做其他操作
		insMap := make(map[string]interface{})
		insMap["Instruction"] = "SuccessTask"
		insMap["Data"] = 1

		return insMap

	}

	args2 := func(args ...interface{}) interface{} {

		fmt.Println("鸡你太美")
		return nil

	}
	cyt := time.Duration(
		tools.RandInt64(1, 5, false)) * time.Millisecond
	tu1 := timerupdate.CreateExcuTu("Books", 1, cyt,false, args)

	cyt2 := time.Duration(
		tools.RandInt64(2, 7, false)) * time.Millisecond
	tu2 := timerupdate.CreateExcuTu("Chapter", 1, cyt2,false, args2)
	ts := []*timerupdate.TimerUpdate{tu1, tu2}

	dispatch.GoRun(dispatch.GoWorkerNum(500), ts...)
}
