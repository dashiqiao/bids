package rpc

import (
	"context"
	"log"

	. "platform_report/config"

	"github.com/smallnest/rpcx/client"
)

type RespData struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func RpcCall(funcName string, body interface{}) (*RespData, error) {
	cf := Conf{}
	config := cf.GetConf()

	d := client.NewPeer2PeerDiscovery("tcp@"+config.RpcHost, "")
	xClient := client.NewXClient("rpcService", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xClient.Close()

	resp := new(RespData)
	err := xClient.Call(context.Background(), funcName, body, resp)
	if err != nil {
		log.Println(err.Error())
		return resp, err
	}
	return resp, nil
}
