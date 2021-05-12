package main

import (
	tools "gitlab.wsmfin.com/DEV/GoLangTools"
	"gongZhongHaoInterface/common"
	consts "gongZhongHaoInterface/conf"
	controller "gongZhongHaoInterface/controller"
	"net/http"
	"os"

	//"strconv"
	//"time"
)

var IpLast string

//var TimeStr string = time.Now().Format("2006-01-02")
//var MyLogger *log.Logger

func main() {
	tools.Logger = tools.WriteLogFile("./wx_blog", "")
	res, localHost, _ := tools.GetLocalIp()
	if res == false {
		tools.OutPutInfo(nil,"WECHAT API,error,获取本机ip失败")
		return
	}

	server := http.Server{Addr: localHost + consts.WECHAT_PORT}
	tools.OutPutInfo(nil,"WECHAT API" + localHost + consts.WECHAT_PORT)

	common.ConRedis, common.Err = common.ConnectRedis()
	if common.Err != nil {
		tools.OutPutInfo(nil,"WECHAT API 连接redis失败")
		os.Exit(1)
	}

	//查看WECHAT API wxConnect
	http.HandleFunc("/wxConnect", controller.WxConnect)

	//GET 获取根据传来的media_id 获取微信临时素材(当前仅下载图片类别)
	//http://10.1.193.136:8088/getTmpMedia?media_id=sFCzs7j5X6WqK7KEAm5rGqDSpymkg7OPqJhT7VX5nP1R77Q9QuR5EihA7suktyFx
	http.HandleFunc("/getTmpMedia", controller.GetTmpMedia)

	//POST 获取根据传来的media_id 获取微信永久素材(当前仅下载图片类别)
	//http://10.1.193.136:8088/getMedia?media_id=sFCzs7j5X6WqK7KEAm5rGqDSpymkg7OPqJhT7VX5nP1R77Q9QuR5EihA7suktyFx
	http.HandleFunc("/getMedia", controller.GetMedia)

	//获取微信素材数量
	//http://10.1.193.136:8088/getMediaCount
	http.HandleFunc("/getMediaCount", controller.WxMediaCount)

	//获取微信素材列表并入库
	//素材的类型，图片（image）、视频（video）、语音 （voice）、图文（news）
	//http://10.1.193.136:8088/getAllMediasToDB?media_type=news
	http.HandleFunc("/getAllMediasToDB", controller.WxMediasToDB)

	//增加客服
	http.HandleFunc("/addKfAccount", controller.WxAddKfAccount)

	//图片识别完毕
	//为朋友圈活动的指定openid，发送客服消息，让其进行下一步活动。
	http.HandleFunc("/addPyqScore", controller.WxActivityPYQAddScore)

	tools.OutPutInfo(nil,"WECHAT API启动成功，监听中...")

	server.ListenAndServe()
}
