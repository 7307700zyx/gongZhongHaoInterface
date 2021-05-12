package conf

//mysql数据库和数据库操作相关配置
const MYSQL_DSN string = "zyx:123456abx@tcp(120.27.94.3:3306)/xykj?charset=utf8"

//strAquDSN: 			   dev:dev1qaz#EDC@tcp(10.1.49.37:3306)/bank_credit_dev?charset=utf8
const PAGE_COUNT int = 20
const PAGE_OFFSET int = 0

const ACCESS_TOKEN string = "access_token" //redis保存key


const REDIS_KEY_MEDIS_VOICE string = "voice_count"
const REDIS_KEY_MEDIS_VIDEO string = "video_count"
const REDIS_KEY_MEDIS_IMAGE string = "image_count"
const REDIS_KEY_MEDIS_NEWS string = "news_count"


//进程API相关

const WECHAT_PATH string = `.\`         //守护进程执行文件放置路径
const WECHAT_LOG string = "weichat.log" //守护进程log日志文件放置路径


//微信配置相关（测试）
const WECHAT_PORT string = ":8088"      //工作端端口号
const FWH_OPENID string = " gh_819c9ce0b8cf" //服务号的微信号(测试)
const APP_ID string = "wx461ebe61665c0054" //测试，app_id
const SECRET  string = "900648ce0c0da8fa37960b8e8d2a3825" //测试第三方用户唯一凭证密钥，即appsecret
const TOKEN  string = "fizG8qUYget4aMafWcI6UlGXrPEyC68C"
const REDIS_PWD string = "123456"
const REDIS_PROT string = "6379"
const REDIS_HOST_KEY string = "redisHost" //工作端配置文件内rediskey
const REDIS_HOST string = "10.1.49.37"    //中心端redis地址
const REDIS_CONN_TIMES int = 10           //redis尝试重连次数


//微信配置相关（生产）
//const WECHAT_PORT string = ":80"      //工作端端口号
//const FWH_OPENID string = "gh_0c7aab3fb892" //服务号的微信号
//const APP_ID string = "wx3f95cdc15cdeaf44"               //app_id
//const SECRET string = "3c64871728e063bd691c075c04656567" //第三方用户唯一凭证密钥，即appsecret
//const TOKEN string = "fizG8qUYget4aMafWcI6UlGXrPEyC68C"
//const REDIS_PWD string = "xykj1234"
//const REDIS_PROT string = "6380"
//const REDIS_HOST_KEY string = "redisHost" //工作端配置文件内rediskey
//const REDIS_HOST string = "127.0.0.1"    //中心端redis地址
//const REDIS_CONN_TIMES int = 10           //redis尝试重连次数


const GRANT_TYPE int = 1 //获取access_token填写client_credential

//微信接口
//获取token
//https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
const GET_ACCESS_TOKEN_URI string = "https://api.weixin.qq.com/cgi-bin/token"

//自动回复
//https://api.weixin.qq.com/cgi-bin/get_current_autoreply_info?access_token=ACCESS_TOKEN
//const AUTO_REPLY_URI string = "https://api.weixin.qq.com/cgi-bin/get_current_autoreply_info?access_token=ACCESS_TOKEN"

//获取永久素材列表 "https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=ACCESS_TOKEN"
const GET_MEDIA_LIST_URI string = "https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token="

//根据media_id获取永久素材 https://api.weixin.qq.com/cgi-bin/material/get_material?access_token=ACCESS_TOKEN
const GET_MEDIA_BY_MEDIAID_URI string = "https://api.weixin.qq.com/cgi-bin/material/get_material?access_token="

//根据media_id获取临时素材 https://api.weixin.qq.com/cgi-bin/media/get?access_token=ACCESS_TOKEN&media_id=MEDIA_ID
const GET_TMP_MEDIA_BY_MEDIAID_URI string = "https://api.weixin.qq.com/cgi-bin/media/get"

//GET https://api.weixin.qq.com/cgi-bin/material/get_materialcount?access_token=ACCESS_TOKEN
const GET_MEDIA_COUNT_URI string = "https://api.weixin.qq.com/cgi-bin/material/get_materialcount?access_token="

//添加客服帐号 POST https://api.weixin.qq.com/customservice/kfaccount/add?access_token=ACCESS_TOKEN
const ADD_KF_ACCOUNT_URI string = "https://api.weixin.qq.com/customservice/kfaccount/add?access_token="

//客服发送消息 POST https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=ACCESS_TOKEN
const KF_SEND_URI string = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="



// 各种消息类型,除了扫带二维码事件
/**
 * 文本消息
 */
const MESSAGE_TEXT string = "text"

/**
 * 图片消息
 */
const MESSAGE_IMAGE string = "image"

/**
 * 图文消息
 */
const MESSAGE_NEWS string = "news"

/**
 * 语音消息
 */
const MESSAGE_VOICE string = "voice"

/**
 * 视频消息
 */
const MESSAGE_VIDEO string = "video"

/**
 * 小视频消息
 */
const MESSAGE_SHORTVIDEO string = "shortvideo"

/**
 * 地理位置消息
 */
const MESSAGE_LOCATION string = "location"

/**
 * 链接消息
 */
const MESSAGE_LINK string = "link"

/**
 * 事件推送消息
 */
const MESSAGE_EVENT string = "event"

/**
 * 事件推送消息中,事件类型，subscribe(订阅)
 */
const MESSAGE_EVENT_SUBSCRIBE string = "subscribe"

/**
 * 事件推送消息中,事件类型，unsubscribe(取消订阅)
 */
const MESSAGE_EVENT_UNSUBSCRIBE string = "unsubscribe"

/**
 * 事件推送消息中,上报地理位置事件
 */
const MESSAGE_EVENT_LOCATION_UP string = "LOCATION"

/**
 * 事件推送消息中,自定义菜单事件,点击菜单拉取消息时的事件推送
 */
const MESSAGE_EVENT_CLICK string = "CLICK"

/**
 * 事件推送消息中,自定义菜单事件,点击菜单跳转链接时的事件推送
 */
const MESSAGE_EVENT_VIEW string = "VIEW"

//版本号
const AUTOLR_VERSION = "v1.0"

func GrantTypeToType(grant_type int) string {
	switch grant_type {
	//获取access_token填写client_credential
	case GRANT_TYPE:
		return "client_credential"
	default:
		return ""
	}
}
