package main

import (
	"context"
	"fmt"
	"os"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"1.1.1.1:9876"})),
		producer.WithRetry(2),
		producer.WithGroupName("pf_cg_app_control"),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: "1111",
			SecretKey: "2222",
		}),
	)

	if err != nil {
		fmt.Println("init producer error: " + err.Error())
		os.Exit(0)
	}

	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	for i := 0; i < 100; i++ {
		res, err := p.SendSync(context.Background(), primitive.NewMessage("a_to_p_control",
			[]byte("{\n    \"busyInfo\":{\n        \"createBy\":\"developer@rexel.com.cn\",\n        \"busyStatus\":\"todo\",\n        \"busyId\":\"20210605185951\",\n        \"createTime\":1622890792000,\n        \"personnelId\":\"13011112221\",\n        \"truckRfid\":\"truck_rfid000001\",\n        \"updateBy\":\"developer@rexel.com.cn\",\n        \"isDelete\":0,\n        \"tenantId\":\"71225760678100992002\",\n        \"busyType\":\"only_in\",\n        \"updateTime\":1622890792000,\n        \"enterTime\":1622890792000\n    },\n    \"truckInfo\":{\n        \"truckModel\":\"1\",\n        \"createBy\":\"admin\",\n        \"createTime\":1620957832000,\n        \"truckRfid\":\"truck_rfid000001\",\n        \"updateBy\":\"developer@rexel.com.cn\",\n        \"isDelete\":0,\n        \"truckDescribe\":\"川A-XXXA0\",\n        \"tenantId\":\"71225760678100992002\",\n        \"updateTime\":1622705488000,\n        \"plateNumber\":\"川A-XXXA0\"\n    },\n    \"personnel\":{\n        \"isDelete\":0,\n        \"sex\":\"male\",\n        \"fullName\":\"司机2\",\n        \"updateTime\":1622515483000,\n        \"idNumber\":\"\",\n        \"createBy\":\"admin\",\n        \"phoneNumber\":\"13011112221\",\n        \"createTime\":1620957832000,\n        \"personnelId\":\"13011112221\",\n        \"updateBy\":\"admin\",\n        \"station\":\"运输车司机\",\n        \"tenantId\":\"71225760678100992002\",\n        \"company\":\"langjiu\"\n    },\n    \"type\":\"enter\",\n    \"deviceInfo\":[\n        {\n            \"iotId\":\"cf5c6ddae45846368285afdb467864d8\",\n            \"pointId\":[\n                \"VT_DB2_LED1_TXT\"\n            ]\n        }\n    ],\n    \"ledContent\":\"请前往1#液体堆场1#行车1#车位执行入库任务,储藏位置:1190\",\n    \"busyStep\":[\n        {\n            \"isDelete\":0,\n            \"parkingId\":\"parking111\",\n            \"busyStepName\":\"入库\",\n            \"updateTime\":1622890792000,\n            \"stamnosRfid\":\"stamnos_rfid000001\",\n            \"busyStepId\":\"inbound\",\n            \"createBy\":\"developer@rexel.com.cn\",\n            \"busyId\":\"20210605185951\",\n            \"createTime\":1622890792000,\n            \"updateBy\":\"developer@rexel.com.cn\",\n            \"tenantId\":\"71225760678100992002\",\n            \"slotId\":\"1190\",\n            \"startTime\":1622890792000,\n            \"busyStepSort\":20\n        }\n    ]\n}")))

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}
