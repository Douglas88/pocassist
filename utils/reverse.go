package utils

import (
	"bytes"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/url"
	"pocassist/basic"
	"time"
)

// use ceye api

func NewReverse() *Reverse {
	ceyeDomain := basic.GlobalConfig.Reverse.Domain
	flag := RandLowLetterNumber(8)
	if ceyeDomain == "" {
		return &Reverse{}
	}
	urlStr := fmt.Sprintf("http://%s.%s", flag, ceyeDomain)
	u, _ := url.Parse(urlStr)
	return &Reverse{
		Flag:				flag,
		Url:                ParseUrl(u),
		Domain:             u.Hostname(),
		Ip:                 "",
		IsDomainNameServer: false,
	}
}

func ReverseCheck(r *Reverse, timeout int64) bool {
	// todo
	ceyeApiToken := basic.GlobalConfig.Reverse.ApiKey
	if ceyeApiToken == "" || r.Domain == "" {
		return false
	}
	// 延迟 x 秒获取结果
	time.Sleep(time.Second * time.Duration(timeout))

	//check dns
	verifyUrl := fmt.Sprintf("http://api.ceye.io/v1/records?token=%s&type=dns&filter=%s", ceyeApiToken, r.Flag)
	if GetReverseResp(verifyUrl){
			return true
	} else {
		//	check request
		verifyUrl := fmt.Sprintf("http://api.ceye.io/v1/records?token=%s&type=http&filter=%s", ceyeApiToken, r.Flag)
		if GetReverseResp(verifyUrl){
			return true
		}
	}
	return false
}

func GetReverseResp(verifyUrl string) bool {
	notExist := []byte(`"data": []`)

	fastReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(fastReq)
	fastReq.SetRequestURI(verifyUrl)
	fastReq.Header.SetMethod(fasthttp.MethodGet)

	resp, err := DoFasthttpRequest(fastReq, false)

	if err != nil {
		return false
	}
	if !bytes.Contains(resp.Body, notExist) { // api返回结果不为空
		return true
	}
	return false
}
