package core

import (
	"fmgolib/fmenum"
	"fmgolib/utils"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)
import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ResponseData struct {
	ReturnCode string `json:"rcode"`
	DataCode   string `json:"data_code"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	FuncCode   string `json:"function_code"`
	Msg        Data   `json:"data"`
	Version    string `json:"version"`
	IsModify   ModifyInfo
}

type Data []Point //自定义类型

type Point struct {
	Index    int64  `json:"index"`
	Name     string `json:"name"`
	NewValue string `json:"new_value"`
	Qty      int64  `json:"qty"`
}

type ModifyInfo struct {
	SuccessCount int64 `json:"succCount"`
	FailCount    int64 `json:"failCount"`
}

type MyPointBox struct {
	Items []Point
}

func (box *MyPointBox) AddPoint(item Point) []Point {
	box.Items = append(box.Items, item)
	return box.Items
}

// SingleRecordCover @title  解析单条的数据
// @description  获取Response 并且按照规则解析
// @auth      作者  张建新    时间（2021/6/1 22:57 ）
// @param     hexData        []string         "hex string 数组"
// @return    result         Point          "构建的多个Point对象"
func SingleRecordCover(hexData []string) Point {
	utils.Debug("hex recode Cover  " + fmt.Sprintln(hexData))
	data, _ := utils.GbkToUtf8(utils.Hex2Char(strings.Join(hexData[3:], "")))
	return Point{
		Index: utils.HexToDec(strings.Join(hexData[0:3], "")),
		Name:  string(data),
	}
}

// CoverStData @title    解析VT数据
// @description   解析VT数据并且构建Point
// @auth      作者  张建新    时间（2021/6/1 22:57 ）
// @param     hexData        string         "hex string 数组"
// @return    result         Point          "构建的多个Point对象"
func CoverStData(hexData []string, dataLength int) []Point {
	startIndex := 0
	var items []Point
	result := MyPointBox{items}
	cursor := 0
	lenData := 0
	for {
		jsonObj := Point{}
		if startIndex+3 >= len(hexData) {
			/// 弹出
			break
		}
		jsonObj.Index = utils.HexToDec(strings.Join(hexData[startIndex:startIndex+3], ""))
		if dataLength == 0 {
			letLenData := utils.HexToDec(strings.Join(hexData[startIndex+3:startIndex+3+1], ""))
			lenData = int(letLenData)
		} else {
			lenData = dataLength
		}
		cursor = lenData + cursor
		jsonObj.NewValue = string(utils.Hex2Char(strings.Join(hexData[startIndex+3+1:startIndex+3+1+lenData], "")))
		jsonObj.Qty = 0 // vt 为内部变量，不存在链接设备 默认0
		startIndex = startIndex + 3 + 1 + lenData
		result.AddPoint(jsonObj)
	}
	return result.Items
}

// CoverModifyPointValueByIndexAndType @title  修改数据
// @description  获取 Response 并且按照规则解析
// @auth      作者  张建新    时间（2021/6/1 22:57 ）
// @param     hexData        []string         "hex string 数组"
// @return    result         Point          "构建的多个Point对象"

func CoverModifyPointValueByIndexAndType(hexData []string) (result ModifyInfo) {
	result.SuccessCount = utils.HexToDec(strings.Join(hexData[6:8], ""))
	result.FailCount = utils.HexToDec(strings.Join(hexData[8:10], ""))
	return result
}

// CoverGetPointValueByIndexAndType @title  处理所有的Response
// @description  获取 Response 并且按照规则解析
// @auth      作者  张建新    时间（2021/6/1 22:57 ）
// @param     hexData        []string         "hex string 数组"
// @return    result         Point          "构建的多个Point对象"
func CoverGetPointValueByIndexAndType(hexData []string) []Point {
	doubles := []string{"01", "02", "03", "07"}
	flouts := []string{"12", "13", "17"}
	ints := []string{"05", "06", "08", "04"}
	str := []string{"09"}
	fCode := hexData[4]
	fData := hexData[6:]

	if utils.StrInSlice(fCode, doubles) {
		return RespCoverDataR4R8(fData, 8)
	}
	if utils.StrInSlice(fCode, flouts) {
		return RespCoverDataR4R8(fData, 4)
	}
	if utils.StrInSlice(fCode, ints) {
		return RespCoverDataR4R8(fData, 1)
	}
	if utils.StrInSlice(fCode, str) {
		return CoverStData(hexData[8:], 0)
	}
	return nil
}

func GetSingleRecord(aps []byte, lens int64, requestConn net.Conn) []byte {
	var temp []byte
	temp = append(temp, aps...)
	data := make([]byte, lens)
	readNum, err := requestConn.Read(data[:])
	if err == nil {
		temp = append(temp, data[:readNum]...)
		strHex := utils.ByteArrayToHex(temp)
		hexArray := utils.SplitSubN(strHex, 2)
		fmt.Println(CoverGetPointValueByIndexAndType(hexArray))
	}
	return nil
}

// CoverConnResponseToJSON @title 获取Response 的hex str 数据并且根据Type进行解析返回JSON
// @description  获取 Response 并且按照规则解析
// @auth      作者  张建新    时间（2021/6/1 22:57 ）
// @param     hexData        []string         "hex string 数组"
// @return    result         Point          "构建的多个Point对象"
func (f *FC) CoverConnResponseToJSON(responseCoverFunc fmenum.RequestType, hex string, isRelease bool) interface{} {

	jsonStr := ResponseData{}
	// 获取新的链接

	send := utils.Unhexlify(hex)
	vCore, _ := f.Pools.Get()
	requestConn := vCore.(net.Conn)
	_, err_w := requestConn.Write(send)
	if err_w != nil {
		return jsonStr
	}

	// 获取 所有的 index 和 point name
	if responseCoverFunc.GetRequestType() == "01" {
		var items []Point
		result := MyPointBox{items}
		header := make([]byte, 9)
		num, err := requestConn.Read(header[:])
		if err == nil {
			if num > 0 && num == 9 {
				hexArray := utils.SplitSubN(utils.ByteArrayToHex(header[:num]), 2)
				jsonStr.Type = fmenum.GetTagType(hexArray[4])
				jsonStr.Status = fmenum.GetMsg(hexArray[5])
				// 检测头部是否是00
				if hexArray[6] == "00" {
					// 解析头部
					IndexPackageLen := utils.HexToDec(strings.Join(hexArray[6:9], ""))
					if IndexPackageLen == 0 {
						// 没有子集
						if isRelease {
							f.Pools.Close(vCore)
						} else {
							f.Pools.Put(vCore)
						}
						return jsonStr
					}
					var TempLen = 0
					for {
						// 找到一个点的的数据长度
						lenHeader := make([]byte, 1)
						num, err = requestConn.Read(lenHeader[:])
						if err == nil && num == 1 {
							if utils.ByteArrayToHex(lenHeader) == "00" {
								break
							}
							// 获取到数据长度
							utils.Debug("一组数据的 Hex " + utils.ByteArrayToHex(lenHeader))
							RecordDataLen := utils.HexToDec(utils.ByteArrayToHex(lenHeader))
							if TempLen == int(IndexPackageLen) {
								break
							}
							TempLen++
							//解析单条数据
							RecordData := make([]byte, RecordDataLen)
							num, err = requestConn.Read(RecordData[:])
							if err == nil {
								pointObj := SingleRecordCover(utils.SplitSubN(utils.ByteArrayToHex(RecordData), 2))
								result.AddPoint(pointObj)
							}
						}
					}
				}
			}
		}
		jsonStr.Msg = result.Items
		if isRelease {
			f.Pools.Close(vCore)
		} else {
			f.Pools.Put(vCore)
		}

		utils.Debug("get index and point name CoverResponse.go requestConn.Close()")
		return jsonStr
	}
	// 读取某些点位数据  release 1.0
	if responseCoverFunc.GetRequestType() == "02" || responseCoverFunc.GetRequestType() == "03" {

		err_read := requestConn.SetReadDeadline(time.Now().Add(200 * time.Millisecond)) // timeout
		if err_read != nil {
			utils.Debug("setReadDeadline failed:", err_read)
		}
		header := make([]byte, 6)
		var result []byte
		// 标准Header
		num1, err := requestConn.Read(header[:])
		if err == nil {
			if num1 > 0 && num1 == 6 {
				result = append(result, header[:num1]...)
				hexArray := utils.SplitSubN(utils.ByteArrayToHex(header[:num1]), 2)
				// 请求正确
				if hexArray[5] == "00" {
					lenData := make([]byte, 2)
					num2, err0 := requestConn.Read(lenData[:])
					if err0 == nil && num2 > 0 && num2 == 2 {
						result = append(result, lenData[:num2]...)
						// 获取数据体
						readDataLen := utils.HexToDec(utils.ByteArrayToHex(lenData[:num2]))
						if readDataLen == 0 {
							jsonStr.DataCode = "NULL"
							if isRelease {
								f.Pools.Close(vCore)
							} else {
								f.Pools.Put(vCore)
							}
							utils.Debug("DataLen = 0 CoverResponse.go requestConn.Close()")
							return jsonStr
						}
						data := make([]byte, readDataLen)
						utils.Debug("Package Data Length >", readDataLen)
						num3, err1 := requestConn.Read(data[:])
						if err1 == nil {
							result = append(result, data[:num3]...)
							debug2 := utils.SplitSubN(utils.ByteArrayToHex(result), 2)
							utils.Debug("result 2 > ", debug2)
							lenData2 := make([]byte, 2)
							num4, err2 := requestConn.Read(lenData2[:])
							debug3 := utils.SplitSubN(utils.ByteArrayToHex(lenData2), 2)
							utils.Debug("result 3 > ", debug3)
							if num4 == 2 && err2 == nil && strings.Join(utils.SplitSubN(utils.ByteArrayToHex(lenData2), 2), "") == "0000" {
								strHex := utils.ByteArrayToHex(result)
								hexArray = utils.SplitSubN(strHex, 2)
								jsonStr.Type = fmenum.GetTagType(hexArray[4])
								jsonStr.Status = fmenum.GetMsg(hexArray[5])
								jsonStr.ReturnCode = strings.Join(hexArray[0:2], "")
								jsonStr.Msg = CoverGetPointValueByIndexAndType(hexArray)
								jsonStr.DataCode = "00"
								utils.Debug("CoverConnResponseToJSON 02 03 ", strHex)
							}
							//回收链接
							if isRelease {
								f.Pools.Close(vCore)
							} else {
								f.Pools.Put(vCore)
							}
							utils.Debug("Successfully CoverResponse.go requestConn.Close()")

							return jsonStr
						}
					}
				} else {
					jsonStr.DataCode = "XX"
					jsonStr.Msg = nil
					//回收链接
					if isRelease {
						f.Pools.Close(vCore)
					} else {
						f.Pools.Put(vCore)
					}
					utils.Debug("DataNull CoverResponse.go requestConn.Close() OK")
					return jsonStr
				}
			}
		} else {
			jsonStr.DataCode = "XX"
			jsonStr.Msg = nil
			if err, ok := err.(net.Error); ok && err.Timeout() {
				//回收链接
				if isRelease {
					f.Pools.Close(vCore)
				} else {
					f.Pools.Put(vCore)
				}
				utils.Debug("TimeOut CoverResponse.go requestConn.Close() OK")
				return jsonStr
			}
			if isRelease {
				f.Pools.Close(vCore)
			} else {
				f.Pools.Put(vCore)
			}
			return jsonStr
		}

		utils.Debug("Default CoverResponse.go requestConn.Close()")
		//回收链接
		if isRelease {
			f.Pools.Close(vCore)
		} else {
			f.Pools.Put(vCore)
		}
		return jsonStr
	}

	//修改数据返回值进行处理 release 1.0
	if responseCoverFunc.GetRequestType() == "05" {
		t4 := time.Now()
		data := make([]byte, 10)
		num, err := requestConn.Read(data[:])
		if err == nil && num == 10 {
			hexArray := utils.SplitSubN(utils.ByteArrayToHex(data[:num]), 2)
			utils.Debug("修改变量返回值：", hexArray)
			jsonStr.IsModify = CoverModifyPointValueByIndexAndType(hexArray)
			jsonStr.Type = fmenum.GetTagType(hexArray[4])
			jsonStr.Status = fmenum.GetMsg(hexArray[5])
			jsonStr.ReturnCode = strings.Join(hexArray[0:2], "")
			jsonStr.DataCode = "00"
			//回收链接
			t5 := time.Now()
			utils.Debug("Use Time = ", t5.Sub(t4))
			utils.Debug("Modify CoverResponse.go requestConn.Close()")
			if isRelease {
				f.Pools.Close(vCore)
			} else {
				f.Pools.Put(vCore)
			}
			return jsonStr
		} else if err != nil {
			if isRelease {
				f.Pools.Close(vCore)
			} else {
				f.Pools.Put(vCore)
			}
			utils.Debug("Error CoverResponse.go requestConn.Close()", err)
			return jsonStr
		}
	}
	if isRelease {
		f.Pools.Close(vCore)
	} else {
		f.Pools.Put(vCore)
	}
	return jsonStr
}

func (resp ResponseData) GetValInt() int64 {
	if len(resp.Msg) > 0 {
		int64Data, _ := strconv.ParseInt(resp.Msg[0].NewValue, 10, 64)
		return int64Data
	}
	return 0
}

func (resp ResponseData) GetValFloat() float64 {
	if len(resp.Msg) > 0 {
		float64Data, _ := strconv.ParseFloat(resp.Msg[0].NewValue, 64)
		return float64Data
	}
	return 0
}

// RespCoverDataR4R8 @title 获取Response 的hex str 单独处理 Double Float Int
// @description  获取 Response 并且按照规则解析
// @auth      作者  张建新    时间（2021/6/1 22:57 ）
// @param     hexData        []string         "hex string 数组"
// @return    result         Point          "构建的多个Point对象"
func RespCoverDataR4R8(hexData []string, r int) []Point {
	cursor := 0
	items := []Point{}
	result := MyPointBox{items}
	rDataLen := utils.HexToDec(strings.Join(hexData[0:2], ""))
	hexData = hexData[2:]
	for {
		__temp := Point{}
		if cursor == int(rDataLen) {
			break
		}
		index_ := utils.HexToDec(strings.Join(hexData[cursor:cursor+3], ""))
		if r == 4 {
			qty := utils.HexToDec(strings.Join(hexData[cursor+3:cursor+3+1], ""))
			r8R4Data, _ := utils.HexToFloat32(strings.Join(hexData[cursor+3+1:cursor+3+1+r], ""))
			cursor = cursor + 3 + 1 + r
			__temp.Index = index_
			__temp.NewValue = r8R4Data
			__temp.Qty = qty
			result.AddPoint(__temp)
			continue
		}

		if r == 8 {
			qty := utils.HexToDec(strings.Join(hexData[cursor+3:cursor+3+1], ""))
			r8R4Data, _ := utils.HexToFloat64(strings.Join(hexData[cursor+3+1:cursor+3+1+r], ""))
			cursor = cursor + 3 + 1 + r
			__temp.Index = index_
			__temp.NewValue = r8R4Data
			__temp.Qty = qty
			result.AddPoint(__temp)
			continue
		}

		if r == 1 {
			qty := CoverI1Qty(strings.Join(hexData[cursor+3:cursor+3+r], ""))
			r8R4Data := int(utils.HexToDec(strings.Join(hexData[cursor+3:cursor+3+r], "")))
			cursor = cursor + 3 + r
			__temp.Index = index_
			__temp.NewValue = strconv.FormatInt(int64(r8R4Data), 10)
			__temp.Qty = qty
			result.AddPoint(__temp)
			continue
		}
	}
	return result.Items
}

func CoverI1Qty(hex string) (qty int64) {
	// qty = 8x 设备链接失败
	if hex[:1] == "8" {
		qty = 1
	}
	// qty = 0x 点位设备链接正常
	if hex[:1] == "0" {
		qty = 0
	}
	return qty
}
