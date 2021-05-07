package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	tools "gitlab.wsmfin.com/DEV/GoLangTools"
	"gongZhongHaoInterface/common"
	consts "gongZhongHaoInterface/conf"
	"os"
	"strconv"

	//"math"
	"net/http"
	//"strconv"
	//"strconv"
	//"time"
)

func WxMediaCount(w http.ResponseWriter, r *http.Request) {
	common.AccessToken, common.Err = common.GetAccessToken(common.ConRedis)
	if common.Err != nil {
		w.Write([]byte("get token fail"))
		//common.MyLogger("get token fail")
		return
	}
	getMedisCountUri := consts.GET_MEDIA_COUNT_URI + common.AccessToken
	res, err := tools.SendPostRequst(getMedisCountUri, map[string]string{})
	if err != nil {
		w.Write([]byte("获取素材数量接口调用失败：" + getMedisCountUri))
		//common.MyTrace("获取素材数量接口调用失败："+getMedisCountUri)
		fmt.Println(err)
	}
	//{"voice_count":0,"video_count":0,"image_count":0,"news_count":0}
	voice_count := gjson.Get(res, "voice_count").String()
	video_count := gjson.Get(res, "video_count").String()
	image_count := gjson.Get(res, "image_count").String()
	news_count := gjson.Get(res, "news_count").String()
	common.ConRedis.Set("voice_count", voice_count, 0)
	common.ConRedis.Set("video_count", video_count, 0)
	common.ConRedis.Set("image_count", image_count, 0)
	common.ConRedis.Set("news_count", news_count, 0)

	w.Write([]byte(res))
	fmt.Println(err)
	fmt.Println(res)
	return
}

//{
//	"type":TYPE,	素材的类型，图片（image）、视频（video）、语音 （voice）、图文（news）
//	"offset":OFFSET,从全部素材的该偏移位置开始返回，0表示从第一个素材 返回
//	"count":COUNT	返回素材的数量，取值在1到20之间
//}
func WxMediaList(offset string, count string, media_type string) bool {

	common.AccessToken, common.Err = common.GetAccessToken(common.ConRedis)
	if common.Err != nil {
		fmt.Println("get token fail")
		return false
	}

	getMedisCountUri := consts.GET_MEDIA_LIST_URI + common.AccessToken
	fmt.Println(getMedisCountUri)

	////map解析成json
	medis := make(map[string]string)
	medis["type"] = media_type
	medis["offset"] = offset
	medis["count"] = count
	b_medis, _ := json.Marshal(medis)
	s_medis := string(b_medis)
	jsonArray := tools.SendPostRequstJson(getMedisCountUri, s_medis)
	fmt.Println("------------start-----------------")

	fmt.Println(jsonArray)

	dataSlice := make([]map[string]string, 0) //初始化数据切片

	if media_type == "news" {
		for _, res := range gjson.Get(jsonArray, "item").Array() {
			mapData := make(map[string]string)
			res.ForEach(func(key, value gjson.Result) bool {
				if key.String() == "media_id" {
					mapData["media_id"] = value.String()
				}
				if key.String() == "content" {
					jsonContents := value.String()
					title := gjson.Get(jsonContents, "news_item.0.title").String()
					mapData["title"] = title
					url := gjson.Get(jsonContents, "news_item.0.url").String()
					mapData["url"] = url
				}
				mapData["type"] = media_type
				//fmt.Println(key.String(), ":", value.String())
				return true
			})
			dataSlice = append(dataSlice, mapData)
		}
	} else {
		for _, res := range gjson.Get(jsonArray, "item").Array() {
			mapData := make(map[string]string)
			res.ForEach(func(key, value gjson.Result) bool {
				if key.String() == "media_id" {
					mapData["media_id"] = value.String()
				}
				if key.String() == "name" {
					mapData["title"] = value.String()
				}
				if key.String() == "url" {
					mapData["url"] = value.String()
				}
				mapData["type"] = media_type
				fmt.Println(key.String(), ":", value.String())
				return true
			})
			dataSlice = append(dataSlice, mapData)
		}
	}

	fmt.Println("------------end-----------------")
	fmt.Println(dataSlice)

	_, err, strSqlTmp := common.InsertAll(consts.MYSQL_DSN, "media", dataSlice)
	fmt.Println(err)
	fmt.Println(strSqlTmp)
	if err != nil {
		fmt.Println("素材列表插入数据库失败：" + strSqlTmp)
		return false
	}
	return true
}

func GetMedia(w http.ResponseWriter, r *http.Request) {
	common.AccessToken, common.Err = common.GetAccessToken(common.ConRedis)
	if common.Err != nil {
		fmt.Println("get token fail")
		return
	}
	query := r.URL.Query()
	media_id := query.Get("media_id")
	if media_id == "" {
		w.Write([]byte("素材id必传"))
		return
	}

	//根据media_id获取永久素材 https://api.weixin.qq.com/cgi-bin/material/get_material?access_token=ACCESS_TOKEN
	getMediaUri := consts.GET_MEDIA_BY_MEDIAID_URI + common.AccessToken
	media_map := make(map[string]string)
	media_map["media_id"] = media_id
	media_json, _ := json.Marshal(media_map)
	media_ids := string(media_json)

	dir, _ := os.Getwd()
	fmt.Println(dir)
	//创建文件
	path := dir + "/public/image/" + media_id + ".jpg"
	f1, err := os.Create(path)
	if err != nil {
		return
	}
	f1.Close()
	var jsonStr = []byte(media_ids)
	fmt.Println(getMediaUri)
	req, err := http.NewRequest("POST", getMediaUri, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 1024)
	f, err1 := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm) //可读写，追加的方式打开（或创建文件）
	if err1 != nil {
		panic(err1)
		return
	}
	defer f.Close()

	for {
		n, _ := resp.Body.Read(buf)
		if 0 == n {
			break
		}
		f.WriteString(string(buf[:n]))
	}

	if err != nil {
		return
	}
	return

}

func GetTmpMedia(w http.ResponseWriter, r *http.Request) {
	common.AccessToken, common.Err = common.GetAccessToken(common.ConRedis)
	if common.Err != nil {
		fmt.Println("get token fail")
		return
	}
	query := r.URL.Query()
	media_id := query.Get("media_id")
	if media_id == "" {
		w.Write([]byte("素材id必传"))
		return
	}
	dir, _ := os.Getwd()
	//创建文件
	path := dir + "/public/image/" + media_id + ".jpg"
	f1, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	f1.Close()

	//根据media_id获取临时素材 https://api.weixin.qq.com/cgi-bin/media/get?access_token=ACCESS_TOKEN&media_id=MEDIA_ID
	getTmpMediaUri := consts.GET_TMP_MEDIA_BY_MEDIAID_URI + "?access_token=" + common.AccessToken + "&media_id=" + media_id
	resp, err := http.Get(getTmpMediaUri)
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 1024)
	f, err1 := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm) //可读写，追加的方式打开（或创建文件）
	if err1 != nil {
		panic(err1)
		return
	}
	defer f.Close()

	for {
		n, _ := resp.Body.Read(buf)
		if 0 == n {
			break
		}
		f.WriteString(string(buf[:n]))
	}
	if err != nil {
		return
	}
	return
}

//const REDIS_KEY_MEDIS_VOICE string = "voice_count"
//const REDIS_KEY_MEDIS_VIDEO string = "video_count"
//const REDIS_KEY_MEDIS_IMAGE string = "image_count"
//const REDIS_KEY_MEDIS_NEWS  string = "news_count"
func WxMediasToDB(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	media_type := query.Get("media_type")
	if media_type == "" {
		w.Write([]byte("素材类别参数必传"))
		return
	}
	fmt.Println(media_type)
	media_total_num_key := r.FormValue("media_type_count")
	media_total_num_str, _ := common.ConRedis.Get(media_total_num_key).Result()
	fmt.Println(media_total_num_str)
	//总素材数
	media_total_num_int, _ := common.GetIntStrRes(media_total_num_str)

	//起始数0
	for page := 0; page <= media_total_num_int; page += consts.PAGE_COUNT {
		fmt.Println(page)
		fmt.Println(consts.PAGE_COUNT)
		res := WxMediaList(strconv.Itoa(page), strconv.Itoa(consts.PAGE_COUNT), media_type)
		if res != true {
			fmt.Println("素材数据获取失败")
			w.Write([]byte("素材列表数据入库失败"))
			return
		}
	}
	w.Write([]byte("素材列表数据入库完毕"))
	return
}
