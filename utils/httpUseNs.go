package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var SpecifiDnsServerTransport http.RoundTripper = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		Resolver: &net.Resolver{PreferGo: true, Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			address = "192.168.0.20:53"
			return d.DialContext(ctx, network, address)
		}},
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var SpecifiDnsServerClient = http.Client{Transport: SpecifiDnsServerTransport}

func Test() {

	//resp, err := SpecifiDnsServerClient.Get("http://rocketmq-dashboard.rocketmq-operator.svc.cluster.local:8080/rocketmq/nsaddr")
	resp, err := SpecifiDnsServerClient.Get("http://10.20.192.111:8080/rocketmq/nsaddr")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
