package utils

import (
	"bytes"
	"strconv"
	"strings"
)

func StrInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MMapUseIndexGetName(x map[string][]map[string]int64, index int64) (pointName string) {
	for _, array := range x { //第一层Map
		for _, data := range array {
			for k, v := range data {
				if v == index {
					pointName = k
				}
			}
		}
	}
	return pointName
}

func Map2IntSlice(m []map[string]int64) ([]string, []string) {
	slK, slV := make([]string, 0), make([]string, 0)
	for _, mMap := range m {
		for k, v := range mMap {
			slK = append(slK, k)
			slV = append(slV, strconv.FormatInt(v, 10))
		}
	}
	return RemoveDuplicateElement2(slK), RemoveDuplicateElement2(slV)
}

func SplitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}

func RemoveDuplicateElement(languages []string) []string {
	result := make([]string, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func RemoveDuplicateElement2(arrays []string) []string {
	check := make(map[string]int)
	res := make([]string, 0)
	for _, val := range arrays {
		check[val] = 1
	}

	for letter := range check {
		letter = strings.Replace(letter, " ", "", -1)
		// 去除换行符
		letter = strings.Replace(letter, "\n", "", -1)
		if letter != "" {
			res = append(res, letter)
		}
	}

	return res
}

func DeleteRepeat(list []string) []string {
	mapdata := make(map[string]interface{})
	if len(list) <= 0 {
		return nil
	}
	// 利用key的唯一性，将key对应的value置为true，同时将重复的数组元素过滤
	for _, v := range list {
		mapdata[v] = "true"
	}
	var datas []string
	// 将mapdata中的key，拼接在新的切片中
	for k := range mapdata {
		if k == "" {
			continue
		}
		datas = append(datas, k)
	}
	return datas
}
