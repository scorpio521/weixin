package controllers

import (
	//"fmt"
	"github.com/astaxie/beego"
	_ "github.com/silenceper/wechat/message"
	"learnbeego/blog/tools"

	//"learnbeego/blog/models/class"
	//"encoding/json"
	//"learnbeego/blog/models/class"

	//"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	//openid := "odfKc1anwlIhIZo3Z26JHEypQZ9w"
	//tools.GetFinalPoster(openid+".jpeg","qrbg.jpeg",openid+"_qr.jpg")

	//tools.DownloadPic("http://test.hexiefamily.xin/site/img/001.png","static/img/openid/3.png",100,100,100)

	tools.GetQrPoster("showqrcode_yougu.jpeg","bg.jpeg","qrbg.jpeg")
	//tools.GetFinalPoster("3.png","qrbg.jpeg","sssssss.jpg")
	//test article poster
	//1、
	//tools.GetFinalPoster()
	//2、
	//article := &tools.Article{}
	//
	//// 先取二维码图片，oldfileName
	////生成新的二维码图片 fileName
	////新的二维码图片变成多大图片xsize,ysize
	////path 是总体的目录，无论是二维码还是生成的poster 目录，都在path目录下
	//path :="static/img/"
	//oldfileName := "qrstatic/showqrcode_yougu.jpeg"
	//fileName :="qrstatic/new_showqrcode_yougu.jpeg"
	//xsize :=100
	//ysize :=100
	//
	//
	////背景图片
	//backgroundimg :="qrstatic/bg.jpg"
	////合并之后的图片
	//posterName := "x1/poster123.jpg"
	//
	//
	//
	//articlePoster := tools.NewArticlePoster(posterName, article)
	//articlePosterBgService := tools.NewArticlePosterBg(
	//	backgroundimg,//背景图片
	//	articlePoster,
	//	//背景图片大小
	//	&tools.Rect{
	//		X0: 0,
	//		Y0: 0,
	//		X1: 400,
	//		Y1: 350,
	//	},
	//	//二维码放在什么位置
	//	&tools.Pt{
	//		X: 50,
	//		Y: 200,
	//	},
	//)
	//
	//_, filePath, err := articlePosterBgService.Generate(path,oldfileName,fileName,xsize,ysize)
	//
	//beego.Info("filePath is:",filePath)
	//if err != nil {
	//	beego.Info("something wrong",err)
	//	return
	//}





	//insertion :=&class.Contacts_answers{
	//	Customerpath: 	"dev_cj",
	//	Contactid:      1,
	//	Aid:1,
	//	Created_on:		time.Now().Local(),
	//}
	//insertid,err := insertion.Create()
	//if err != nil {
	//	beego.Info("插入出错了")
	//}
	//beego.Info("插入contact_answers成功")
	//beego.Info(insertid)

	//customerpath := "dev_cj"
	//maps,queryerror := class.GetAllQuestion(customerpath)
	//if queryerror!=""{
	//	beego.Info("something wrong or 没有数据")
	//}else{
	//	beego.Info("something right",maps)
	//}
	//var quesids []int64
	//quesidmap := make(map[int64]interface{})
	//for _,v :=range maps{
	//	idint,err := strconv.ParseInt(v["id"].(string), 10, 64)
	//	beego.Info("transfer to int",err)
	//	if err == nil{
	//		quesids = append(quesids,idint)
	//		quesidmap[idint] = v["t"]
	//	}
	//
	//}
	//beego.Info(quesids)
	//beego.Info(quesidmap)

	//aid :=1
	//qid,queryerror := class.GetAnswerQuestionid(int64(aid))
	//if queryerror!=""{
	//	beego.Info("something wrong or 没有数据")
	//}else{
	//	beego.Info("something right",qid)
	//}

	//customerpath := "dev_cj"
	//qid :=1
	//maps,queryerror := class.GetOneAnswer(customerpath,qid)
	//if queryerror!=""{
	//	beego.Info("something wrong or 没有数据")
	//}else{
	//	beego.Info("something right",maps)
	//}
	//mjson,_ :=json.Marshal(maps)
	//mString :=string(mjson)
	//beego.Info("print mString:%s",mString)

	//contactid := 1
	//customerpath := "dev_cj"
	//num, queryerr := class.GetIdNum(int64(contactid),customerpath)
	//if queryerr!="" {
	//	fmt.Println("something wrong ")
	//}
	//fmt.Println("something right",num)


	//insertion := &class.Contacts{
	//	Customerpath: "dev_cj",
	//	Openid:"dawdawdwadwadawdaw",
	//	Created_on:		time.Now().Local(),
	//}
	//var contactid int64
	//contactid,err := insertion.Create()
	//if err != nil {
	//	beego.Info("插入失败")
	//}else {
	//	beego.Info("插入成功")
	//	beego.Info(contactid)
	//}


	//
	//
	//config :=&class.Config{
	//	Customerpath: "dev0000",
	//}
	//
	//
	//configmap, queryerr := config.QueryByPath()
	////fmt.Println(customermap, err)
	////returnparmas:[{1 dev4 12121 1212121}]
	////returnparmas: [] code = 1 ; msg = the sql sentence is wrong,please check
	//beego.Info(configmap)
	//beego.Info(queryerr)
	//if queryerr != ""{
	//	c.TplName = "bad.html"
	//	return
	//}
	//if configmap == nil {
	//	beego.Info("configmap is empty")
	//	c.TplName = "bad.html"
	//	return
	//}


	c.Data["Website"] = "kkkk"
	c.Data["Email"] = "kkkkkkk@qq.com"
	c.TplName = "index.tpl"
}

