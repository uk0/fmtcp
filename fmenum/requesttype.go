package fmenum

type RequestType string

const (
	//""" 获取 所有的 index 和 point name"""

	GET_ALL_INDEX_AND_POINT_NAME RequestType = "01"

	//"""通过变量索引读取变量值,或某段时间内发生变化的变量值 1c"""

	GET_POINT_VALUE_BY_INDEX RequestType = "02"

	//通过变量索引,选择读取某些变量值 1e

	GET_POINT_VALUE_BY_INDEX_LIST RequestType = "03"

	//通过变量索引,批量修改变量值、

	MODIFY_POINT_VALUE RequestType = "05"
)

func (c RequestType) GetRequestType() string {
	return string(c)
}
