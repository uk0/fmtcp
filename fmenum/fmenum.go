package fmenum

type MessageType string

type Enum interface {
	name() string
	ordinal() int
	values() *[]string
}

const (
	VT MessageType = "09"

	//"""下面是 modify get index 使用"""
	R8_AI MessageType = "01"
	R8_AO MessageType = "02"
	R8_AR MessageType = "03"
	R8_VA MessageType = "07"
	I1_DI MessageType = "04"
	I1_DO MessageType = "05"
	I1_DR MessageType = "06"
	I1_VD MessageType = "08"

	//""" R4 区分"""

	R4_AI MessageType = "11"
	R4_AO MessageType = "12"
	R4_AR MessageType = "13"
	R4_VA MessageType = "17"
)

func (c MessageType) GetTagType() string {
	return string(c)
}

func (c MessageType) GetTag() string {
	return string(c)
}

var tagNameMapALL = make(map[string]string)
var tagNameMapUse = make(map[string]string)

func init() {
	tagNameMapALL["09"] = "VT"
	//"""下面是 modify get index 使用"""
	tagNameMapALL["01"] = "R8_AI"
	tagNameMapALL["02"] = "R8_AO"
	tagNameMapALL["03"] = "R8_AR"
	tagNameMapALL["07"] = "R8_VA"
	tagNameMapALL["04"] = "I1_DI"
	tagNameMapALL["05"] = "I1_DO"
	tagNameMapALL["06"] = "I1_DR"
	tagNameMapALL["08"] = "I1_VD"

	//""" R4 区分"""

	tagNameMapALL["11"] = "R4_AI"
	tagNameMapALL["12"] = "R4_AO"
	tagNameMapALL["13"] = "R4_AR"
	tagNameMapALL["17"] = "R4_VA"

	tagNameMapUse["01"] = "AI"
	tagNameMapUse["02"] = "AO"
	tagNameMapUse["03"] = "AR"
	tagNameMapUse["07"] = "VA"
	tagNameMapUse["04"] = "DI"
	tagNameMapUse["05"] = "DO"
	tagNameMapUse["06"] = "DR"
	tagNameMapUse["08"] = "VD"
	tagNameMapUse["09"] = "VT"

}

func GetTagType(code string) string {
	return tagNameMapUse[code]
}

/***
AI=01,AO=02,AR=03,VA=07
DI=04,DO=05,DR=06,VD=08
VT=09
*/
func GetALL() map[string]string {
	return tagNameMapALL
}

func GetUse() map[string]string {
	return tagNameMapUse
}
