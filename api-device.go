package hikvision

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/antchfx/xmlquery"
)

func GetDeviceInfo(host, username, password string) (*xmlquery.Node, error) {
	ctx := context.Background()
	r := NewDigestRequest(ctx, username, password) // username & password
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", host, "/ISAPI/System/deviceInfo"), nil)
	resp, _ := r.Do(req)
	defer resp.Body.Close()

	rsp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return xmlquery.Parse(strings.NewReader(string(rsp)))
}

// 获取系统能力集
func GetSystemCapability(host, username, password string) (*xmlquery.Node, error) {
	ctx := context.Background()
	r := NewDigestRequest(ctx, username, password) // username & password
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", host, "/ISAPI/System/capabilities"), nil)
	resp, _ := r.Do(req)
	defer resp.Body.Close()

	rsp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return xmlquery.Parse(strings.NewReader(string(rsp)))
}

// 获取车辆抓拍识别能力
func GetItcCapability(host, username, password string) (*xmlquery.Node, error) {
	ctx := context.Background()
	r := NewDigestRequest(ctx, username, password) // username & password
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", host, "/ISAPI/ITC/capabilities"), nil)
	resp, _ := r.Do(req)
	defer resp.Body.Close()

	rsp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return xmlquery.Parse(strings.NewReader(string(rsp)))
}

// 获取交通服务总能力
func GetTrafficCapability(host, username, password string) (*xmlquery.Node, error) {
	ctx := context.Background()
	r := NewDigestRequest(ctx, username, password) // username & password
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", host, "/ISAPI/Traffic/capabilities"), nil)
	resp, _ := r.Do(req)
	defer resp.Body.Close()

	rsp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return xmlquery.Parse(strings.NewReader(string(rsp)))
}

// 获取触发模式能力集
func GetTriggerModeCapability(host, username, password string) (*xmlquery.Node, error) {
	ctx := context.Background()
	r := NewDigestRequest(ctx, username, password) // username & password
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", host, "/ISAPI/ITC/TriggerMode/capabilities"), nil)
	resp, _ := r.Do(req)
	defer resp.Body.Close()

	rsp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return xmlquery.Parse(strings.NewReader(string(rsp)))
}
