package service

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
	. "platform_report/config"
	"strconv"
)

type UserService struct {
	Meta string
	Host string
}

type RespData struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func NewUserService() *UserService {
	cf := Conf{}
	config := cf.GetConf()
	return &UserService{Host: config.RpcHost, Meta: "rpcService"}
}

func (c *UserService) GetUserName(userId interface{}) interface{} {
	reqMap := make(map[string]interface{})
	switch vType := userId.(type) {
	case int:
		reqMap["userId"] = int64(vType)
	case int64:
		reqMap["userId"] = vType
	case string:
		reqMap["userId"], _ = strconv.ParseInt(vType, 10, 64)
	default:
		fmt.Printf("unexpected type %T\n", vType)
	}

	userInfo, err := c.Call(c.Host, c.Meta, "GetUserInfo", reqMap)
	if err != nil {
		return ""
	}
	return (userInfo.Data.(map[string]interface{}))["Realname"]
}

func (c *UserService) Call(host, meta, method string, body interface{}) (*RespData, error) {
	d := client.NewPeer2PeerDiscovery("tcp@"+host, "")
	xclient := client.NewXClient(meta, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	resp := new(RespData)
	err := xclient.Call(context.Background(), method, body, resp)
	return resp, err
}
