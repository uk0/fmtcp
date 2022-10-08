package main

import (
	"fmgolib/core"
	"time"

	ext "github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"
)

func main() {
	fm := core.NewFmCore("127.0.0.1:5002", "/stream/config.json")

	StartTruckScalesWorker(fm)
	// 构建输入源：该输入源每隔一段时间 扫描一次
	source := ext.NewChanSource(tickerChanSource(time.Microsecond*500, fm, []string{"VT1"}))
	slidingWindow := flow.NewSlidingWindow(time.Second*2, time.Second*2)

	// 结果输出：将结果输出到标准输出
	sink := ext.NewStdoutSink()

	// 将流式处理各个步骤串联起来
	source.Via(slidingWindow).To(sink)

	select {}
}

func StartTruckScalesWorker(fc *core.FC) {
	/**
	测点类型 VT

	request
		地磅设备
		RFID CRE (车辆的RFID ，酒罐RFID)


	result
		DeviceList->find Contains LED-> send value to LED POINT
		返回值 设备列表 找到设备列表包含LED并且发送相关信息。
		LED
	*/

}

func tickerChanSource(repeat time.Duration, fc *core.FC, subPoint []string) chan interface{} {
	ticker := time.NewTicker(repeat)
	oc := ticker.C
	nc := make(chan interface{})
	go func() {
		for range oc {
			nc <- fc.FuncGetPointValueByListName(subPoint)
		}
	}()
	return nc
}
