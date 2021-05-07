package common

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
	"gitlab.wsmfin.com/DEV/GoLangTools"
	consts "gongZhongHaoInterface/conf"
	"log"
	"strconv"
	"time"
)

var ConRedis *redis.Client

//设置全局的error句柄
var Err error
var MyLogger *log.Logger
var AccessToken string

func ConnectRedis() (*redis.Client, error) {
	redisHost := consts.REDIS_HOST
	for j := 1; j <= consts.REDIS_CONN_TIMES; j++ {
		ConRedis, Err = tools.Conredis(redisHost+":"+consts.REDIS_PROT, consts.REDIS_PWD, "0")
		if Err != nil {
			log.Println(Err, "第", j, "次连接redis【异常】host:"+redisHost, "密码md5后:", tools.GetMd5Sum(consts.REDIS_PWD))
			tools.Sleep(1000)
		} else {
			log.Println(nil, "第", j, "次连接redis【正常】host:"+redisHost)
			break
		}
	}

	return ConRedis, Err
}

func ReplyText(MsgFromUserName string, MsgToUserName string, new_content string) string {
	create_doc := etree.NewDocument()
	xml := create_doc.CreateElement("xml")
	//xml.CreateComment("These are all known people")

	ToUserName := xml.CreateElement("ToUserName")
	ToUserName.CreateText(MsgFromUserName)

	FromUserName := xml.CreateElement("FromUserName")
	FromUserName.CreateText(MsgToUserName)

	now := time.Now().Unix() //获取时间戳
	now_string := strconv.FormatInt(now, 10)
	CreateTime := xml.CreateElement("CreateTime")
	CreateTime.CreateText(now_string)

	MsgType := xml.CreateElement("MsgType")
	MsgType.CreateText("text")

	rep_Content := xml.CreateElement("Content")
	rep_Content.CreateText(new_content)

	create_doc.Indent(2)
	res, _ := create_doc.WriteToString()
	return res
}

func ReplyImg(MsgFromUserName string, MsgToUserName string, media_id string) string {
	create_doc := etree.NewDocument()
	xml := create_doc.CreateElement("xml")
	//xml.CreateComment("These are all known people")

	ToUserName := xml.CreateElement("ToUserName")
	ToUserName.CreateText(MsgFromUserName)

	FromUserName := xml.CreateElement("FromUserName")
	FromUserName.CreateText(MsgToUserName)

	now := time.Now().Unix() //获取时间戳
	now_string := strconv.FormatInt(now, 10)
	CreateTime := xml.CreateElement("CreateTime")
	CreateTime.CreateText(now_string)

	MsgType := xml.CreateElement("MsgType")
	MsgType.CreateText("image")

	Image := xml.CreateElement("Image")
	MediaId := Image.CreateElement("MediaId")
	MediaId.CreateText(media_id)

	create_doc.Indent(2)
	res, _ := create_doc.WriteToString()
	return res
}

func ReplyVoice(MsgFromUserName string, MsgToUserName string, media_id string) string {
	create_doc := etree.NewDocument()
	xml := create_doc.CreateElement("xml")
	//xml.CreateComment("These are all known people")

	ToUserName := xml.CreateElement("ToUserName")
	ToUserName.CreateText(MsgFromUserName)

	FromUserName := xml.CreateElement("FromUserName")
	FromUserName.CreateText(MsgToUserName)

	now := time.Now().Unix() //获取时间戳
	now_string := strconv.FormatInt(now, 10)
	CreateTime := xml.CreateElement("CreateTime")
	CreateTime.CreateText(now_string)

	MsgType := xml.CreateElement("MsgType")
	MsgType.CreateText("voice")

	Voice := xml.CreateElement("Voice")
	MediaId := Voice.CreateElement("MediaId")
	MediaId.CreateText(media_id)
	create_doc.Indent(2)
	res, _ := create_doc.WriteToString()
	return res
}

func ReplyVideo(MsgFromUserName string, MsgToUserName string, media_id string, title string, description string) string {
	create_doc := etree.NewDocument()
	xml := create_doc.CreateElement("xml")
	//xml.CreateComment("These are all known people")

	ToUserName := xml.CreateElement("ToUserName")
	ToUserName.CreateText(MsgFromUserName)

	FromUserName := xml.CreateElement("FromUserName")
	FromUserName.CreateText(MsgToUserName)

	now := time.Now().Unix() //获取时间戳
	now_string := strconv.FormatInt(now, 10)
	CreateTime := xml.CreateElement("CreateTime")
	CreateTime.CreateText(now_string)

	MsgType := xml.CreateElement("MsgType")
	MsgType.CreateText("video")

	Video := xml.CreateElement("Video")
	MediaId := Video.CreateElement("MediaId")
	MediaId.CreateText(media_id)
	Title := Video.CreateElement("Title")
	Title.CreateText(title)
	Description := Video.CreateElement("Description")
	Description.CreateText(description)

	create_doc.Indent(2)
	res, _ := create_doc.WriteToString()
	return res
}

func ReplyNews(MsgFromUserName string, MsgToUserName string, title string, description string, picUrl string, url string) string {
	create_doc := etree.NewDocument()
	xml := create_doc.CreateElement("xml")
	//xml.CreateComment("These are all known people")

	ToUserName := xml.CreateElement("ToUserName")
	ToUserName.CreateText(MsgFromUserName)

	FromUserName := xml.CreateElement("FromUserName")
	FromUserName.CreateText(MsgToUserName)

	now := time.Now().Unix() //获取时间戳
	now_string := strconv.FormatInt(now, 10)
	CreateTime := xml.CreateElement("CreateTime")
	CreateTime.CreateText(now_string)

	MsgType := xml.CreateElement("MsgType")
	MsgType.CreateText("news")

	ArticleCount := xml.CreateElement("ArticleCount")
	ArticleCount.CreateText("1")

	Articles := xml.CreateElement("Articles")
	item := Articles.CreateElement("item")
	Title := item.CreateElement("Title")
	Title.CreateText(title)
	Description := item.CreateElement("Description")
	Description.CreateText(description)

	PicUrl := item.CreateElement("PicUrl")
	PicUrl.CreateText(picUrl)
	Url := item.CreateElement("Url")
	Url.CreateText(url)

	create_doc.Indent(2)
	res, _ := create_doc.WriteToString()
	return res
}

func ConMysql(dsn string) (*sql.DB, error) {
	Db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return Db, nil
}

func MysqlQuery(sqlWhere string, sqlOrder string, sqlLimit string, sqlTable string, sqlField string, sqlQuery string, strDsn string) ([]map[string]string, error) {
	//defer myTrace(sqlTable+"数据表查询开始")()
	//连接数据库
	//db, err := sql.Open(driverName, dataSourceName)
	db, err := ConMysql(strDsn)
	if err != nil {
		log.Println("连接数据库失败", err.Error())
		return nil, err
	}
	defer db.Close()
	var null string
	if sqlQuery == null {
		sqlQuery = "SELECT " + sqlField + " FROM " + sqlTable + " WHERE 1=1 " + sqlWhere + " " + sqlOrder + " Limit " + sqlLimit
	}
	//loger.Println(sqlQuery)
	//查询数据库
	query, err := db.Query(sqlQuery)
	if err != nil {
		log.Println("查询数据库失败", err.Error())
		return nil, err
	}
	defer query.Close()
	//读出查询出的列字段名
	cols, _ := query.Columns()
	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(cols))
	//query.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(cols))
	//让每一行数据都填充到[][]byte里面
	for i := range values {
		scans[i] = &values[i]
	}
	//最后得到的map
	// := make(map[int]map[string]string)
	results := [](map[string]string){}
	i := 0
	for query.Next() { //循环，让游标往下推
		if err := query.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			log.Println(err)
			return nil, err
		}
		row := make(map[string]string) //每行数据
		for k, v := range values {     //每行数据是放在values里面，现在把它挪到row里
			key := cols[k]
			row[key] = string(v)
		}
		results = append(results, row) //装入结果集中
		i++
	}
	db.Close() //用完关闭
	return results, nil
}

func DelTableData(dsn string, table_name string, whereSql string) (int64, error) {
	//defer fmt.Println(table_name + "表数据删除")()
	db, err := ConMysql(dsn)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	res, err := db.Exec("DELETE FROM " + table_name + " WHERE " + whereSql)
	if err != nil {
		return 0, err
	}
	intResult, err := res.RowsAffected()

	return intResult, nil
}

func InsertAll(strDsn string, table_name string, data []map[string]string) (int, error, string) {
	//defer MyTrace(table_name + "表数据插入")()
	//链接数据库
	sql_str_first := "INSERT INTO " + table_name + "(" //开始sql文
	sql_str_next := "("                                //后续sql文
	data_len := len(data)                              //整个数据长度
	if data_len == 0 {
		return data_len, nil, ""
	}
	maxitems := 300
	remainder := data_len % maxitems
	loops := data_len / maxitems

	frow := data[0]
	fnum := len(frow)
	var fields []string
	var counter int
	for fieldname, _ := range frow {
		counter++
		sql_str_first += fieldname
		sql_str_next += "?"
		if counter < fnum {
			sql_str_first += ","
			sql_str_next += ","
		}
		fields = append(fields, fieldname)
	}
	fmt.Println(fields)
	sql_str_first += ")"
	sql_str_next += ")"

	sql_loops := ""
	sql_remai := ""
	var rows []map[string]string
	vlen := maxitems * fnum
	values := make([]interface{}, remainder*fnum)
	fmt.Println("-----------111--------------")

	db, err := ConMysql(strDsn)
	fmt.Println("-----------222--------------")
	defer db.Close()
	//事务开始
	tx, err := db.Begin()
	if err != nil {
		return 0, err, ""
	}
	fmt.Println("-----------33--------------")

	if remainder > 0 {
		sql_remai += sql_str_first + "VALUES"
		for k := 1; k < remainder; k++ {
			sql_remai += sql_str_next + ","
		}
		sql_remai += sql_str_next
		rows = data[0:remainder]
		counter = 0
		for _, row := range rows {
			for _, fieldname := range fields {
				values[counter] = row[fieldname]
				counter++
			}
		}
		_, err := tx.Exec(sql_remai, values...)
		if nil != err {
			tx.Rollback()
			return 0, err, sql_remai + ";" + fmt.Sprintf("%s", values)
		}

	}
	fmt.Println("-----------666--------------")

	values = make([]interface{}, vlen, vlen)
	if loops > 0 {
		sql_loops += sql_str_first + "VALUES"
		for k := 1; k < maxitems; k++ {
			sql_loops += sql_str_next + ","
		}
		sql_loops += sql_str_next
		startpos := remainder
		endpos := remainder + maxitems
		//fmt.Println(sql_loops)
		stmtIns, err := tx.Prepare(sql_loops)
		if nil != err {
			return 0, err, sql_loops
		}
		defer stmtIns.Close()
		for k := 0; k < loops; k++ {
			rows = data[startpos:endpos]
			counter = 0
			for _, row := range rows {
				for _, fieldname := range fields {
					values[counter] = row[fieldname]
					counter++
				}
			}
			_, err := stmtIns.Exec(values...)
			if nil != err {
				tx.Rollback()
				return 0, err, sql_loops + ";" + fmt.Sprintf("%s", values)
			}
			startpos = endpos
			endpos += maxitems
		}
	}
	tx.Commit()
	return data_len, nil, ""
}

func ReplaceIntoAll(strDsn string, table_name string, data []map[string]string) (int, error, string) {
	//defer MyTrace(table_name + "表数据插入或更新")()
	//链接数据库
	sql_str_first := "REPLACE INTO " + table_name + "(" //开始sql文
	sql_str_next := "("                                 //后续sql文
	data_len := len(data)                               //整个数据长度
	if data_len == 0 {
		return data_len, nil, ""
	}
	maxitems := 300
	remainder := data_len % maxitems
	loops := data_len / maxitems

	frow := data[0]
	fnum := len(frow)
	var fields []string
	var counter int
	for fieldname, _ := range frow {
		counter++
		sql_str_first += fieldname
		sql_str_next += "?"
		if counter < fnum {
			sql_str_first += ","
			sql_str_next += ","
		}
		fields = append(fields, fieldname)
	}
	sql_str_first += ")"
	sql_str_next += ")"

	sql_loops := ""
	sql_remai := ""
	var rows []map[string]string
	vlen := maxitems * fnum
	values := make([]interface{}, remainder*fnum)
	db, err := ConMysql(strDsn)
	defer db.Close()
	//事务开始
	tx, err := db.Begin()
	if err != nil {
		return 0, err, ""
	}

	if remainder > 0 {
		sql_remai += sql_str_first + "VALUES"
		for k := 1; k < remainder; k++ {
			sql_remai += sql_str_next + ","
		}
		sql_remai += sql_str_next
		rows = data[0:remainder]
		counter = 0
		for _, row := range rows {
			for _, fieldname := range fields {
				values[counter] = row[fieldname]
				counter++
			}
		}
		_, err := tx.Exec(sql_remai, values...)
		if nil != err {
			tx.Rollback()
			return 0, err, sql_remai + ";" + fmt.Sprintf("%s", values)
		}

	}
	values = make([]interface{}, vlen, vlen)
	if loops > 0 {
		sql_loops += sql_str_first + "VALUES"
		for k := 1; k < maxitems; k++ {
			sql_loops += sql_str_next + ","
		}
		sql_loops += sql_str_next
		startpos := remainder
		endpos := remainder + maxitems
		//fmt.Println(sql_loops)
		stmtIns, err := tx.Prepare(sql_loops)
		if nil != err {
			return 0, err, sql_loops
		}
		defer stmtIns.Close()
		for k := 0; k < loops; k++ {
			rows = data[startpos:endpos]
			counter = 0
			for _, row := range rows {
				for _, fieldname := range fields {
					values[counter] = row[fieldname]
					counter++
				}
			}
			_, err := stmtIns.Exec(values...)
			if nil != err {
				tx.Rollback()
				return 0, err, sql_loops + ";" + fmt.Sprintf("%s", values)
			}
			startpos = endpos
			endpos += maxitems
		}
	}
	tx.Commit()
	return data_len, nil, ""
}

func GetIntStrRes(str string) (int, float64) {
	//var deci decimal.Decimal
	int, _ := strconv.Atoi(str)
	float, _ := strconv.ParseFloat(str, 64)
	return int, float
}

func WxGet(redis *redis.Client, key string) (string, error) {
	var result string
	result, Err = redis.Get(key).Result()
	if Err != nil || result == "" {
		return result, errors.New("未找到")
	}
	return result, nil
}

func WxSet(redis *redis.Client, key string, val interface{}, exp_time time.Duration) (string, error) {
	var result string
	result, Err = redis.Set(key, val, exp_time).Result()
	if Err != nil {
		return result, Err
	}

	return result, nil
}

func GetAccessToken(redis *redis.Client) (string, error) {
	get_accesstoken_uri := consts.GET_ACCESS_TOKEN_URI +
		"?grant_type=" + consts.GrantTypeToType(consts.GRANT_TYPE) +
		"&appid=" + consts.APP_ID +
		"&secret=" + consts.SECRET
	AccessToken, Err = WxGet(redis, consts.ACCESS_TOKEN)

	if Err != nil || AccessToken == "" {
		AccessToken_json, Err := tools.HttpGet(get_accesstoken_uri)

		AccessToken = gjson.Get(AccessToken_json, "access_token").String()

		_, Err = WxSet(redis, consts.ACCESS_TOKEN, AccessToken, time.Duration(7200)*time.Second)
		if Err != nil {
			return "", Err
		}
		return AccessToken, Err
	}
	return AccessToken, Err
}

func WxHMGet(redis *redis.Client, key string, fields ...string) ([]interface{}, error) {
	var result []interface{}
	result, Err = redis.HMGet(key, fields...).Result()
	if Err != nil {
		return nil, Err
	}
	for _, v := range result {
		if v == nil {
			return nil, errors.New("redis取值为nil")
		}
	}
	return result, nil
}
