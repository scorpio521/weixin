package controllers
import (
	"learnbeego/blog/models/class"
    "github.com/astaxie/beego"
    "strings"
    "sort"
    "crypto/sha1"
    "io"
    "fmt"
	_"encoding/json"
	"encoding/xml"
    "time"
)
type WXController struct {
	beego.Controller
}

type JSONStruct struct {
	Code int  `json:"code"`
	Msg  string  `json:"msg"`
}
//`这个代表前端传来的参数格式是这样子的`
type XMLStruct struct {
	XMLName xml.Name `xml:"xml"`
	ToUserName string  `xml:"ToUserName"`
	FromUserName string  `xml:"FromUserName"`
	CreateTime string  `xml:"CreateTime"`
	MsgType string  `xml:"MsgType"`
	Event string  `xml:"Event"`
	ExpiredTime string  `xml:"ExpiredTime"`
}

type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}
type TextResponseBody struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}

const Token = "weixin"
func (c *WXController) Get(){
	timestamp,nonce,signatureIn := c.GetString("timestamp"),c.GetString("nonce"),c.GetString("signature")
	//s :="nonce str is:"
	//joinstr :=strings.Join([]string{s,nonce},"")
	//beego.Info(joinstr)
	signatureGen := makeSignature(timestamp, nonce)
	beego.Info("xxxxxxxxxxxx")
	//将加密后获得的字符串与signature对比，如果一致，说明该请求来源于微信
	if signatureGen != signatureIn {
		fmt.Printf("signatureGen != signatureIn signatureGen=%s,signatureIn=%s\n", signatureGen, signatureIn)
		c.Ctx.WriteString("")

	} else {
		//如果请求来自于微信，则原样返回echostr参数内容 以上完成后，接入验证就会生效，开发者配置提交就会成功。
		echostr := c.GetString("echostr")
		c.Ctx.WriteString(echostr)
	}
}
func makeSignature(timestamp, nonce string) string {

	//1. 将 plat_token、timestamp、nonce三个参数进行字典序排序
	sl := []string{Token, timestamp, nonce}
	sort.Strings(sl)
	//2. 将三个参数字符串拼接成一个字符串进行sha1加密
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))

	return fmt.Sprintf("%x", s.Sum(nil))
}

/**
//1
s11 := new(Persons)
s11.Name = "cj"
s11.getName()
//2
s0 := &Persons{}
s0.Name = "cj"
s0.getName()
//3
s0 := &Persons{Name:"cj"}
s0.getName()
 */

func (c *WXController) Post() {
	fmt.Println("config",beego.AppConfig.String("log_path"))
	//解析json demo
	//	var class JSONStruct
	//	var err error
	//	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &class); err == nil {
	//		c.Data["json"] = &class
	//		fmt.Println(class.Code)
	//		//c.Data["json"] = &JSONStruct{007, "hello"}
	//	} else {
	//		fmt.Println(err)
	//		c.Data["json"] = err
	//	}
	//	c.ServeJSON()
	// 解析xml demo
	//var class XMLStruct
	//beego.Info(beego.AppConfig.String("log_path"))
	//err := xml.Unmarshal(c.Ctx.Input.RequestBody, &class)//好方便，就这样就解析了xml，so关键是xml的结构
	//if err != nil {
	//	fmt.Println("xml.Unmarshal is err:", err.Error())
	//	c.Data["json"] = err
	//}else{
	//	c.Data["json"] = &class
	//	fmt.Println(class.ExpiredTime)
	//	beego.Info(class.ExpiredTime)
	//}
	//c.ServeJSON()

	var classtext TextRequestBody
	beego.Info(beego.AppConfig.String("log_path"))
	err := xml.Unmarshal(c.Ctx.Input.RequestBody, &classtext)//好方便，就这样就解析了xml，so关键是xml的结构
	if err != nil {
		fmt.Println("xml.Unmarshal is err:", err.Error())
		c.Data["json"] = err
	}else{
		customerpath := beego.AppConfig.String("customerpath")
		if beego.AppConfig.String("runmode")=="dev" {
			customerpath = beego.AppConfig.String("customerpathdev")
		}
		con := &class.Config{Customerpath:customerpath}
		customermap, err := con.QueryByPath()
		//fmt.Println(customermap, err)
		//returnparmas:[{1 dev4 12121 1212121}]
		//returnparmas: [] code = 1 ; msg = the sql sentence is wrong,please check
		beego.Info(customermap,err)
		appconfig := make(map[string]string)
		for _, v := range customermap {
			appconfig["appid"] = v.Appid
			appconfig["appsec"] = v.Appsec
		}
		fmt.Println(appconfig["appid"])



		res := &TextResponseBody{}
		res.MsgType = classtext.MsgType
		res.ToUserName = classtext.ToUserName
		res.Content = classtext.Content
		res.CreateTime = classtext.CreateTime
		res.MsgId = classtext.MsgId
		c.Data["json"] = &res
		beego.Info(res)
		fmt.Println(res.MsgType)
		beego.Info(res.MsgType)
	}
	c.ServeJSON()

}

