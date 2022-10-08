package main

import (
	"fmgolib/core"
	"fmgolib/utils"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	jsoniter "github.com/json-iterator/go"
	"time"
)

func main() {

	// RocketMQ 日志级别
	rlog.SetLogLevel("error")
	// 初始化Logrs
	utils.LogInit("error")

	time.Sleep(1 * time.Second)

	_ = core.NewFmCore("192.168.2.245:5002", "/Users/firshme/Desktop/IotPlatformLite/fmgolib/test_get_data/config.json")
	//
	//data := core.FmgolibGlobal.FuncGetPointValueByNameListLink([]string{"VT_2Q3KA_LED_VM", "VT_1Q1KA_LED_VM", "VT_1Q2KA_LED_VM", "VT_1Q1KA_LED_VM", "VT_1Q3KA_LED_VM"})
	//
	//data2 := core.FmgolibGlobal.FuncGetPointValueByListName([]string{"VT_2Q3KA_LED_VM", "VT_1Q1KA_LED_VM", "VT_1Q2KA_LED_VM", "VT_1Q1KA_LED_VM", "VT_1Q3KA_LED_VM"})
	//
	//utils.Debug("----------------------------------------------------")
	//utils.Debug(data)
	//
	//utils.Debug(data2)
	//core.FmgolibGlobal.FuncModifyValueByName([]core.ModifyStruct{core.ModifyStruct{
	//	PointName: "AR3",
	//	Value:     "100.121",
	//}})

	// 3e2a271d09000700005d03616263
	// 3e2a271d09000700005d03616263
	go func() {
		for i := 0; i < 100000000; i++ {

			var json = jsoniter.ConfigCompatibleWithStandardLibrary

			data := core.FmgolibGlobal.FuncGetPointValueByNameListLink([]string{"%RAND"})

			tt, _ := json.Marshal(data)
			//core.FmgolibGlobal.FuncModifyValueByName([]core.ModifyStruct{core.ModifyStruct{
			//	PointName: "VA_5",
			//	Value:     data[0].Msg[0].NewValue,
			//}})

			utils.Debug("----------------------------------------------------")
			utils.Error(string(tt))
			utils.Debug("----------------------------------------------------")

		}
	}()

	select {}

}
