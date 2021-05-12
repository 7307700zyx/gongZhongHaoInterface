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
	"strconv"
	"time"
)

var ConRedis *redis.Client

//设置全局的error句柄
var Err error
//var MyLogger tools.WriteLogFile("./wx_blog", "")
var AccessToken string

func GetSPDJDetail()([]string,string,[]map[string]map[string]string){
	dataSlice := make([]map[string]map[string]string, 0) //初始化数据切片

	prizeMap1 := make(map[string]map[string]string)
	prizeMap1["001"] = map[string]string{
		"name": "全店通用免单券",
		"cost": "200",
	}

	prizeMap2 := make(map[string]map[string]string)
	prizeMap2["002"] = map[string]string{
		"name": "本周推荐主题免单券",
		"cost": "150",
	}
	prizeMap3 := make(map[string]map[string]string)
	prizeMap3["003"] = map[string]string{
		"name": "铜币兑换券",
		"cost": "120",
	}
	prizeMap4 := make(map[string]map[string]string)
	prizeMap4["004"] = map[string]string{
		"name": "全店通用减免20元券",
		"cost": "100",
	}
	prizeMap5 := make(map[string]map[string]string)
	prizeMap5["005"] = map[string]string{
		"name": "本周推荐主题减免20元券",
		"cost": "50",
	}
	prizeMap6 := make(map[string]map[string]string)
	prizeMap6["006"] = map[string]string{
		"name": "饮品券",
		"cost": "30",
	}
	prizeMap7 := make(map[string]map[string]string)
	prizeMap7["007"] = map[string]string{
		"name": "棒棒糖券",
		"cost": "20",
	}
	dataSlice = append(dataSlice, prizeMap1)
	dataSlice = append(dataSlice, prizeMap2)
	dataSlice = append(dataSlice, prizeMap3)
	dataSlice = append(dataSlice, prizeMap4)
	dataSlice = append(dataSlice, prizeMap5)
	dataSlice = append(dataSlice, prizeMap6)
	dataSlice = append(dataSlice, prizeMap7)

	prizeDetailStr := ""
	prizeCodeList := make([]string, 0, len(dataSlice))

	for i := 0; i < len(dataSlice); i++ {
		for k,v := range dataSlice[i] {
			prizeCodeList = append(prizeCodeList, k)
			prizeDetailStr = prizeDetailStr + k+"："+v["name"]+" "+v["cost"]+"碎片\n"
		}
	}

	return prizeCodeList,prizeDetailStr,dataSlice
}
func ConnectRedis() (*redis.Client, error) {
	redisHost := consts.REDIS_HOST
	for j := 1; j <= consts.REDIS_CONN_TIMES; j++ {
		ConRedis, Err = tools.Conredis(redisHost+":"+consts.REDIS_PROT, consts.REDIS_PWD, "0")
		if Err != nil {
			tools.OutPutInfo(Err, "第", j, "次连接redis【异常】host:"+redisHost, "密码md5后:", tools.GetMd5Sum(consts.REDIS_PWD))
			tools.Sleep(1000)
		} else {
			tools.OutPutInfo(nil, "第", j, "次连接redis【正常】host:"+redisHost)
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

func ReplyTextCommon(MsgFromUserName string, MsgToUserName string) string {
	res := ReplyText(MsgFromUserName, MsgToUserName, consts.REPLY_TEXT_COMMON)
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
		tools.OutPutInfo(nil,"连接数据库失败", err.Error())
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
		tools.OutPutInfo(nil,"查询数据库失败", err.Error())
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
			tools.OutPutInfo(err)
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

func MysqlExec(sql string, strDsn string) (int64, error) {
	db, err := ConMysql(strDsn)
	if err != nil {
		return 0, err
	}
	response, err := db.Exec(sql)
	if err != nil {
		return 0, err
	}
	num, err := response.RowsAffected()

	return num, err
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

func GetActivity(act_name string,open_id string) map[string]string {
	find_user_activity_sql := "SELECT * FROM fa_customer_activity " +
		"WHERE open_id='"+open_id+"' AND act_name = '"+act_name+"'"
	result, _ := MysqlQuery("", "", "", "", "", find_user_activity_sql, consts.MYSQL_DSN)
	resCnt := len(result)
	if resCnt == 0{
		return map[string]string{}
	}
	return result[0]
}

func GetCustomer(open_id string,phone_num string) (map[string]string,bool) {
	find_user_sql := "SELECT * FROM fa_customer " +
		"WHERE openid='"+open_id+"' AND phone = '"+phone_num+"'"
	fmt.Println(find_user_sql)
	result, _ := MysqlQuery("", "", "", "", "", find_user_sql, consts.MYSQL_DSN)
	resCnt := len(result)
	if resCnt == 0{
		return map[string]string{},false
	}
	return result[0] , true
}

func GetSecretKeys() (map[string]string,bool) {
	find_secret_keys_sql := "SELECT * FROM fa_secret_keys " +
		"WHERE status=1"
	result, _ := MysqlQuery("", "", "", "", "", find_secret_keys_sql, consts.MYSQL_DSN)
	resCnt := len(result)
	if resCnt == 0{
		return map[string]string{},false
	}

	//用户表：修改用户积分
	set_secret_keys_sql := " UPDATE fa_secret_keys SET status = 2 " +"WHERE id = "+ result[0]["id"]
	_, err := MysqlExec(set_secret_keys_sql, consts.MYSQL_DSN)
	if err != nil{
		return map[string]string{},false
	}
	return result[0],true
}

func AddCustomer(open_id string,phone_num string) bool{
	find_customer_sql := "SELECT * FROM fa_customer " +
		"WHERE phone = '"+phone_num+"'"
	fmt.Println(find_customer_sql)
	result, _ := MysqlQuery("", "", "", "", "", find_customer_sql, consts.MYSQL_DSN)
	resCnt := len(result)
	if resCnt == 0{
		add_customer_sql := "INSERT fa_customer (openid,phone,score) VALUES ('"+open_id+"','"+phone_num+"',0)"
		_, err := MysqlExec(add_customer_sql, consts.MYSQL_DSN)
		if err != nil{
			tools.OutPutInfo(nil,"创建新用户失败：" + add_customer_sql)
			return false
		}
	}
	return true
}

func ActivityPyqAddScore(open_id string,content string)(string,bool){
	//活动状态校验，并判断是否是已经活动结束状态
	//today := tools.GetCurrentStringDateFormat(tools.GetCurrentStringDate())
	today := tools.GetCurrentStringDateFormat("2006-01-02")
	find_user_activity_sql := "SELECT * FROM fa_customer_activity " +
		"WHERE open_id='"+open_id+"' AND phone_num = '"+content+"'  AND date_format(upd_timestamp,'%Y-%m-%d') = '"+today +"'"+
		" AND act_name = '"+consts.ACTIVITY_PYQ+"'"
	fmt.Println(find_user_activity_sql)
	result, _ := MysqlQuery("", "", "", "", "", find_user_activity_sql, consts.MYSQL_DSN)
	resCnt := len(result)
	//当天已经有数据，并且是完结状态，不可以继续此活动
	if resCnt > 0 {
		//已经完结状态了，不允许继续参加
		if result[0]["final_result"] == "2" {
			tools.OutPutInfo(nil,"用户当天活动已经完成，并且分数已经增加，无法重复增加")
			return  "小可爱这个任务您已经兑换过碎片啦~",false
		}
	}
	if result[0]["step"] == "3" && result[0]["img_check"] == "2"{
		//		step4-4-3：校验通过
		//				   修改数据	(1)用户表：增加积分
		//	step4-3：修改数据库状态，step=3
		customer_add_score_sql := " UPDATE fa_customer SET score=score+"+consts.PYQ_SCORE+
			" WHERE openid = '"+open_id+"' AND phone = '"+content+"'"
		fmt.Println(customer_add_score_sql)
		_, err := MysqlExec(customer_add_score_sql, consts.MYSQL_DSN)
		if err != nil{
			tools.OutPutInfo(nil,"为用户增加积分失败：" + customer_add_score_sql)
			return "系统错误，请重新发送手机号",false
		}
		//				   			（2）用户活动表：final_result=2，score=xx
		customer_activity_step_sql := " UPDATE fa_customer_activity SET step=3,score="+consts.PYQ_SCORE+",final_result = 2"+
			" WHERE act_name = '"+consts.ACTIVITY_PYQ +"' AND open_id = '"+open_id+"'"
		fmt.Println(customer_activity_step_sql)
		_, err = MysqlExec(customer_activity_step_sql, consts.MYSQL_DSN)
		if err != nil{
			tools.OutPutInfo(nil,"用户参与朋友圈活动，为用户增加积分 失败：" + customer_activity_step_sql)
			return "系统错误，请重新发送手机号",false
		}
		return "success",true
	}else{
		tools.OutPutInfo(nil,"用户当天活动图片验证与手机号步骤未完成，无法增加积分")
		return "系统错误，请重新发送手机号", false
	}
}

//朋友圈图片截图识别
//TODO
func VerImg(imgUrl,openid string)error{
	return nil
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
	if Err != nil || result == "" || len(result) == 0 {
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
	tools.OutPutInfo(nil,"====================读取redis内的token-开始============================")
	tools.OutPutInfo(nil,"请求地址为："+get_accesstoken_uri)
	AccessToken, Err = WxGet(redis, consts.ACCESS_TOKEN)
	tools.OutPutInfo(nil,"token取值为："+AccessToken)
	tools.OutPutInfo(nil,"====================读取redis内的token-结束============================")

	if AccessToken == "" || len(AccessToken)==0  {
		AccessToken_json, Err := tools.HttpGet(get_accesstoken_uri)
		fmt.Println(AccessToken_json)
		tools.OutPutInfo(nil,"请求url返回值为："+AccessToken)

		AccessToken = gjson.Get(AccessToken_json, "access_token").String()
		fmt.Println(AccessToken)
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
