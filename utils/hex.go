package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"strconv"
	"unsafe"
)

func Unhexlify(str string) []byte {
	res := make([]byte, 0)
	for i := 0; i < len(str); i += 2 {
		x, _ := strconv.ParseInt(str[i:i+2], 16, 32)
		res = append(res, byte(x))
	}
	return res
}

func Hex2Char(str string) []byte {
	return Unhexlify(str)
}

func ByteArrayToHex(hexByte []byte) string {
	return hex.EncodeToString(hexByte)
}

func CharToHex(hexByte []byte, types string) string {
	var re []byte
	if types == "utf8->gbk" {
		re, _ = Utf8ToGbk(hexByte)
	}
	if types == "gbk->utf8" {
		re, _ = GbkToUtf8(hexByte)
	}
	return hex.EncodeToString(re)
}

func DoubleToHex(dstr string) string {
	dl_, _ := strconv.ParseFloat(dstr, 64)
	b1 := make([]byte, 8)
	binary.LittleEndian.PutUint64(b1, uint64(math.Float64bits(dl_)))
	//fmt.Println(b1)
	return fmt.Sprintf("%x", b1)
}

func Float32ToHex(fstr string) string {
	fl_, _ := strconv.ParseFloat(fstr, 32)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(math.Float32bits(float32(fl_))))
	return fmt.Sprintf("%x", b)
}

func HexToFloat32(hexStr string) (f string, err error) {
	b, _ := hex.DecodeString(hexStr)
	decimalNum := decimal.NewFromFloat32(Float32frombytes(b))
	f = decimalNum.String()
	return
}

func HexToFloat64(hexStr string) (f string, err error) {
	b, _ := hex.DecodeString(hexStr)
	decimalNum := decimal.NewFromFloat(Float64frombytes(b))
	f = decimalNum.String()
	return
}
func HexToByteArray(hexStr string) []byte {
	b, _ := hex.DecodeString(hexStr)
	return b
}

func Int64ToHex(fstr string) string {
	i, _ := strconv.ParseInt(fstr, 10, 64)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return fmt.Sprintf("%x", b)
}

func BoolToHex(str string) string {
	if str == "1" {
		return "01"
	}
	if str == "0" {
		return "00"
	}
	// TODO Default value setting
	return "00"
}

/// Dview 特殊处理
func Int32ToHex(fstr string) string {
	i, _ := strconv.ParseInt(fstr, 10, 32)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	reInt := fmt.Sprintf("%x", b)
	return reInt[2:]
}

func Int32ToHex2(dec int64) string {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(dec))
	reInt := fmt.Sprintf("%x", b)
	return reInt[4:]
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func IsLittleEndian() bool {
	var i int32 = 0x01020304
	// 下面这两句是为了将int32类型的指针转换为byte类型的指针
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb // 取得pb位置对应的值
	// 由于b是byte类型的,最多保存8位,那么只能取得开始的8位
	// 小端: 04 (03 02 01)
	// 大端: 01 (02 03 04)
	return (b == 0x04)
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-2; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func HexToDec(str string) int64 {
	num, _ := strconv.ParseInt(str, 16, 64)
	return num
}
