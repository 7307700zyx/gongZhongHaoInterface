package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	tools "gitlab.wsmfin.com/DEV/GoLangTools"

	"gongZhongHaoInterface/common"
	consts "gongZhongHaoInterface/conf"
	"net/http"
)

func WxAddKfAccount(w http.ResponseWriter, r *http.Request) {
	common.AccessToken, common.Err = common.GetAccessToken(common.ConRedis)
	if common.Err != nil {
		w.Write([]byte("123123123"))
		return
	}

	getMedisCountUri := consts.ADD_KF_ACCOUNT_URI + common.AccessToken
	fmt.Println(getMedisCountUri)

	////map解析成json
	kfAccount := make(map[string]string)
	kfAccount["kf_account"] = "test1@test"
	kfAccount["nickname"] = "客服1"
	kfAccount["password"] = "pswmd5"
	b_kfAccount, _ := json.Marshal(kfAccount)
	s_kfAccount := string(b_kfAccount)
	jsonArray := tools.SendPostRequstJson(getMedisCountUri, s_kfAccount)
	fmt.Println(jsonArray)
	errcode := gjson.Get(jsonArray, "errcode").String()
	errmsg := gjson.Get(jsonArray, "errmsg").String()
	fmt.Println(errcode)
	fmt.Println(errmsg)

	w.Write([]byte("123123123"))
	return
}


func WxKfSend(msg_type string,media_id string,MsgToUserName string , MsgFromUserName string  ){
	common.AccessToken, common.Err = common.GetAccessToken(common.ConRedis)
	if common.Err != nil {
		fmt.Println("获取token失败")
		return
	}
	getMedisCountUri := consts.KF_SEND_URI + common.AccessToken
	fmt.Println(getMedisCountUri)
	msg := make(map[string]interface{})
	son_map := make(map[string]string)
	msg["touser"] = MsgToUserName

	var send_data string
	switch msg_type {
	case "news":
		////map解析成json
		msg["msgtype"] = msg_type
		son_map["media_id"] = media_id
		msg["mpnews"] = son_map
		b_msg, _ := json.Marshal(msg)
		send_data = string(b_msg)
		break
	case "text":
		////map解析成json
		msg["msgtype"] = msg_type
		son_map["content"] = media_id
		msg["text"] = son_map
		b_msg, _ := json.Marshal(msg)
		send_data = string(b_msg)
		break
	case "image":
		////map解析成json
		msg["msgtype"] = msg_type
		son_map["media_id"] = media_id
		msg["image"] = son_map
		b_msg, _ := json.Marshal(msg)
		send_data = string(b_msg)
		break
	default:
		////map解析成json
		msg["msgtype"] = msg_type
		son_map["content"] = media_id
		msg["text"] = son_map
		b_msg, _ := json.Marshal(msg)
		send_data = string(b_msg)
		break
	}

	fmt.Println(send_data)
	jsonArray := tools.SendPostRequstJson(getMedisCountUri, send_data)
	fmt.Println(jsonArray)

}