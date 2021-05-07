package controllers

import (
	"crypto/sha1"
	"fmt"
	//"localHost"
	consts "gongZhongHaoInterface/conf"
	//controller "gongZhongHaoInterface/controller"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

func makeSignature(timestamp, nonce string) string { //本地计算signature
	si := []string{consts.TOKEN, timestamp, nonce}
	sort.Strings(si)            //字典序排序
	str := strings.Join(si, "") //组合字符串
	s := sha1.New()             //返回一个新的使用SHA1校验的hash.Hash接口
	io.WriteString(s, str)      //WriteString函数将字符串数组str中的内容写入到s中
	return fmt.Sprintf("%x", s.Sum(nil))

}

func validateUrl(w http.ResponseWriter, r *http.Request) bool {
	timestamp := strings.Join(r.Form["timestamp"], "")
	nonce := strings.Join(r.Form["nonce"], "")
	signature := strings.Join(r.Form["signature"], "")
	echostr := strings.Join(r.Form["echostr"], "")
	fmt.Println("-----------------------------")

	fmt.Println(signature)

	signatureGen := makeSignature(timestamp, nonce)
	fmt.Println(signatureGen)
	fmt.Println("-----------------------------")
	if signatureGen != signature {
		return false
	}
	fmt.Fprintf(w, echostr) //原样返回eechostr给微信服务器
	return true
}

func WxConnect(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Request需要解析
	if !validateUrl(w, r) {
		log.Println("Wechat Service: This http request is not from wechat platform")
		return
	}
	log.Println("validateUrl Ok")
	//WxMediaList(w, r)
	fmt.Println("---------------------")
	WxMsg(w, r)

}

//func main() {
//	log.Println("Wechat Service: Start!")
//	http.HandleFunc("/", procSignature)
//	err := http.ListenAndServe(":80", nil)
//	if err != nil {
//		log.Println("Wechat Service: ListenAndServe Error: ", err)
//	}
//	log.Println("Wechat Service: Stop!")
//}
