package test_assert

import (
	//"container/ring"
	"fmt"
	"github.com/google/uuid"
	"sort"
	"testing"
)

func Test_Something(t *testing.T) {
	//var r = ring.New(2)

	var i int
	var cMap = map[string]string{}
	var keys []string
	for i < 1000 {

		// 找到数据将本次数据发送到Point
		if k, ok := cMap["VT_DB1_CRE_DAT"]; ok {
			//tag1:=fmt.Sprintf("%s1","VT_DB1_CRE_DAT")
			//tag2:=fmt.Sprintf("%s2","VT_DB1_CRE_DAT")
			UDI2, _ := uuid.NewRandom()
			keys = append(keys, UDI2.String())
			fmt.Println(UDI2.String())
			if k == UDI2.String() {
				// 两次一样
			} else {
				keys = append(keys, k)
				//fmt.Println("------------------------------")
			}
			sort.Strings(keys)
			//fmt.Println(keys)
			keys = []string{} // 重置
			delete(cMap, "VT_DB1_CRE_DAT")
		} else {
			// 找不到Key 将数据放入
			UDI1, _ := uuid.NewRandom()
			cMap["VT_DB1_CRE_DAT"] = UDI1.String()
		}

	}

	//fmt.Println(utils.Int64ToHex("1"))
	//fmt.Println(utils.Int64ToHex("0"))
	//fmt.Println(utils.Int64ToHex("100"))
	//fmt.Println( fmt.Sprintf("%06x",100))
	//fmt.Println(utils.IsLittleEndian())
	//s := "00000000"
	//fmt.Println(utils.SplitSubN(s,2)[0])
	//fmt.Println(utils.SplitSubN(s,2)[1])
	select {}
}

func Test2(t *testing.T) {
	windowDataCRP := []string{"1", "2", "3"}
	windowData := []string{"4", "5", "6"}
	windowDataCRP = append(windowDataCRP, windowData...)
	fmt.Println(windowDataCRP)
	data := []string{"336a165a-b558-4153-a227-6688824d2197", "336a165a-b558-4153-a227-6688824d2197", "f3b83f7f-ec1d-4cbc-98cf-5beee7f29dc4", "14248a02-9549-4358-bf0e-f269e6af073e"}
	for i, _ := range data {

		tag1 := fmt.Sprintf("%s%d", "VT_DB1_CRE_DAT", i)
		fmt.Println(tag1)
		fmt.Println(i)
	}

}
