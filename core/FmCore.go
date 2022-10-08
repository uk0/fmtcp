package core

import (
	"fmgolib/fmenum"
	"fmgolib/utils"
	"fmt"
	pool "github.com/silenceper/pool"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	FmgolibGlobal *FC
)

type FC struct {
	Address    string
	PNameCache map[string]Data
	Pools      pool.Pool
	GlobalMx   sync.Mutex
}

var DefaultFailReadResponse = func() ResponseData {
	return ResponseData{
		Status:   fmenum.GetMsg("GG"),
		DataCode: "XX",
	}
}

type FmTuple struct {
	Index int
	Value string
}
type ModifyTemp struct {
	Type  string
	Value []FmTuple
}

type ModifyStruct struct {
	PointName string
	Value     string
}

func NewFmCore(host string, config string) (fc *FC) {

	factory := func() (interface{}, error) { return net.Dial("tcp", host) }

	close := func(v interface{}) error { return v.(net.Conn).Close() }
	//创建一个连接池： 初始化5，最大连接40，空闲连接数是20
	poolConfig := &pool.Config{
		InitialCap: 5,
		MaxIdle:    20,
		MaxCap:     40,
		Factory:    factory,
		Close:      close,
		//Ping:       ping,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: 58 * time.Second,
	}
	pool, err := pool.NewChannelPool(poolConfig)

	if err != nil {
		utils.Info("err=", err)
	}

	fc = &FC{
		Address: host,
		Pools:   pool,
	}
	if _, err := os.Stat(config); err == nil || os.IsExist(err) {
		data, _ := os.ReadFile(config)
		json.Unmarshal(data, &fc.PNameCache)
		utils.Info("配置文件已经存在，准备加载。")
		utils.Info("配置文件内容如下", string(data))
	} else {
		fc.PNameCache = map[string]Data{}
		utils.Info("没有找到配置文件config.json 准备第一次初始化生成文件。")
		for key := range fmenum.GetUse() {
			hex := fc.GetALlPointAndIndexByTypeTagStr(key, 0, 0)
			temp := fc.CoverConnResponseToJSON(fmenum.GET_ALL_INDEX_AND_POINT_NAME, hex, false)
			jsonStr := temp.(ResponseData)
			fc.PNameCache[jsonStr.Type] = jsonStr.Msg
		}

		jsonData, err := json.Marshal(fc.PNameCache)
		utils.Info("參數已經打印")
		utils.Info(string(jsonData))
		jsonFile, err := os.Create(config)
		if err != nil {
			panic(err)
		}
		defer jsonFile.Close()

		jsonFile.Write(jsonData)
		jsonFile.Close()
		fmt.Println("JSON data written to ", jsonFile.Name())
	}

	FmgolibGlobal = fc
	go fc.getLenConn()
	fmt.Println("FmCore.go Cache Init successfully")
	return
}

func (f *FC) CreateConn() (*net.TCPConn, error) {
	addr, _ := net.ResolveTCPAddr("tcp", f.Address)
	conn1, err1 := net.DialTCP("tcp", nil, addr)
	if err1 == nil {
		_ = conn1.SetKeepAlive(true)
		utils.Info("FmCore.go 创建连接成功")
	} else {
		utils.Info("FmCore.go DView服务已经掉线，等待下次轮训，")
	}
	return conn1, err1
}

// GetALlPointAndIndexByTypeTag 获取某个类型的所有点位 的 Index 和 Name
// 获取某个类型的所有点位Name 以及 Index
//  每次只能get一种tag的所有名字 和 index 注意 ：高字节前,低字节后
//**
func (f *FC) GetALlPointAndIndexByTypeTag(tag string, startIndex int, getCount int) string {

	/***"""
	tag = AI=01,AO=02,AR=03,VA=07
	DI=04,DO=05,DR=06,VD=08
	VT=09
	"""
	**/
	result := "3e2a271b"
	result += tag
	//# 默认大端
	result += fmt.Sprintf("%06x", startIndex)
	result += fmt.Sprintf("%06x", getCount)

	return result
}

func (f *FC) GetALlPointAndIndexByTypeTagStr(tag string, startIndex int, getCount int) string {

	/***"""
	tag = AI=01,AO=02,AR=03,VA=07
	DI=04,DO=05,DR=06,VD=08
	VT=09
	"""
	**/
	result := "3e2a271b"
	result += tag
	//# 默认大端
	result += fmt.Sprintf("%06x", startIndex)
	result += fmt.Sprintf("%06x", getCount)

	return result
}

func (f *FC) GenGetPointValueByIndexAndTypeTag(tag string, getIndex int, getCount int, timeWindows int) string {
	result := "3e2a271c"
	result += tag
	// 默认大端
	result += fmt.Sprintf("%06x", getIndex)
	result += fmt.Sprintf("%06x", getCount)
	result += fmt.Sprintf("%04x", timeWindows)
	utils.Debug("send > ", result)
	return result
}

func (f *FC) GenModifyPointByIndexAndTypeTag(tag string, _value FmTuple) string {

	doubles := []string{"02", "03", "07"}
	flouts := []string{"12", "13", "17"}
	bool := []string{"05", "06", "08", "04"}
	str := []string{"09"}
	result := ""
	if utils.StrInSlice(tag, doubles) {
		result += fmt.Sprintf("%06x", _value.Index)
		result += utils.DoubleToHex(_value.Value)
		utils.Debug("send > ", result)
		return result
	}
	if utils.StrInSlice(tag, flouts) {
		result += fmt.Sprintf("%06x", _value.Index)
		result += utils.Float32ToHex(_value.Value)
		utils.Debug("send > ", result)
		return result
	}
	if utils.StrInSlice(tag, bool) {
		result += fmt.Sprintf("%06x", _value.Index)
		//TODO 临时处理方式
		result += utils.BoolToHex(_value.Value)
		//result += utils.Int64ToHex(_value.Value)
		utils.Debug("send > ", result)
		return result
	}
	if utils.StrInSlice(tag, str) {
		result += fmt.Sprintf("%06x", _value.Index)
		sendData := utils.CharToHex([]byte(_value.Value), "utf8->gbk")
		sendLen := len(sendData) / 2
		result += fmt.Sprintf("%02x", sendLen)
		result += sendData
		utils.Debug("send > ", result)
		return result
	}
	return result
}

func (f *FC) GenModifyPointValueByIndex(tag string, lists []FmTuple) string {
	result := "3e2a271d"
	result += tag
	temp := ""
	for _, tuple_ := range lists {
		temp += f.GenModifyPointByIndexAndTypeTag(tag, tuple_)
	}
	send_len := len(temp)
	send_len = send_len / 2
	lenHex := fmt.Sprintf("%04x", send_len)
	result += lenHex
	result += temp
	utils.Debug("send > ", result)
	return result
}

func (f *FC) GenGetPointValueByIndexList(tag string, timeWindows int, indexList []string) string {
	result := "3e2a271e"
	result += tag
	//默认大端
	result += fmt.Sprintf("%04x", timeWindows)
	result += fmt.Sprintf("%04x", len(indexList)*3)
	for _, data := range indexList {
		result += utils.Int32ToHex(data)
	}
	utils.Debug("send > ", result)
	return result
}

//////////////////////////////////////////////////////////////////

func (f *FC) FuncTestCache(pointName string) (string, int64) {
	key, index := f.GetCacheByPointName(pointName)
	return key, index
}

func (f *FC) FuncGetPointValueByName(pointName string) ResponseData {
	key, index := f.GetCacheByPointName(pointName)
	if key == "NOT_FOUND_POINT" {
		//TODO 长度为0 找不到点位不进行查询
		return DefaultFailReadResponse()
	}
	hex := f.GenGetPointValueByIndexAndTypeTag(key, int(index), 1, 0)
	jsonStr := f.CoverConnResponseToJSON(fmenum.GET_POINT_VALUE_BY_INDEX, hex, false)
	return jsonStr.(ResponseData)
}

func (f *FC) FuncGetPointValueByNameUseOne(pointName string) ResponseData {
	key, index := f.GetCacheByPointName(pointName)
	if key == "NOT_FOUND_POINT" {
		//TODO 长度为0 找不到点位不进行查询
		return DefaultFailReadResponse()
	}
	hex := f.GenGetPointValueByIndexAndTypeTag(key, int(index), 1, 0)
	jsonStr := f.CoverConnResponseToJSON(fmenum.GET_POINT_VALUE_BY_INDEX, hex, true)
	return jsonStr.(ResponseData)
}

//GetCacheByPointName 获取缓存并且返回点位类型 以及要操作的点位坐标
func (f *FC) GetCacheByPointName(pointName string) (k string, idx int64) {
	for types, v := range f.PNameCache {
		for key, value := range fmenum.GetUse() {
			if value == types {
				for _, d := range v {
					if strings.ToUpper(strings.TrimSpace(d.Name)) == strings.ToUpper(strings.TrimSpace(pointName)) {
						k = key
						idx = d.Index
						return key, d.Index
					}
				}
			}
		}
	}
	return "NOT_FOUND_POINT", 0
}

//FuncGetPointValueByListName 获取某些点位的数据信息 @pointNames
func (f *FC) FuncGetPointValueByListName(pointNames []string) (jsonStr []ResponseData) {

	for _, pointName := range pointNames {
		key, index := f.GetCacheByPointName(pointName)
		if key == "NOT_FOUND_POINT" {
			//TODO 长度为0 找不到点位不进行查询
			return []ResponseData{DefaultFailReadResponse()}
		}
		// 注意当 Index = 0 的时候 Count 必定不能为 0
		hex := f.GenGetPointValueByIndexAndTypeTag(key, int(index), 1, 0)
		interfaceData := f.CoverConnResponseToJSON(fmenum.GET_POINT_VALUE_BY_INDEX, hex, false)

		temp := interfaceData.(ResponseData)
		letPoints := []Point{}
		for _, point := range temp.Msg {
			point.Name = pointName
			letPoints = append(letPoints, point)
		}
		temp.Msg = letPoints
		jsonStr = append(jsonStr, temp)
	}
	return jsonStr
}

func (f *FC) FuncGetPointValueByNameListLink(pointNames []string) (jsonStr []ResponseData) {
	//TODO 自动Close
	GetDataMap := map[string][]map[string]int64{}
	for _, pointName := range pointNames {
		key, index := f.GetCacheByPointName(pointName)
		if key == "NOT_FOUND_POINT" {
			// 长度为0 找不到点位不进行查询
			return []ResponseData{DefaultFailReadResponse()}
		}
		if _, ok := GetDataMap[key]; ok {
			temp := GetDataMap[key]
			temp = append(temp, map[string]int64{pointName: index})
			GetDataMap[key] = temp
		} else {
			var indexList []map[string]int64
			indexList = append(indexList, map[string]int64{pointName: index})
			GetDataMap[key] = indexList
		}

	}
	utils.Debug("GetDataMap>", GetDataMap)
	for k, v := range GetDataMap {
		_, tmpV := utils.Map2IntSlice(v)
		hex := f.GenGetPointValueByIndexList(k, 0, tmpV)
		interfaceData := f.CoverConnResponseToJSON(fmenum.GET_POINT_VALUE_BY_INDEX, hex, true)
		utils.Debug(interfaceData)
		temp := interfaceData.(ResponseData)
		var letPoints []Point
		for _, point := range temp.Msg {
			point.Name = utils.MMapUseIndexGetName(GetDataMap, point.Index)
			letPoints = append(letPoints, point)
		}
		temp.Msg = letPoints
		jsonStr = append(jsonStr, temp)
	}
	return jsonStr
}

//FuncGetPointValueByListNameUseTimeRange 获取某些点位的数据信息 @pointNames @sec
func (f *FC) FuncGetPointValueByListNameUseTimeRange(pointNames []string, sec int) (jsonStr []ResponseData) {
	//TODO 自动Close
	for _, pointName := range pointNames {
		key, index := f.GetCacheByPointName(pointName)
		if key == "NOT_FOUND_POINT" {
			// 长度为0 找不到点位不进行查询
			return []ResponseData{DefaultFailReadResponse()}
		}
		// 注意当 Index = 0 的时候 Count 必定不能为 0
		hex := f.GenGetPointValueByIndexAndTypeTag(key, int(index), 1, sec)
		interfaceData := f.CoverConnResponseToJSON(fmenum.GET_POINT_VALUE_BY_INDEX, hex, true)
		temp := interfaceData.(ResponseData)
		var letPoints []Point
		for _, point := range temp.Msg {
			point.Name = pointName
			letPoints = append(letPoints, point)
		}
		temp.Msg = letPoints
		jsonStr = append(jsonStr, temp)
	}

	return jsonStr
}

func (f *FC) FuncModifyValueByName(indexListAndValue []ModifyStruct) (resp ResponseData) {
	// 写锁定
	var modifyList []ModifyTemp
	for _, modifyStruct := range indexListAndValue {
		t1 := time.Now()
		key, index := f.GetCacheByPointName(modifyStruct.PointName)
		if key == "NOT_FOUND_POINT" {
			//TODO 长度为0 找不到点位不进行查询
			return DefaultFailReadResponse()
		}
		t2 := time.Now()
		utils.Info("GetCacheByPointName 单次命中缓存消耗时间", t2.Sub(t1))
		var modifyListTemp []FmTuple
		modifyListTemp = append(modifyListTemp, FmTuple{
			Index: int(index),
			Value: modifyStruct.Value,
		})
		modifyList = append(modifyList, ModifyTemp{
			Type:  key,
			Value: modifyListTemp,
		})
	}
	// 同一个类型的 需要修改的数据
	for _, ModifyTempData := range modifyList {
		hexGenResp := f.GenModifyPointValueByIndex(ModifyTempData.Type, ModifyTempData.Value)
		// 修改使用短链接
		resp = f.CoverConnResponseToJSON(fmenum.MODIFY_POINT_VALUE, hexGenResp, false).(ResponseData)
	}

	return
}

func (f *FC) getLenConn() {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			utils.Debug("当前剩余连接数", f.Pools.Len())
		}
	}
}

// 备用
func (f *FC) FuncModifyValueByNameUseSingle(indexListAndValue []ModifyStruct, rw string) (resp ResponseData) {

	var modifyList []ModifyTemp
	for _, modifyStruct := range indexListAndValue {
		key, index := f.GetCacheByPointName(modifyStruct.PointName)
		if key == "NOT_FOUND_POINT" {
			//TODO 长度为0 找不到点位不进行查询
			return DefaultFailReadResponse()
		}
		var modifyListTemp []FmTuple
		modifyListTemp = append(modifyListTemp, FmTuple{
			Index: int(index),
			Value: modifyStruct.Value,
		})
		modifyList = append(modifyList, ModifyTemp{
			Type:  key,
			Value: modifyListTemp,
		})
	}
	// 同一个类型的 需要修改的数据
	for _, ModifyTempData := range modifyList {
		hexGenResp := f.GenModifyPointValueByIndex(ModifyTempData.Type, ModifyTempData.Value)
		//3C 2A 27 1D 09 00 00 01 00 00
		resp = f.CoverConnResponseToJSON(fmenum.MODIFY_POINT_VALUE, hexGenResp, false).(ResponseData)
	}

	return
}
