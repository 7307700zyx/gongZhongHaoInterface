package controllers

import (
	tools "gitlab.wsmfin.com/DEV/GoLangTools"

	//"fmt"
	"github.com/beevik/etree"
	"gongZhongHaoInterface/common"
	consts "gongZhongHaoInterface/conf"
	"io/ioutil"
	"net/http"
)

func textMsg(w http.ResponseWriter, RequestBody string) {
	msg_type_str := "【文本消息】-"
	//tools.OutPutInfo(nil,msg_type_str+"xml内容:"+RequestBody)
	doc := etree.NewDocument()
	if err := doc.ReadFromString(RequestBody); err != nil {
		tools.OutPutInfo(nil, msg_type_str+" 系统异常，xml内容解析失败:"+RequestBody)
		return
	}
	var Content string
	var MsgToUserName string
	var MsgFromUserName string
	for _, xmls := range doc.SelectElements("xml") {
		if msg_content := xmls.SelectElement("Content"); msg_content != nil {
			Content = msg_content.Text()
		}
		if msg_ToUserName := xmls.SelectElement("ToUserName"); msg_ToUserName != nil {
			MsgToUserName = msg_ToUserName.Text()
		}
		if msg_FromUserName := xmls.SelectElement("FromUserName"); msg_FromUserName != nil {
			MsgFromUserName = msg_FromUserName.Text()
		}
	}
	//tools.OutPutInfo(nil,msg_type_str+"xml内容解析成功:用户openid："+MsgFromUserName+",发送内容："+Content)

	//if(Content == "kftext"){
	//	WxKfSend("text","这里是客服发送text消息",MsgFromUserName,MsgToUserName)
	//	return
	//}
	//if(Content == "kftvoice"){
	//	WxKfSend("voice","这里是客服发送voice消息",MsgFromUserName,MsgToUserName)
	//	return
	//}
	//if(Content == "kfvideo"){
	//	WxKfSend("video","这里是客服发送video消息",MsgFromUserName,MsgToUserName)
	//	return
	//}
	//if(Content == "kfimage"){
	//	WxKfSend("image","这里是客服发送image消息",MsgFromUserName,MsgToUserName)
	//	return
	//}
	//if(Content == "kfmpnews"){
	//	WxKfSend("news","这里是客服发送news消息",MsgFromUserName,MsgToUserName)
	//	return
	//}

	keyWordSql := "SELECT * FROM fa_media WHERE keyword = '" + Content + "' LIMIT 0,1"
	//fmt.Println(keyWordSql)
	result, err := common.MysqlQuery("", "", "", "", "", keyWordSql, consts.MYSQL_DSN)
	if err != nil {
		tools.OutPutInfo(err, msg_type_str+"查询关键字表查询失败，关键字："+Content+"，异常sql："+keyWordSql)
		return
	}

	//如果不是关键字，判断是否是参加活动
	var xmlStr string
	resCnt := len(result)
	if resCnt <= 0 {
		tools.OutPutInfo(err, msg_type_str+Content+"不在关键字列表中，查看是否是活动内的关键字")
		WxActivityText(w, MsgFromUserName, MsgToUserName, Content)
		return
	}
	res := result[0]
	//tools.OutPutInfo(err,msg_type_str+Content+"在关键字列表中，查看关键字详情：",res)
	if res["type"] != "" {
		if res["type"] == "news" {
			tools.OutPutInfo(err, msg_type_str+Content+"在关键字列表中，回复内容为【图文】类型，发送客服消息：", res)
			//xmlStr = common.ReplyNews(MsgFromUserName, MsgToUserName, res["title"], res["description"], res["picUrl"], res["url"])
			WxKfSend("news", res["media_id"], MsgFromUserName, MsgToUserName)
		}
		if res["type"] == "voice" {
			tools.OutPutInfo(err, msg_type_str+Content+"在关键字列表中，回复内容为【音频】类型：", res)
			xmlStr = common.ReplyVoice(MsgFromUserName, MsgToUserName, res["media_id"])
		}
		if res["type"] == "image" {
			tools.OutPutInfo(err, msg_type_str+Content+"在关键字列表中，回复内容为【图片】类型：", res)
			xmlStr = common.ReplyImg(MsgFromUserName, MsgToUserName, res["media_id"])
		}
		if res["type"] == "text" {
			tools.OutPutInfo(err, msg_type_str+Content+"在关键字列表中，回复内容为【文本】类型：", res)
			xmlStr = common.ReplyText(MsgFromUserName, MsgToUserName, res["title"])
		}
		if res["type"] == "video" {
			tools.OutPutInfo(err, msg_type_str+Content+"在关键字列表中，回复内容为【视频】类型：", res)
			xmlStr = common.ReplyVideo(MsgFromUserName, MsgToUserName, res["media_id"], res["title"], res["description"])
		}

		//tools.OutPutInfo(err,msg_type_str+"在关键字列表中，回复内容为【视频】类型，完整XML：",xmlStr)
		w.Write([]byte(xmlStr))
	}
}

func imgMsg(w http.ResponseWriter, RequestBody string) {
	msg_type_str := "【图片消息】-"
	doc := etree.NewDocument()
	if err := doc.ReadFromString(RequestBody); err != nil {
		tools.OutPutInfo(nil, msg_type_str+" 系统异常，xml内容读取失败:"+RequestBody)
		return
	}
	var PicUrl string
	var MediaId string
	var MsgFromUserName string
	var MsgToUserName string
	for _, xmls := range doc.SelectElements("xml") {
		if msg_picurl := xmls.SelectElement("PicUrl"); msg_picurl != nil {
			PicUrl = msg_picurl.Text()
		}
		if msg_mediaid := xmls.SelectElement("MediaId"); msg_mediaid != nil {
			MediaId = msg_mediaid.Text()
		}
		if msg_FromUserName := xmls.SelectElement("FromUserName"); msg_FromUserName != nil {
			MsgFromUserName = msg_FromUserName.Text()
		}
		if msg_ToUserName := xmls.SelectElement("ToUserName"); msg_ToUserName != nil {
			MsgToUserName = msg_ToUserName.Text()
		}
	}

	tools.OutPutInfo(nil, "PicUrl:"+PicUrl+";MediaId:"+MediaId)
	WxActivityImage(w, MsgFromUserName, MsgToUserName, PicUrl, MediaId)
	return
}

func voiceMsg(w http.ResponseWriter, RequestBody string) {
	msg_type_str := "【音频消息】-"
	doc := etree.NewDocument()
	if err := doc.ReadFromString(RequestBody); err != nil {
		tools.OutPutInfo(nil, msg_type_str+" 系统异常，xml内容读取失败:"+RequestBody)
		return
	}
	var MsgFromUserName string
	var MsgToUserName string
	for _, xmls := range doc.SelectElements("xml") {
		if msg_FromUserName := xmls.SelectElement("FromUserName"); msg_FromUserName != nil {
			MsgFromUserName = msg_FromUserName.Text()
		}
		if msg_ToUserName := xmls.SelectElement("ToUserName"); msg_ToUserName != nil {
			MsgToUserName = msg_ToUserName.Text()
		}
	}
	xmlStr := common.ReplyTextCommon(MsgFromUserName, MsgToUserName)
	//tools.OutPutInfo(nil,msg_type_str+"调用默认回复的完整XML：",xmlStr)
	w.Write([]byte(xmlStr))
}

func videoMsg(w http.ResponseWriter, RequestBody string) {
	msg_type_str := "【视频消息】-"
	doc := etree.NewDocument()
	if err := doc.ReadFromString(RequestBody); err != nil {
		tools.OutPutInfo(nil, msg_type_str+" 系统异常，xml内容读取失败:"+RequestBody)
		return
	}
	var MsgFromUserName string
	var MsgToUserName string
	for _, xmls := range doc.SelectElements("xml") {
		if msg_FromUserName := xmls.SelectElement("FromUserName"); msg_FromUserName != nil {
			MsgFromUserName = msg_FromUserName.Text()
		}
		if msg_ToUserName := xmls.SelectElement("ToUserName"); msg_ToUserName != nil {
			MsgToUserName = msg_ToUserName.Text()
		}
	}
	xmlStr := common.ReplyTextCommon(MsgFromUserName, MsgToUserName)
	//tools.OutPutInfo(nil,msg_type_str+"调用默认回复的完整XML：",xmlStr)
	w.Write([]byte(xmlStr))
}

func eventMsg(w http.ResponseWriter, RequestBody string) {
	msg_type_str := "【事件消息】-"
	doc := etree.NewDocument()
	if err := doc.ReadFromString(RequestBody); err != nil {
		tools.OutPutInfo(nil, msg_type_str+" 系统异常，xml内容读取失败:"+RequestBody)
		return
	}
	var EventFromUserName string
	var EventToUserName string
	var EventType string
	for _, xmls := range doc.SelectElements("xml") {
		if msg_FromUserName := xmls.SelectElement("FromUserName"); msg_FromUserName != nil {
			EventFromUserName = msg_FromUserName.Text()
		}
		if msg_ToUserName := xmls.SelectElement("ToUserName"); msg_ToUserName != nil {
			EventToUserName = msg_ToUserName.Text()
		}
		if msg_Event := xmls.SelectElement("Event"); msg_Event != nil {
			EventType = msg_Event.Text()
		}
	}

	if EventType == consts.MESSAGE_EVENT_SUBSCRIBE {
		xmlStr := common.ReplyText(EventFromUserName, EventToUserName, consts.REPLY_TEXT_SUBSCRIBE)
		//tools.OutPutInfo(nil,msg_type_str+"关注事件 回复的完整XML：",xmlStr)
		w.Write([]byte(xmlStr))
	} else {
		xmlStr := common.ReplyTextCommon(EventFromUserName, EventToUserName)
		w.Write([]byte(xmlStr))
	}

}

func WxMsg(w http.ResponseWriter, r *http.Request) {
	common.AccessToken, common.Err = common.GetAccessToken(common.ConRedis)
	if common.Err != nil {
		tools.OutPutInfo(common.Err, "获取token失败")
		w.Write([]byte("获取token失败"))
		return
	}

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			tools.OutPutInfo(err, "读取传送的数据失败")
			w.Write([]byte("读取传送的数据失败"))
			return
		}
		msg_xml := string(body)
		doc := etree.NewDocument()
		if err := doc.ReadFromString(msg_xml); err != nil {
			tools.OutPutInfo(err, "解析用户发送的信息失败")
			w.Write([]byte("解析用户发送的信息失败"))
			return
		}
		var msg_type_str string
		for _, xmls := range doc.SelectElements("xml") {
			if msg_type := xmls.SelectElement("MsgType"); msg_type != nil {
				msg_type_str = msg_type.Text()
			}
		}
		//tools.OutPutInfo(nil,"msg_type_str:"+msg_type_str)
		//tools.OutPutInfo(nil,"=====================================")

		switch msg_type_str {
		case consts.MESSAGE_TEXT:
			textMsg(w, msg_xml)
		case consts.MESSAGE_IMAGE:
			imgMsg(w, msg_xml)
		case consts.MESSAGE_NEWS:
			imgMsg(w, msg_xml)
		case consts.MESSAGE_VOICE:
			voiceMsg(w, msg_xml)
		case consts.MESSAGE_VIDEO:
			videoMsg(w, msg_xml)
			//事件推送
		case consts.MESSAGE_EVENT:
			eventMsg(w, msg_xml)
		default:
			textMsg(w, msg_xml)
		}
		return
	}
	return
}
