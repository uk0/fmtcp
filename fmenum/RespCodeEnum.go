package fmenum

var statusMap = make(map[string]string)

func init() {
	statusMap["00"] = "无错误"
	statusMap["01"] = "变量类型不正确"
	statusMap["02"] = "变量实际数量零"
	statusMap["03"] = "变量开始索引错误"
	statusMap["04"] = "读取变量数量错误"
	statusMap["05"] = "禁止修改"
	statusMap["06"] = "修改变量全部失败"
	statusMap["07"] = "修改变量部分失败"
	statusMap["08"] = "批量变量名称格式错误"
	statusMap["FE"] = "加密狗授权不支持此协议"
	statusMap["FF"] = "访问运行数据库失败"
	statusMap["GG"] = "请注意:点位没有在Func[GetCacheByPointName]找到Index请确认[config.json]文件是否正确,发送请求被取消."
}

func GetMsg(code string) string {
	return statusMap[code]
}
