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

//func textMsg(w http.ResponseWriter,RequestBody string){
//	doc := etree.NewDocument()
//	if err := doc.ReadFromString(RequestBody); err != nil {
//		w.Write([]byte("success"))
//	}
//	var Content string
//	//var MsgId string
//	var MsgToUserName string
//	var MsgFromUserName string
//	for _, xmls := range doc.SelectElements("xml") {
//		if msg_content := xmls.SelectElement("Content"); msg_content != nil {
//			Content = msg_content.Text()
//		}
//		if msg_ToUserName := xmls.SelectElement("ToUserName"); msg_ToUserName != nil {
//			MsgToUserName = msg_ToUserName.Text()
//		}
//		if msg_FromUserName := xmls.SelectElement("FromUserName"); msg_FromUserName != nil {
//			MsgFromUserName = msg_FromUserName.Text()
//		}
//	}
//	productInooSql := "SELECT * FROM media WHERE keyword = '"+Content+"' LIMIT 0,1"
//	fmt.Println(productInooSql)
//	result,err := common.MysqlQuery("","","","","",productInooSql,consts.MYSQL_DSN)
//	if(err != nil){
//		fmt.Println(err)
//		return
//	}
//
//	var xmlStr string
//	resCnt := len(result)
//	if(resCnt <= 0){
//		fmt.Println("未找到对应关键字")
//		return
//	}
//	res := result[0]
//	fmt.Println(res)
//	if(res["type"] != ""){
//		if(res["type"] == "news"){
//			xmlStr = common.ReplyNews(MsgFromUserName,MsgToUserName , res["title"], res["description"], res["picUrl"], res["url"])
//		}
//		if(res["type"] == "voice"){
//			xmlStr = common.ReplyVoice(MsgFromUserName,MsgToUserName , res["media_id"])
//		}
//		if(res["type"] == "image"){
//			xmlStr = common.ReplyImg(MsgFromUserName,MsgToUserName , res["media_id"])
//		}
//		if(res["type"] == "text"){
//			xmlStr = common.ReplyText(MsgFromUserName,MsgToUserName , res["title"])
//		}
//		if(res["type"] == "video"){
//			xmlStr = common.ReplyVideo(MsgFromUserName,MsgToUserName , res["media_id"], res["title"], res["description"])
//		}
//
//		w.Write([]byte(xmlStr))
//	}
//}
//
//func imgMsg(w http.ResponseWriter,RequestBody string){
//	doc := etree.NewDocument()
//	if err := doc.ReadFromString(RequestBody); err != nil {
//		//fmt.Println(err)
//		w.Write([]byte("success"))
//	}
//	var PicUrl string
//	var MediaId string
//	for _, xmls := range doc.SelectElements("xml") {
//		if msg_picurl := xmls.SelectElement("PicUrl"); msg_picurl != nil {
//			//fmt.Println(msg_picurl.Text())
//			//"WECHAT API 连接redis失败"
//			PicUrl = msg_picurl.Text()
//		}
//		if msg_mediaid := xmls.SelectElement("MediaId"); msg_mediaid != nil {
//			//fmt.Println(msg_mediaid.Text())
//			MediaId = msg_mediaid.Text()
//		}
//	}
//	log.Println(PicUrl)
//	log.Println(MediaId)
//	//fmt.Println(PicUrl)
//	//fmt.Println(MediaId)
//	//fmt.Println("----------imgMsg  msg_type_str------------")
//	w.Write([]byte("123123123"))
//}
//
//func voiceMsg(w http.ResponseWriter,RequestBody string) {
//	//fmt.Println("----------voiceMsg------------")
//}
//func videoMsg(w http.ResponseWriter,RequestBody string)  {
//	//fmt.Println("----------videoMsg------------")
//}

func WxAddKfAccount(w http.ResponseWriter, r *http.Request) {
	common.AccessToken, common.Err = common.GetAccessToken(common.ConRedis)
	if common.Err != nil {
		w.Write([]byte("123123123"))
		return
	}

	getMedisCountUri := consts.GET_MEDIA_LIST_URI + common.AccessToken
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
	//
	//if r.Method == "POST" {
	//	body, err := ioutil.ReadAll(r.Body)
	//	if err != nil {
	//		log.Fatal(err)
	//		w.Write([]byte("123123123"))
	//		return
	//	}
	//	msg_xml := string(body)
	//	doc := etree.NewDocument()
	//	if err := doc.ReadFromString(msg_xml); err != nil {
	//		w.Write([]byte("123123123"))
	//		return
	//	}
	//	var msg_type_str string
	//	for _, xmls := range doc.SelectElements("xml") {
	//		if msg_type := xmls.SelectElement("MsgType"); msg_type != nil {
	//			msg_type_str = msg_type.Text()
	//		}
	//	}
	//	switch msg_type_str {
	//	case consts.MESSAGE_TEXT:
	//		textMsg(w,msg_xml)
	//	case consts.MESSAGE_IMAGE:
	//		imgMsg(w,msg_xml)
	//	case consts.MESSAGE_NEWS:
	//		imgMsg(w,msg_xml)
	//	case consts.MESSAGE_VOICE:
	//		voiceMsg(w,msg_xml)
	//	case consts.MESSAGE_VIDEO:
	//		videoMsg(w,msg_xml)
	//	default:
	//		textMsg(w,msg_xml)
	//	}
	//	return
	//}

	w.Write([]byte("123123123"))
	return
}
