package main

import (
	/*
		#cgo CFLAGS: -I.
		#include <stdio.h>
		#include <stdlib.h>
	*/
	"C"
	"encoding/json"
	"fmgolib/core"
	"fmgolib/utils"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"unsafe"
)

//export Init
func Init(host *C.char, filePath *C.char) {
	c_host := C.GoString(host)
	c_filePath := C.GoString(filePath)
	// RocketMQ 日志级别
	rlog.SetLogLevel("error")
	// 初始化Logrs
	utils.LogInit("info")
	core.NewFmCore(c_host, c_filePath)
	return
}

//export GetValueByPointName
func GetValueByPointName(charData *C.char) *C.char {
	bytesList := C.GoString(charData)
	byte2, _ := json.Marshal(core.FmgolibGlobal.FuncGetPointValueByName(bytesList))
	cs := C.CString(string(byte2))
	defer C.free(unsafe.Pointer(cs))
	return cs
}

//export FuncTest
func FuncTest() {
	fmt.Println("hi frishme")
}

func main() {}
