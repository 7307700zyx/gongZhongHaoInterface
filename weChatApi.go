package main

import (
	"fmt"
	tools "gitlab.wsmfin.com/DEV/GoLangTools"
	"gongZhongHaoInterface/common"
	consts "gongZhongHaoInterface/conf"
	controller "gongZhongHaoInterface/controller"
	"log"
	"net/http"
	"os"

	//"strconv"
	"time"
)

var IpLast string

var TimeStr string = time.Now().Format("2006-01-02")
var MyLogger *log.Logger

///**
// * 查看是否开启
// */
//func PingPong(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("PONG"))
//	writeLogSimple("PingPong", "访问pingpong接口", true)
//}
//
//func ShutDownHandle(w http.ResponseWriter, r *http.Request) {
//	tools.ShutDownEXE()
//	w.Write([]byte("success"))
//	writeLogSimple("ShutDownHandle", "关闭计算机完毕", true)
//}

func ReShutDownHandle(w http.ResponseWriter, r *http.Request) {
	//验证用户名密码，如果成功则header里返回session，失败则返回StatusUnauthorized状态码

	w.WriteHeader(http.StatusOK)
	if (r.Form.Get("user") == "admin") && (r.Form.Get("pass") == "888") {
		w.Write([]byte("hello,验证成功！"))
	} else {
		w.Write([]byte("hello,验证失败了！"))
	}
}

func main() {

	file, err := os.Create(consts.WECHAT_PATH + "\\" + consts.WECHAT_LOG + "_" + TimeStr + ".log")
	if err != nil {
		log.Fatalln("fail to create " + TimeStr + ".log")
	}
	MyLogger = log.New(file, "", log.Ldate|log.Ltime)

	MyLogger.SetFlags(log.LstdFlags) // 设置写入文件的log日志的格式

	res, localHost, err := tools.GetLocalIp()
	if res == false {
		//fmt.Println("error,获取本机ip失败")
		log.Println("error,获取本机ip失败")
		MyLogger.Println("WECHAT API,error,获取本机ip失败")
		return
	}

	server := http.Server{Addr: localHost + consts.WECHAT_PORT}
	fmt.Println(localHost + consts.WECHAT_PORT)
	//fmt.Println(localHost + consts.WECHAT_PORT)
	log.Println(localHost + consts.WECHAT_PORT)
	//fmt.Println("WECHAT API" + localHost + consts.WECHAT_PORT)

	log.Println("WECHAT API" + localHost + consts.WECHAT_PORT)
	common.ConRedis, common.Err = common.ConnectRedis()
	if common.Err != nil {
		log.Println("WECHAT API 连接redis失败")
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

	log.Println("WECHAT API启动成功，监听中...")
	server.ListenAndServe()
}
