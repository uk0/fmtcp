package utils

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func GetIPAddressByNS(domain string) []string {
	var nameSvr []string
	// 解析ip地址
	domainArray := strings.Split(domain, ":")
	ns, err := net.LookupHost(domainArray[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err: %s", err.Error())
	}

	// 对域名解析进行控制判断
	// 有些域名通常会先使用cname解析到一个别名上，然后再解析到实际的ip地址上
	switch {
	case len(ns) != 0:
		for _, n := range ns {
			nameSvr = append(nameSvr, fmt.Sprintf("%s:%s", n, domainArray[1]))
		}
		return nameSvr
	default:
		fmt.Println("default")
	}
	return nameSvr
}
