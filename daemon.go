package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"gitlab.wsmfin.com/DEV/AutoLR3Consts"
	"gitlab.wsmfin.com/DEV/GoLangTools"
	"net/http"
	"os"
	"strings"
	//"strconv"
	"bufio"
	//"regexp"
	//"log"
	//"io"
	"io/ioutil"
	"log"
	//"os/exec"
	"os/exec"
	"time"
)

var ConRedis *redis.Client
var IpLast string

var TimeStr string = time.Now().Format("2006-01-02")
var MyLogger *log.Logger

func writeLogSimple(funcName string, msg string, res bool) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	if res == false {
		// 写入log文件格式： 2018/07/31 17:28:21 4.Println log without log.LstdFlags ...
		ConRedis.LPush("logs:"+IpLast, "["+timeStr+"] 【error AutoLR3-daemon】 Api "+funcName+": "+msg+"\n ")
		MyLogger.Println("logs:" + IpLast + "[" + timeStr + "] 【error AutoLR3-daemon】 Api " + funcName + ": " + msg + "\n ")
	} else {
		ConRedis.LPush("logs:"+IpLast, "["+timeStr+"] 【success AutoLR3-daemon】 Api "+funcName+": "+msg+"\n ")
		MyLogger.Println("logs:" + IpLast + "[" + timeStr + "] 【success AutoLR3-daemon】 Api " + funcName + ": " + msg + "\n ")
	}
}

/**
 * 查看是否开启
 */
func PingPong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG"))
	writeLogSimple("PingPong", "访问pingpong接口", true)
}

/**
 * 工作端操作
 * get 获取传递来的 type 参数。
 * 如果传递参数 ?type=start。则杀死并开启工作端
 * 如果未传递参数或?type=不为start，则仅查看工作端是否启动
 * 响应 error/success
 */
func WorkHandle(w http.ResponseWriter, r *http.Request) {
	writeLogSimple("WorkHandle", "访问WorkHandle接口", true)
	query := r.URL.Query()
	start := query.Get("type")
	if start == "start" {
		//step1.杀死进程
		res, _, err := tools.KillProcessByName(consts.APP_TITLE_WORK + ".exe")
		if res == false {
			w.Write([]byte("error"))
			writeLogSimple("WorkHandle?type=start", "杀死进程失败"+err.Error(), false)
		}
		//_ , err = tools.StartProcess(consts.WORK_PATH+"/"+consts.APP_TITLE_WORK+".exe")
		//if err != nil{
		//	w.Write([]byte("error"))
		//	writeLogSimple("WorkHandle?type=start","开启进程失败"+err.Error(),false)
		//}
		err = exec.Command(`cmd`, `/c`, `start`, consts.WORK_PATH+"\\"+consts.APP_TITLE_WORK+".exe").Start()
		if err != nil {
			w.Write([]byte("error"))
			writeLogSimple("UpdateHandle", "启动"+consts.WORK_PATH+"/"+consts.APP_TITLE_WORK+".exe"+"失败"+err.Error(), false)
			return
		}

		w.Write([]byte("success"))
		writeLogSimple("WorkHandle?type=start", "重启工作端进程成功", true)
	} else {
		res, _, _ := tools.IsProcessExist(consts.APP_TITLE_WORK + ".exe")
		if res == true {
			writeLogSimple("WorkHandle", "访问WorkHandle接口SUCCESS", true)
			w.Write([]byte("success"))
		} else {
			w.Write([]byte("error"))
			writeLogSimple("WorkHandle", "访问WorkHandle接口ERROR", false)
		}
	}
}

/*
 *	工作端更新
 */
func UpdateHandle(w http.ResponseWriter, r *http.Request) {

	writeLogSimple("UpdateHandle", "调用守护进程UpdateHandle接口，执行中...", true)
	//step1.杀死进程
	res, out, err := tools.KillProcessByName(consts.APP_TITLE_WORK + ".exe")
	if res == false {
		w.Write([]byte(out))
		writeLogSimple("UpdateHandle", "杀死进程失败->"+out+"->"+err.Error(), false)
		return
	}
	tools.Sleep(1000)
	//step2.清空目录
	err = os.RemoveAll(consts.WORK_PATH)
	if err != nil {
		//w.Write([]byte("清空目录"+consts.WORK_PATH+"失败"))
		w.Write([]byte("error"))
		writeLogSimple("UpdateHandle", "清空目录"+consts.WORK_PATH+"失败 "+err.Error(), false)
		return
	}

	////step3.创建目录
	err = os.Mkdir(consts.WORK_PATH, 0777)
	//err = tools.MakeAllDir(consts.WORK_PATH)
	if err != nil {
		w.Write([]byte("error"))
		writeLogSimple("UpdateHandle", "创建目录"+consts.WORK_PATH+"失败"+err.Error(), false)
		return
	}

	//step4.读取传递来的文件并放到指定地址
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//获取文件流,第三个返回值是错误对象
	//file, header, err := r.FormFile(consts.WORK_PARAMETERS)
	file, _, err := r.FormFile(consts.WORK_PARAMETERS)
	if err != nil {
		w.Write([]byte("error"))
		writeLogSimple("UpdateHandle", "文件读取失败失败"+err.Error(), false)
		return
	}

	//读取文件流为[]byte
	b, _ := ioutil.ReadAll(file)
	//把文件保存到指定位置	"C:/AUTOLR3/"+"AutoLR3.zip"
	ioutil.WriteFile(consts.WORK_PATH+"/"+consts.ZIP_NAME, b, 0777)
	//关闭文件
	err = file.Close()
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//解压上传文件
	err = tools.DeCompressByPath(consts.WORK_PATH+"/"+consts.ZIP_NAME, consts.WORK_PATH+"/")
	if err != nil {
		w.Write([]byte("error"))
		writeLogSimple("UpdateHandle", "解压"+consts.WORK_PATH+"/"+consts.ZIP_NAME+"失败"+err.Error(), false)
		return
	}

	//_,err = tools.StartProcess(consts.WORK_PATH+"/"+consts.APP_TITLE_WORK+".exe")
	err = exec.Command(`cmd`, `/c`, `start`, consts.WORK_PATH+"\\"+consts.APP_TITLE_WORK+".exe").Start()
	if err != nil {
		w.Write([]byte("error"))
		writeLogSimple("UpdateHandle", "启动"+consts.WORK_PATH+"/"+consts.APP_TITLE_WORK+".exe"+"失败"+err.Error(), false)
		return
	}

	w.Write([]byte("success"))
	writeLogSimple("UpdateHandle", "更新升级成功", true)
	return
}

func ShutDownHandle(w http.ResponseWriter, r *http.Request) {
	tools.ShutDownEXE()
	w.Write([]byte("success"))
	writeLogSimple("ShutDownHandle", "关闭计算机完毕", true)
}

func ReShutDownHandle(w http.ResponseWriter, r *http.Request) {
	tools.ReShutDownEXE()
	w.Write([]byte("success"))
	writeLogSimple("ShutDownHandle", "重启计算机完毕", true)
}

func main() {
	//file, err := os.Create(consts.DAEMON_PATH+"\\"+TimeStr+".log")
	//if err != nil {
	//	log.Fatalln("fail to create "+TimeStr+".log")
	//}
	//MyLogger = log.New(file, "", log.Ldate|log.Ltime)
	//
	//MyLogger.SetFlags(log.LstdFlags)    // 设置写入文件的log日志的格式

	res, workeHost, err := tools.GetLocalIp()
	if res == false {
		fmt.Println("error,获取本机ip失败")
		MyLogger.Println("守护进程,error,获取本机ip失败")
		return
	}
	fmt.Println(err)
	IpLast = strings.Split(workeHost, ".")[3]

	//读取工作端配置文件内的 redisHost字段
	//redisHost := ReadLineFile(consts.WORK_PATH+"/conf/app.conf")
	//fmt.Println(redisHost)
	//if len(redisHost) <= 0{
	//	redisHost = consts.REDIS_HOST_TMP
	//	fmt.Println("error,读取工作端配置文件内的 redisHost字段失败")
	//	fmt.Println("使用默认redis地址配置"+redisHost)
	//	MyLogger.Println("守护进程,error,读取工作端配置文件内的 redisHost字段失败")
	//	MyLogger.Println("使用默认redis地址配置"+redisHost)
	//}
	//for i:= 1;i<=120;i++{
	//	//连接redis
	//	ConRedis, err = tools.Conredis(redisHost+":6380", consts.REDIS_PWD, "0")
	//	if err != nil {
	//		fmt.Println("守护进程,error,连接redis失败,失败次数"+strconv.Itoa(i))
	//		MyLogger.Println("守护进程,error,连接redis失败,失败次数"+strconv.Itoa(i))
	//		tools.Sleep(1000)
	//		continue
	//	}else{
	//		break
	//	}
	//}
	//fmt.Println("守护进程,success,连接redis成功")

	fmt.Println("守护进程 start")
	server := http.Server{Addr: workeHost + consts.DAEMON_PORT}
	fmt.Println(workeHost + consts.DAEMON_PORT)
	fmt.Println("守护进程" + workeHost + consts.DAEMON_PORT)
	//查看守护进程pingpong
	http.HandleFunc("/"+consts.DAEMON_PING, PingPong)
	//更新工作端
	http.HandleFunc("/"+consts.DAEMON_UPDATE, UpdateHandle)
	//查看工作端是否工作
	http.HandleFunc("/"+consts.DAEMON_WORK_HANDLE, WorkHandle)
	//关闭计算机
	http.HandleFunc("/"+consts.SHUTDOWN, ShutDownHandle)
	//重启计算机
	http.HandleFunc("/"+consts.RE_SHUTDOWN, ReShutDownHandle)

	fmt.Println("守护进程启动成功，监听中...")
	server.ListenAndServe()
}

func ReadLineFile(fileName string) string {
	redisHost := ""
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("read file fail", err)
		return redisHost
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), consts.REDIS_HOST_KEY) {
			splitRedisConf := strings.Split(scanner.Text(), "=")
			redisHostList := strings.Fields(splitRedisConf[1])
			redisHost := strings.Join(redisHostList, ",")
			return redisHost
		}
	}
	return redisHost
}
