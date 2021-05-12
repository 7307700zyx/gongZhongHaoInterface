package controllers

import (
	tools "gitlab.wsmfin.com/DEV/GoLangTools"
	"gongZhongHaoInterface/common"
	consts "gongZhongHaoInterface/conf"
	"net/http"
	"regexp"
	"strconv"
)

func WxActivityText(w http.ResponseWriter, open_id string, ggh_id string, content string) {
	msg_type_str := "【活动记录】-"
	today := tools.GetCurrentStringDateFormat("2006-01-02")
	find_user_activity_sql := "SELECT * FROM fa_customer_activity WHERE open_id='" + open_id + "' " +
		"AND date_format(upd_timestamp,'%Y-%m-%d') = '" + today + "'"
	result, _ := common.MysqlQuery("", "", "", "", "", find_user_activity_sql, consts.MYSQL_DSN)
	resCnt := len(result)
	//没有用户参与活动的数据
	if resCnt == 0 {
		//用户传来关键字 朋友圈
		if content == consts.ACTIVITY_PYQ {
			tools.OutPutInfo(nil, msg_type_str+" 参加【"+content+"】活动1")
			WxActivityPyqStep1(w, open_id, ggh_id, content)
			return
			//用户传来关键字 碎片兑奖
		}
		if content == consts.ACTIVITY_SPDJ {
			//TODO 碎片兑奖
			tools.OutPutInfo(nil, msg_type_str+" 参加【"+content+"】活动2")
			WxActivitySPDJStep1(w, open_id, ggh_id, content)
			return
		}

		//不是参与活动的关键字，也不在自动回复列表中
		//直接调用默认回复
		//tools.OutPutInfo(nil,msg_type_str+" 不是关键字，也没有加入活动，调用默认回复")
		xmlStr := common.ReplyTextCommon(open_id, ggh_id)
		w.Write([]byte(xmlStr))
		return

		//已经参加了活动，判断传来的文字信息所属什么活动-关键字
	} else {

		//用户传来关键字 朋友圈
		if content == consts.ACTIVITY_PYQ {
			tools.OutPutInfo(nil, msg_type_str+" 参加【"+content+"】活动1")
			WxActivityPyqStep1(w, open_id, ggh_id, content)
			return
			//用户传来关键字 碎片兑奖
		}
		if content == consts.ACTIVITY_SPDJ {
			tools.OutPutInfo(nil, msg_type_str+" 参加【"+content+"】活动2")
			WxActivitySPDJStep1(w, open_id, ggh_id, content)
			return
		}

		isInActivity := true
		for i := 0; i < resCnt; i++ {

			//判断【朋友圈】活动当前步骤数，判断需要恢复的信息
			if result[i]["act_name"] == consts.ACTIVITY_PYQ {
				//1：发送【朋友圈】
				//2：发送图片
				//3：发送【手机号】
				if result[i]["step"] == "1" {
					//用户发送了不属于关键字的信息过来，
					//并且当前活动步骤是1，【朋友圈】活动 的下一步 是需要发送朋友圈截图图片
					//直接调用默认回复
					//WxActivityPyqStep1(w,open_id,ggh_id,content)
					tools.OutPutInfo(nil, msg_type_str+"【"+consts.ACTIVITY_PYQ+"】当前进行中，步骤（1），调用默认回复")
					xmlStr := common.ReplyTextCommon(open_id, ggh_id)
					w.Write([]byte(xmlStr))
					return
				} else if result[i]["step"] == "2" {
					matchResult, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, content)
					if matchResult {
						//进入朋友圈活动手机回复步骤
						tools.OutPutInfo(nil, msg_type_str+"【"+consts.ACTIVITY_PYQ+"】当前进行中，步骤（2），并且发送的是手机号")
						WxActivityPyqStep3(w, open_id, ggh_id, content)
						return
					} else {
						tools.OutPutInfo(nil, msg_type_str+"【"+consts.ACTIVITY_PYQ+"】当前进行中，步骤（2），并且发送的不是手机号")
						//不是正常手机号类型
						//直接调用默认回复
						xmlStr := common.ReplyTextCommon(open_id, ggh_id)
						w.Write([]byte(xmlStr))
						return
					}
				} else {
					tools.OutPutInfo(nil, msg_type_str+"未匹配【"+consts.ACTIVITY_PYQ+"】活动")
					isInActivity = false
				}
			}
			if result[i]["act_name"] == consts.ACTIVITY_SPDJ {
				if content == consts.ACTIVITY_SPDJ {
					tools.OutPutInfo(nil, msg_type_str+"匹配【"+consts.ACTIVITY_SPDJ+"】活动，开始")
					WxActivitySPDJStep1(w, open_id, ggh_id, content)
					return
					//如果当前是step2，并且发送的 是【手机号】
				}
				if result[i]["step"] == "1" {
					//3：发送【手机号】
					matchResult, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, content)
					if matchResult {
						//tools.OutPutInfo(nil,msg_type_str+"【"+consts.ACTIVITY_SPDJ+"】当前进行中，步骤（1），并且发送的是手机号")

						//进入碎片兑奖 活动 手机回复步骤
						WxActivitySPDJStep2(w, open_id, ggh_id, content)
						return
					} else {
						tools.OutPutInfo(nil, msg_type_str+"【"+consts.ACTIVITY_SPDJ+"】当前进行中，步骤（1），发送的不是手机号")

						//不是正常手机号类型
						//直接调用默认回复
						xmlStr := common.ReplyTextCommon(open_id, ggh_id)
						w.Write([]byte(xmlStr))
						return
					}
				} else if result[i]["step"] == "2" {
					//tools.OutPutInfo(nil,msg_type_str+"【"+consts.ACTIVITY_SPDJ+"】当前进行中，步骤（2），判断是否在奖品列表中："+content)
					prizeCodeList, _, _ := common.GetSPDJDetail()
					for i := 0; i < len(prizeCodeList); i++ {
						if content == prizeCodeList[i] {
							//碎片兑奖
							WxActivitySPDJStep3(w, open_id, ggh_id, i, content)
							return
						}
					}
				} else {
					tools.OutPutInfo(nil, msg_type_str+"未匹配【"+consts.ACTIVITY_SPDJ+"】活动1")
					isInActivity = false
				}
			}
			isInActivity = false
		}
		if isInActivity == false {
			tools.OutPutInfo(nil, msg_type_str+"未匹配任意活动")

			//不是正常手机号类型
			//直接调用默认回复
			xmlStr := common.ReplyTextCommon(open_id, ggh_id)
			w.Write([]byte(xmlStr))
			return
		}
	}
}

func WxActivityImage(w http.ResponseWriter, open_id string, ggh_id string, picurl string, mediaid string) {
	msg_type_str := "【朋友圈活动-步骤2-图片】"
	today := tools.GetCurrentStringDateFormat("2006-01-02")
	find_user_activity_sql := "SELECT * FROM fa_customer_activity WHERE open_id='" + open_id + "' AND date_format(upd_timestamp,'%Y-%m-%d') = '" + today + "'" +
		" AND act_name = '" + consts.ACTIVITY_PYQ + "'"
	result, _ := common.MysqlQuery("", "", "", "", "", find_user_activity_sql, consts.MYSQL_DSN)
	resCnt := len(result)
	//没有用户参与活动的数据
	if resCnt == 0 {
		tools.OutPutInfo(nil, msg_type_str+"未找到用户当天参与【朋友圈】活动记录")

		//不是参与活动的关键字，也不在自动回复列表中
		//直接调用默认回复
		xmlStr := common.ReplyTextCommon(open_id, ggh_id)
		w.Write([]byte(xmlStr))
		return
		//}
		//已经参加了活动，判断传来的文字信息所属什么活动-关键字
	} else {
		//如果当前是step1，并且发送的 是【图片】
		if result[0]["step"] == "1" {
			//进入朋友圈活动【图片】识别步骤
			tools.OutPutInfo(nil, msg_type_str+"，调用图片识别接口")
			WxActivityPyqStep2(w, open_id, ggh_id, picurl, mediaid)
		} else {
			//不是接收图片截图的步骤
			//直接调用默认回复
			tools.OutPutInfo(nil, msg_type_str+"，当前步骤【1】不是接收图片截图的流程")
			xmlStr := common.ReplyTextCommon(open_id, ggh_id)
			w.Write([]byte(xmlStr))
			return
		}
	}
}

func isJoinPYQactivity(open_id string) bool {
	msg_type_str := "【朋友圈活动-判断是否已加入活动】"
	//判断当前用户，当天是否参加过活动
	today := tools.GetCurrentStringDateFormat("2006-01-02")

	find_user_activity_sql := "SELECT * FROM fa_customer_activity " +
		"WHERE open_id='" + open_id + "' AND date_format(upd_timestamp,'%Y-%m-%d') = '" + today + "'" +
		" AND act_name = '" + consts.ACTIVITY_PYQ + "'"

	result, _ := common.MysqlQuery("", "", "", "", "", find_user_activity_sql, consts.MYSQL_DSN)
	resCnt := len(result)
	//当天已经有数据，并且是完结状态，不可以继续此活动
	if resCnt > 0 {
		//已经完结状态了，不允许继续参加
		if result[0]["final_result"] == "2" {
			return true
		} else {
			//删除数据
			sqlToDelAct := "DELETE FROM fa_customer_activity WHERE act_name = '" + consts.ACTIVITY_PYQ + "' AND open_id='" + open_id + "'"
			_, err := common.MysqlExec(sqlToDelAct, consts.MYSQL_DSN)
			if err != nil {
				tools.OutPutInfo(nil, msg_type_str+"用户参与朋友圈活动，清除数据库数据失败："+sqlToDelAct)
			}
			return false
		}
	} else {
		tools.OutPutInfo(nil, msg_type_str+"，数据查询失败："+find_user_activity_sql)
		return false
	}
}

//发送关键字【朋友圈】，操作
func WxActivityPyqStep1(w http.ResponseWriter, open_id string, ggh_id string, content string) {
	msg_type_str := "【朋友圈】【step1】-"

	//tools.OutPutInfo(nil,msg_type_str+"用户发送“朋友圈”关键字，活动开始" )

	//step1、先判断当前(open_id) 是否当天已经参加过一次完整活动，并已经完结状态
	//step1-1，已经参加过，被动回复消息 “小可爱这个任务您已经兑换过碎片啦~”
	//			中断
	res := isJoinPYQactivity(open_id)
	if res == true {
		//tools.OutPutInfo(nil,msg_type_str+"用户发送“朋友圈”关键字，当天活动已经完成" )
		xmlStr := common.ReplyText(open_id, ggh_id, "小可爱这个任务您已经兑换过碎片啦~")
		w.Write([]byte(xmlStr))
		return
	}
	dataSlice := make([]map[string]string, 0) //初始化数据切片
	mapData := make(map[string]string)
	mapData["step"] = "1"
	mapData["act_name"] = consts.ACTIVITY_PYQ
	mapData["open_id"] = open_id

	dataSlice = append(dataSlice, mapData)
	_, err, strSqlTmp := common.InsertAll(consts.MYSQL_DSN, "fa_customer_activity", dataSlice)
	if err != nil {
		tools.OutPutInfo(nil, msg_type_str+"用户发送“朋友圈”关键字，插入数据库失败："+strSqlTmp)
		w.Write([]byte("请重新发送“朋友圈”关键字参与活动"))
		return
	}

	//step2、插入数据库一条数据，step=1，img_check=1，result=1
	//		  回复：“请发送朋友圈截图”
	//		  监听
	xmlStr := common.ReplyText(open_id, ggh_id, "请发送朋友圈截图~")
	w.Write([]byte(xmlStr))
	return
}

//发送朋友圈图片识别操作
func WxActivityPyqStep2(w http.ResponseWriter, open_id string, ggh_id string, img_path string, mediaid string) {
	//tools.Logger = tools.WriteLogFile("./act_pyq_blog", "")
	//step1、先判断当前(open_id) 是否当天已经参加过一次完整活动，并已经完结状态
	//step1-1，已经参加过，被动回复消息 “小可爱这个任务您已经兑换过碎片啦~”
	//			中断
	//fmt.Println(open_id)
	//fmt.Println(ggh_id)
	//fmt.Println(img_path)
	//fmt.Println(mediaid)
	//
	//fmt.Println("--------1---------------")
	//res := isJoinPYQactivity(open_id)
	//if(res == true){
	//	xmlStr := common.ReplyText(open_id, ggh_id, "小可爱这个任务您已经兑换过碎片啦~")
	//	w.Write([]byte(xmlStr))
	//	return
	//}
	//fmt.Println("--------2---------------")
	//step3-2：是图片，
	msg_type_str := "【朋友圈】【step2】-"
	//tools.OutPutInfo(nil,msg_type_str+"用户发送图片，调用图片识别方法")

	//（1）发送给图片识别脚本
	//发送给系统校验
	//不需要等待返回值，直接回复用户信息
	go common.VerImg(img_path, open_id)

	//（2）修改数据库步骤，step=2，img_path=xxx
	customer_img_msg_sql := " UPDATE fa_customer_activity SET step=2,img_check = 1,img_path='" + img_path + "',img_media_id = '" + mediaid + "' " +
		"WHERE act_name = '" + consts.ACTIVITY_PYQ + "' AND open_id = '" + open_id + "'"
	//tools.OutPutInfo(nil,msg_type_str+"更新用户参与活动状态信息")
	_, err := common.MysqlExec(customer_img_msg_sql, consts.MYSQL_DSN)
	if err != nil {
		tools.OutPutInfo(err, msg_type_str+"更新用户参与活动状态信息数据操作失败："+customer_img_msg_sql)
		xmlStr := common.ReplyText(open_id, ggh_id, "请发送朋友圈截图")
		w.Write([]byte(xmlStr))
		return
	}
	//（3）回复：“小编收到啦，请回复您的手机号。”
	xmlStr := common.ReplyText(open_id, ggh_id, "小编收到啦，请回复您的手机号。")
	w.Write([]byte(xmlStr))
	return
}

//发送【手机号】，操作
func WxActivityPyqStep3(w http.ResponseWriter, open_id string, ggh_id string, content string) {
	msg_type_str := "【朋友圈】【step3】-"
	//tools.OutPutInfo(nil,msg_type_str+"用户发送手机号后续处理开始")

	//判断用户表，对应的手机号和openid 的用户是否存在，不存在则创建
	//tools.OutPutInfo(nil,msg_type_str+"判断用户表，对应的手机号和openid 的用户是否存在，不存在则创建")
	boo := common.AddCustomer(open_id, content)
	if boo == false {
		//tools.OutPutInfo("当前手机号已被占用，请更换手机号" )
		tools.OutPutInfo(nil, msg_type_str+"当前手机号已被占用，请更换手机号")
		xmlStr := common.ReplyText(open_id, ggh_id, "当前手机号已被占用，请更换手机号")
		w.Write([]byte(xmlStr))
		return
	}
	////	step4-3：修改数据库状态，step=3
	//用户活动表：修改用户手机号
	customer_phone_num_sql := " UPDATE fa_customer_activity SET step=3,phone_num = '" + content + "'" +
		" WHERE act_name = '" + consts.ACTIVITY_PYQ + "' AND open_id = '" + open_id + "'"
	//tools.OutPutInfo(nil,msg_type_str+"更新用户参与活动状态信息")

	_, err := common.MysqlExec(customer_phone_num_sql, consts.MYSQL_DSN)
	if err != nil {
		tools.OutPutInfo(nil, msg_type_str+"更新用户参与活动状态信息，失败："+customer_phone_num_sql)
		xmlStr := common.ReplyText(open_id, ggh_id, "请发送手机号")
		w.Write([]byte(xmlStr))
		return
	}

	result := common.GetActivity(consts.ACTIVITY_PYQ, open_id)

	//	step4-4：判断图片校验状态img_check
	//1=》校验中；2=》校验成功；3=》机器审核失败；4=》人工审核失败
	switch result["img_check"] {
	//		step4-4-1：未校验中
	case "1":
		xmlStr := common.ReplyText(open_id, ggh_id, "正在审核图片中~")
		w.Write([]byte(xmlStr))
		return
	//校验通过，为用户增加积分，并修改当天任务完成状态
	case "2":
		msg, res := common.ActivityPyqAddScore(open_id, content)
		if res == false {
			xmlStr := common.ReplyText(open_id, ggh_id, msg)
			w.Write([]byte(xmlStr))
		} else {
			xmlStr := common.ReplyText(open_id, ggh_id, "您本次获得5个碎片，回复“碎片兑奖”即可兑换奖品哦。")
			w.Write([]byte(xmlStr))
		}
		return
	//机器审核图片失败，交给人工审核
	case "3":
		xmlStr := common.ReplyText(open_id, ggh_id, "朋友圈活动截图识别失败，正在转人工审核中~")
		w.Write([]byte(xmlStr))
		return
	//人工审核失败，图片完全有问题
	case "4":
		xmlStr := common.ReplyText(open_id, ggh_id, "人工审核朋友圈活动截图失败，请重新发送“朋友圈”关键字参与此次活动~")
		w.Write([]byte(xmlStr))
		return
	}
}

//为指定用户增加积分
//当天已经有数据，并且是完结状态，不可以继续此活动
//如果当前步骤数不是第三步，并且不是图片验证成功状态 ，不增加积分
func WxActivityPYQAddScore(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	open_id := query.Get("open_id")
	phone_num := query.Get("phone_num")

	_, boo := common.ActivityPyqAddScore(open_id, phone_num)
	if boo == true {
		WxKfSend("text", "您本次获得5个碎片，回复“碎片兑奖”即可兑换奖品哦。", open_id, consts.FWH_OPENID)
		return
	}
}

func delSPDJactivity(open_id string) {
	//删除数据
	sqlToDelAct := "DELETE FROM fa_customer_activity WHERE act_name = '" + consts.ACTIVITY_SPDJ + "' AND open_id='" + open_id + "'"
	_, err := common.MysqlExec(sqlToDelAct, consts.MYSQL_DSN)
	if err != nil {
		tools.OutPutInfo(nil, "用户参与碎片兑奖活动，清除数据库数据失败："+sqlToDelAct)
	}
}

func WxActivitySPDJStep1(w http.ResponseWriter, open_id string, ggh_id string, content string) {
	msg_type_str := "【碎片兑奖】【step1】-"
	//tools.OutPutInfo(nil,msg_type_str+"活动开始" )

	delSPDJactivity(open_id)

	dataSlice := make([]map[string]string, 0) //初始化数据切片
	mapData := make(map[string]string)
	mapData["step"] = "1"
	mapData["act_name"] = consts.ACTIVITY_SPDJ
	mapData["open_id"] = open_id

	dataSlice = append(dataSlice, mapData)
	//step2、插入数据库一条数据，step=1，img_check=1，result=1
	_, err, strSqlTmp := common.InsertAll(consts.MYSQL_DSN, "fa_customer_activity", dataSlice)
	if err != nil {
		tools.OutPutInfo(nil, msg_type_str+"用户参与碎片兑奖活动，插入数据库失败："+strSqlTmp)
		w.Write([]byte("异常，请重新发送“碎片兑奖”参与活动"))
		return
	}
	//回复：“请回复您的手机号”
	xmlStr := common.ReplyText(open_id, ggh_id, "请回复您的手机号~")
	w.Write([]byte(xmlStr))
	return
}

func WxActivitySPDJStep2(w http.ResponseWriter, open_id string, ggh_id string, content string) {
	msg_type_str := "【碎片兑奖】【step2】-"
	//创建用户
	common.AddCustomer(open_id, content)

	////	step4-3：修改数据库状态，step=2
	//用户活动表：修改用户手机号
	customer_phone_num_sql := " UPDATE fa_customer_activity SET step=2,phone_num = '" + content + "'" +
		" WHERE act_name = '" + consts.ACTIVITY_SPDJ + "' AND open_id = '" + open_id + "'"
	_, err := common.MysqlExec(customer_phone_num_sql, consts.MYSQL_DSN)
	if err != nil {
		tools.OutPutInfo(nil, msg_type_str+"用户参与碎片兑奖活动，修改用户手机号，当前步骤 失败："+customer_phone_num_sql)
		xmlStr := common.ReplyText(open_id, ggh_id, "系统异常，请重新发送手机号")
		w.Write([]byte(xmlStr))
		return
	}

	result, boo := common.GetCustomer(open_id, content)
	if boo == false {
		res := common.AddCustomer(open_id, content)
		if res == false {
			tools.OutPutInfo(nil, msg_type_str+"当前用户手机号，微信号已经占用，无法创建用户")
			xmlStr := common.ReplyText(open_id, ggh_id, "此手机号已经被占用了，无法增加积分")
			w.Write([]byte(xmlStr))
			return
		}

		result, boo = common.GetCustomer(open_id, content)
	}
	_, prizeList, _ := common.GetSPDJDetail()

	if result["score"] == "" {
		result["score"] = "0"
	}
	repluMsg := "您当前碎片数量为：" + result["score"] + "碎片，兑奖的清单如下：\n" + prizeList
	xmlStr := common.ReplyText(open_id, ggh_id, repluMsg)
	w.Write([]byte(xmlStr))
	return

}

//根据传递的兑奖编码,发送奖品秘钥
func WxActivitySPDJStep3(w http.ResponseWriter, open_id string, ggh_id string, i int, content string) {
	msg_type_str := "【碎片兑奖】【step3】-"

	_, _, prizeDetailMap := common.GetSPDJDetail()
	prizeName := prizeDetailMap[i][content]["name"]
	prizeCost := prizeDetailMap[i][content]["cost"]

	res := common.GetActivity(consts.ACTIVITY_SPDJ, open_id)
	result, boo := common.GetCustomer(open_id, res["phone_num"])
	if boo == false {
		tools.OutPutInfo(nil, msg_type_str+"用户参与碎片兑奖活动，未找到用户信息：")
		xmlStr := common.ReplyText(open_id, ggh_id, "未找到活动信息，请回复“碎片兑奖”重新开始")
		w.Write([]byte(xmlStr))
		return
	}
	score, _ := strconv.Atoi(result["score"])
	prizeCostInt, _ := strconv.Atoi(prizeCost)
	if score < prizeCostInt {
		//fmt.Println("碎片兑奖活动，当前用户积分："+result["score"]+",兑换物品："+prizeName+",物品所需积分"+prizeCost+",无法兑换")
		//tools.OutPutInfo(nil,msg_type_str+"碎片兑奖活动，当前用户积分："+result["score"]+",兑换物品："+prizeName+",物品所需积分"+prizeCost+",无法兑换")
		xmlStr := common.ReplyText(open_id, ggh_id, "碎片兑奖活动，当前用户积分："+result["score"]+",兑换物品："+prizeName+",物品所需积分"+prizeCost+",无法兑换")
		w.Write([]byte(xmlStr))
		return
	}

	left_score := score - prizeCostInt

	//获取兑换奖品秘钥，链接
	keys, boo := common.GetSecretKeys(content)
	if boo == false {
		tools.OutPutInfo(nil, msg_type_str+"获取兑奖秘钥 失败")
		xmlStr := common.ReplyText(open_id, ggh_id, "系统异常，请重新发兑奖物品编号")
		w.Write([]byte(xmlStr))
		return
	}

	//插入用户兑奖历史表
	tools.OutPutInfo(nil, msg_type_str+"插入用户兑奖历史表")
	add_customer_activity_exchange_sql := "INSERT INTO fa_customer_activity_exchange " +
		"(act_name,open_id,phone_num,total_score,prize_score,prize_code,prize_name,secret_key,link) VALUES" +
		"('" + consts.ACTIVITY_SPDJ + "','" + open_id + "','" + res["phone_num"] + "','" + result["score"] + "','" + prizeCost + "', " +
		"'" + content + "', '" + prizeName + "', '" + keys["secret_key"] + "', '" + keys["link"] + "' )"
	//tools.OutPutInfo(nil,msg_type_str+"插入用户兑奖历史表,sql:"+add_customer_activity_exchange_sql)
	_, err := common.MysqlExec(add_customer_activity_exchange_sql, consts.MYSQL_DSN)
	if err != nil {
		tools.OutPutInfo(nil, msg_type_str+"插入用户兑奖历史表 失败："+add_customer_activity_exchange_sql)
		xmlStr := common.ReplyText(open_id, ggh_id, "系统异常，请重新发兑奖物品编号")
		w.Write([]byte(xmlStr))
		return
	}

	//用户表：修改用户积分
	//tools.OutPutInfo(nil,msg_type_str+"修改用户积分")
	customer_score_sql := " UPDATE fa_customer SET score = score-" + prizeCost +
		" WHERE phone = '" + res["phone_num"] + "' AND openid = '" + open_id + "'"
	//tools.OutPutInfo(nil,msg_type_str+"修改用户积分,sql:"+customer_score_sql)
	_, err = common.MysqlExec(customer_score_sql, consts.MYSQL_DSN)
	if err != nil {
		tools.OutPutInfo(nil, msg_type_str+"修改用户积分失败："+customer_score_sql)
		xmlStr := common.ReplyText(open_id, ggh_id, "系统异常，请重新发兑奖物品编号")
		w.Write([]byte(xmlStr))
		return
	}
	//顾客：001
	//公：您兑换的奖品为 全店通用免单券
	//公：已扣除对应碎片200个 您的碎片数量为xxx个（奖品对应碎片的扣除）
	//公：兑换密钥：xxxxxxx（每次发送不同密钥）
	repluMsg := "您兑换的奖品为：" + prizeName + "\n" +
		"已扣除对应碎片" + prizeCost + "个，您的剩余碎片数量为" + strconv.Itoa(left_score) + "个\n" +
		"兑换密钥：" + keys["secret_key"] + "\n" +
		"兑奖链接：" + keys["link"]
	xmlStr := common.ReplyText(open_id, ggh_id, repluMsg)
	w.Write([]byte(xmlStr))
	return
}
