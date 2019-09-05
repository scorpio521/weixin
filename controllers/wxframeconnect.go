package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/cache"
	"github.com/silenceper/wechat/material"
	"github.com/silenceper/wechat/message"
	"learnbeego/blog/models/class"
	"learnbeego/blog/tools"
	"strconv"
	"strings"
	"time"
)
type WXframeController struct {
	beego.Controller
}

func (c *WXframeController) Hello(){
	customerpath := beego.AppConfig.String("customerpath")
	if beego.AppConfig.String("runmode")=="dev" {
		customerpath = beego.AppConfig.String("customerpathdev")
	}
	con := &class.Config{Customerpath:customerpath}
	customermap, queryerr := con.QueryByPath()
	//fmt.Println(customermap, err)
	//returnparmas:[{1 dev4 12121 1212121}]
	//returnparmas: [] code = 1 ; msg = the sql sentence is wrong,please check
	beego.Info(customermap,queryerr)
	appconfig := make(map[string]string)
	for _, v := range customermap {
		appconfig["appid"] = v.Appid
		appconfig["appsec"] = v.Appsec
	}
	beego.Info("wechat account appid is",appconfig["appid"],"wechat account appsec is:",appconfig["appsec"])

	memCache :=cache.NewMemcache("127.0.0.1:11211")
	//配置微信参数
	config := &wechat.Config{
		AppID:          appconfig["appid"],
		AppSecret:      appconfig["appsec"],
		Token:          "weixin",
		EncodingAESKey: "NliWjMcxkisTaDWjswxHY0K6CZgesA9XI5Li7Az7ynl",
		Cache:          memCache,
	}
	wc := wechat.NewWechat(config)
	// 传入request和responseWriter
	server := wc.GetServer(c.Ctx.Request, c.Ctx.ResponseWriter)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
		beego.Info("msg mixmessage is:",msg)
		content := msg.Content
		openid := msg.FromUserName
		beego.Info(openid)

		//检测是否有用户头像
		//if tools.CheckNotExist("static/img/openid/"+openid+".jpeg")==true {//没有用户头像，获取
		//	beego.Info("用户头像不存在")
		//	user :=wc.GetUser()
		//	userinfo,usererr :=user.GetUserInfo(openid)
		//	if usererr!=nil{
		//		beego.Info("获取用户信息错误：",usererr)
		//		text := message.NewText("用户信息获取出错了")
		//		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		//	}
		//
		//
		//	imgPath := "static/img/openid/"
		//	imgUrl := userinfo.Headimgurl
		//	title =userinfo.Nickname+"邀请您来一起参加"
		//
		//	//fileName := path.Base(imgUrl)
		//	fileName := openid+".jpeg"
		//	destpath := imgPath+fileName
		//	//downloadfileerr	:=tools.WrongDownloadPic(imgUrl,imgPath,fileName)
		//	downloadfileerr	:=tools.DownloadPic(imgUrl,destpath,100,100,100)
		//	if downloadfileerr!=nil{
		//		text := message.NewText("下载用户头像出错")
		//		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		//	}
		//
		//}

		beego.Info("不检测用户头像存不存在，直接覆盖")
		user :=wc.GetUser()
		userinfo,usererr :=user.GetUserInfo(openid)
		if usererr!=nil{
			beego.Info("获取用户信息错误：",usererr)
			text := message.NewText("用户信息获取出错了")
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		}


		imgPath := "static/img/openid/"
		imgUrl := userinfo.Headimgurl
		title :=userinfo.Nickname+"邀请您来一起参加"

		//fileName := path.Base(imgUrl)
		fileName := openid+".jpeg"
		destpath := imgPath+fileName
		//downloadfileerr	:=tools.WrongDownloadPic(imgUrl,imgPath,fileName)
		downloadfileerr	:=tools.DownloadPic(imgUrl,destpath,100,100,100)
		if downloadfileerr!=nil{
			text := message.NewText("下载用户头像出错")
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		}

		var contactid int64
		if msg.MsgType == message.MsgTypeText {
			//回复消息：演示回复用户发送的消息



			//查看是否有bizmsgmenuid
			//也就是这里一定会传过来值，只不过默认值会变为0
			//t := reflect.TypeOf(msg)
			//if _, ok := t.FieldByName("BizMsgMenuId") ; ok {
			//	beego.Info("存在bizmsgmenuid----------")
			//}else{
			//	beego.Info("不存在bizmsgmenuid----------")
			//}
			if msg.BizMsgMenuId!=0 {
				beego.Info("存在bizmsgmenuid----------")




				content = "欢迎关注公众号"
				contact :=&class.Contacts{
					Customerpath: customerpath,
					Openid:       openid,
				}

				contactmap, queryerr := contact.QueryByIdPath()

				if queryerr!="" {
					text := message.NewText("查询出错了")
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}
				//fmt.Println(customermap, err)
				//returnparmas:[{1 dev4 12121 1212121}]
				//returnparmas: [] code = 1 ; msg = the sql sentence is wrong,please check
				beego.Info(contactmap,queryerr)

				if contactmap == nil {
					insertion :=&class.Contacts{
						Customerpath: 	customerpath,
						Openid:       	openid,
						Created_on:		time.Now().Local(),
					}
					contactid,err := insertion.Create()
					if err != nil {
						text := message.NewText("插入出错了")
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					}
					beego.Info("插入contact成功")
					beego.Info(contactid)
				}else {
					for _, v := range contactmap {
						contactid = v.Id
					}
					beego.Info("获取contact成功,contactid is :",contactid)
					//beego.Info(contactid)
				}
				//bizmsgmenu id start
				menuid := msg.BizMsgMenuId
				quesmenuid, queryerr := class.GetAnswerQuestionid(menuid)
				beego.Info("Get bizmsgmenuid is quesmenuid:",quesmenuid)
				//aidnumint,err := strconv.Atoi(aidnum.(string))
				if queryerr!="" {
					text := message.NewText("quesmenuid 语句出错")
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}
				quesmap, queryerr := class.GetAllQuestion(customerpath)
				beego.Info("Get question map :",quesmap)
				if queryerr!="" {
					text := message.NewText("quesimon map 语句出错")
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}
				var quesids []int64
				quesidmap := make(map[int64]interface{})
				for _,v :=range quesmap{
					idint,err := strconv.ParseInt(v["id"].(string), 10, 64)
					if err == nil{
						quesids = append(quesids,idint)
						quesidmap[idint] = v["t"]
					}

				}
				beego.Info("quesid is :",quesids,"quesidmap is:",quesidmap)

				objquesmaps, queryerr := class.GetObjQuestionid(contactid,customerpath)
				beego.Info("Get objquesmap map :",objquesmaps)
				if queryerr!="" {
					text := message.NewText("objquesmap map 语句出错")
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}

				inobjectflag := 0
				var objquesids []int64
				for _,v :=range objquesmaps{
					idint,err := strconv.ParseInt(v["qid"].(string), 10, 64)
					objquesids = append(objquesids,idint)
					if err == nil && idint==quesmenuid {
						inobjectflag = 1
					}
				}
				if inobjectflag==0{
					insertion :=&class.Contact_answers{
						Customerpath: 	customerpath,
						Contactid:      contactid,
						Aid:menuid,
						Created_on:		time.Now().Local(),
					}
					insertid,err := insertion.Create()
					if err != nil {
						text := message.NewText("插入出错了")
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					}
					objquesids = append(objquesids,quesmenuid)
					beego.Info("插入contact_answers成功")
					beego.Info(insertid)
				}
				
				numsarr := &tools.Arrcommonfunc{}
				leftarr :=numsarr.Diff(objquesids,quesids)
				beego.Info("left num arr is :",leftarr)
				if len(leftarr)>0{
					length := len(leftarr)
					var choosenumid int64
					if(length==1){
						choosenumid = leftarr[0]
					}else{
						choosenumarr, err := numsarr.Random(leftarr,1)
						choosenumid =choosenumarr[0]
						beego.Info("choose num arr is:",choosenumarr,"err is:",err)
					}
					beego.Info("choose numid is :",choosenumid)


					if _, ok := quesidmap[choosenumid]; !ok {
						text := message.NewText("问题数据不全")
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					}
					quest :=quesidmap[choosenumid].(string)

					amaps,queryerror := class.GetOneAnswer(customerpath, int(choosenumid))
					beego.Info(amaps)
					if queryerror!="" {
						text := message.NewText("查询到answers的sql 语句出错")
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					}
					//mjson,_ :=json.Marshal(amaps)
					//amapstring :=string(mjson)
					//beego.Info("2Get something in answers table:",amapstring)

					kfsendmsg := new(message.KfMsg)
					kfsendmsg.Touser = openid
					kfsendmsg.Msgtype = "msgmenu"
					kfsendmsg.Msgmenu.Headcontent = quest
					kfsendmsg.Msgmenu.Tailcontent = ""
					for _,v :=range amaps{
						ids := v["id"].(string)
						ts := v["t"].(string)
						//beego.Info("id:",v["id"],"t:",v["t"],"------string part is-----,id:",ids,"t:",ts)
						kfsendmsg.Msgmenu.List = append(kfsendmsg.Msgmenu.List , message.SingleList{Id: ids, Content: ts})
					}
					beego.Info(kfsendmsg)
					jsonmenu, _ := json.Marshal(kfsendmsg)
					beego.Info("json:,%s", string(jsonmenu))
					kfcustomer := message.NewKfcustomer(wc.Context)

					kfcustomer.KfSend(kfsendmsg)




				}else{
					//先查mediaid
					mediaid, mediaerr := class.GetMediaid(customerpath,openid)
					if mediaerr!="" {
						text := message.NewText("查询到media的sql 语句出错")
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					}
					beego.Info("Get something in media table:",mediaid,"error:",mediaerr)
					//有则调取mediaid,如果没有生成海报图片并上传给微信，同时获取mediaid,并且创建mediaid,同时返还给前端。
					if mediaid==""{
						//这个接口是本地有或者没有都会有一个图片文件呢
						tools.GetFinalPoster(openid+".jpeg","qrbg.jpeg",openid+"_qr.jpg",title)
						filename := "static/img/openid/"+openid+"_qr.jpg"
						if tools.CheckNotExist(filename)==false{
							mate := wc.GetMaterial()
							mediares,err:=mate.MediaUpload(material.MediaTypeImage,filename)
							if err!=nil{
								text := message.NewText("media 上传失败，请重新上传")
								return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
							}
							//insert into mediaid 数据库
							mediaid = mediares.MediaID
							t := time.Now().Unix()
							t += 86400*2//两天之后失效
							insertion :=&class.Media{
								Customerpath: 	customerpath,
								Openid		:   openid,
								Mediaid		:	mediaid,
								Expiretime	:	time.Unix(t, 0).Format("2006-01-02 15:04:05"),
							}
							mediainsertid,err := insertion.Create()
							if err != nil {
								text := message.NewText("插入出错了")
								return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
							}
							beego.Info("插入media成功")
							beego.Info(mediainsertid)

						}
					}



					//存mediaid,如果存在则发送mediaid

					//发送图片
					kfimgsendmsg := new(message.KfImgMsg)
					kfimgsendmsg.Touser = openid
					kfimgsendmsg.Msgtype = "image"
					kfimgsendmsg.Image.Mediaid = mediaid
					beego.Info(kfimgsendmsg)
					imgjsonmenu, _ := json.Marshal(kfimgsendmsg)
					beego.Info("json:,%s", string(imgjsonmenu))
					imgkfcustomer := message.NewKfcustomer(wc.Context)

					imgkfcustomer.KfImgSend(kfimgsendmsg)

					text := message.NewText("生成图片")

					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}


			} else {
				beego.Info("不存在bizmsgmenuid------")
			}





		}else if msg.MsgType == message.MsgTypeEvent {
			//回复消息：演示回复用户发送的消息
			if msg.Event==message.EventSubscribe {
				content = "欢迎关注该公众号"
				contact :=&class.Contacts{
					Customerpath: customerpath,
					Openid:       "",
				}
				contactmap, queryerr := contact.QueryByIdPath()
				//fmt.Println(customermap, err)
				//returnparmas:[{1 dev4 12121 1212121}]
				//returnparmas: [] code = 1 ; msg = the sql sentence is wrong,please check
				beego.Info(contactmap)
				beego.Info(queryerr)


				for _, v := range contactmap {
					contactid = v.Id
				}
				beego.Info(contactid)



			}else if msg.Event==message.EventClick {

				content = "欢迎关注该公众号"
				contact :=&class.Contacts{
					Customerpath: customerpath,
					Openid:       openid,
				}
				//先获取头像



				contactmap, queryerr := contact.QueryByIdPath()

				if queryerr!="" {
					text := message.NewText("查询出错了")
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}
				//fmt.Println(customermap, err)
				//returnparmas:[{1 dev4 12121 1212121}]
				//returnparmas: [] code = 1 ; msg = the sql sentence is wrong,please check
				beego.Info(contactmap,queryerr)

				if contactmap == nil {
					insertion :=&class.Contacts{
						Customerpath: 	customerpath,
						Openid:       	openid,
						Created_on:		time.Now().Local(),
					}
					contactid,err := insertion.Create()
					if err != nil {
						text := message.NewText("插入出错了")
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					}
					beego.Info("插入contact成功")
					beego.Info(contactid)
				}else {
					for _, v := range contactmap {
						contactid = v.Id
					}
					beego.Info("获取contact成功,contactid is :",contactid)
					//beego.Info(contactid)
				}
				aidnum, queryerr := class.GetAidNum(contactid,customerpath)
				beego.Info("Get something in questions table:",aidnum)
				aidnumint,err := strconv.Atoi(aidnum.(string))
				if queryerr!="" {
					text := message.NewText("查询到num的sql 语句出错")
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}
				beego.Info("1Get something in questions table:",aidnumint,"error:",err)// 这是个interface,不是int类型
				qidnum, queryerr := class.GetQidNum(customerpath)
				if queryerr!="" {
					text := message.NewText("查询到num的sql 语句出错")
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}
				qidnumint,err := strconv.Atoi(qidnum.(string))
				beego.Info("1Get something in answers table:",qidnumint,"error:",err)
				//res := make(orm.Params)
				//nums, err := o.Raw("SELECT name, value FROM options_table").RowsToMap(&res, "name", "value")
				// res is a map[string]interface{}{
				//	"total": 100,
				//	"found": 200,
				// }
				//interface 类型的不允许比较
				if aidnumint<qidnumint {
					qmaps,queryerror := class.GetOneQuestion(customerpath)
					if queryerror!="" {
						text := message.NewText("查询到questions的sql 语句出错")
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					}
					qid,err := strconv.Atoi(qmaps["id"].(string))
					beego.Info("2Get something in questions table:",qmaps,"qid is:",qid,"error:",err)
					beego.Info("customerpath is :",customerpath,"qid is :",qid)
					quest := qmaps["t"].(string)
					amaps,queryerror := class.GetOneAnswer(customerpath,qid)
					beego.Info(amaps)
					if queryerror!="" {
						text := message.NewText("查询到answers的sql 语句出错")
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					}
					//mjson,_ :=json.Marshal(amaps)
					//amapstring :=string(mjson)
					//beego.Info("2Get something in answers table:",amapstring)

					kfsendmsg := new(message.KfMsg)
					kfsendmsg.Touser = openid
					kfsendmsg.Msgtype = "msgmenu"
					kfsendmsg.Msgmenu.Headcontent = quest
					kfsendmsg.Msgmenu.Tailcontent = ""
					for _,v :=range amaps{
						ids := v["id"].(string)
						ts := v["t"].(string)
						//beego.Info("id:",v["id"],"t:",v["t"],"------string part is-----,id:",ids,"t:",ts)
						kfsendmsg.Msgmenu.List = append(kfsendmsg.Msgmenu.List , message.SingleList{Id: ids, Content: ts})
					}
					beego.Info(kfsendmsg)
					jsonmenu, _ := json.Marshal(kfsendmsg)
					beego.Info("json:,%s", string(jsonmenu))
					kfcustomer := message.NewKfcustomer(wc.Context)

					kfcustomer.KfSend(kfsendmsg)

					//return &message.Reply{MsgType: message.MsgTypeText, MsgData: tc}




				}else{
					text := message.NewText("您已经做过测试了哦")
					//先查mediaid
					mediaid, mediaerr := class.GetMediaid(customerpath,openid)
					beego.Info("Get something in media table:",mediaid,"error:",mediaerr)
					if mediaerr!="" {
						text := message.NewText("查询到media的sql 语句出错")
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					}

					//有则调取mediaid,如果没有生成海报图片并上传给微信，同时获取mediaid,并且创建mediaid,同时返还给前端。
					if mediaid==""{
						//这个接口是本地有或者没有都会有一个图片文件呢
						tools.GetFinalPoster(openid+".jpeg","qrbg.jpeg",openid+"_qr.jpg",title)
						filename := "static/img/openid/"+openid+"_qr.jpg"
						if tools.CheckNotExist(filename)==false{
							mate := wc.GetMaterial()
							mediares,err:=mate.MediaUpload(material.MediaTypeImage,filename)
							if err!=nil{
								beego.Info("上传失败的原因是:",err)
								text := message.NewText("media 上传失败，请重新上传")
								return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
							}
							//insert into mediaid 数据库
							mediaid = mediares.MediaID
							t := time.Now().Unix()
							t += 86400*2//两天之后失效
							insertion :=&class.Media{
								Customerpath: 	customerpath,
								Openid		:   openid,
								Mediaid		:	mediaid,
								Expiretime	:	time.Unix(t, 0).Format("2006-01-02 15:04:05"),
							}
							mediainsertid,err := insertion.Create()
							if err != nil {
								text := message.NewText("插入出错了")
								return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
							}
							beego.Info("插入media成功")
							beego.Info(mediainsertid)

						}
					}

					kfimgsendmsg := new(message.KfImgMsg)
					kfimgsendmsg.Touser = openid
					kfimgsendmsg.Msgtype = "image"
					kfimgsendmsg.Image.Mediaid = mediaid
					beego.Info(kfimgsendmsg)
					imgjsonmenu, _ := json.Marshal(kfimgsendmsg)
					beego.Info("json:,%s", string(imgjsonmenu))
					imgkfcustomer := message.NewKfcustomer(wc.Context)

					imgkfcustomer.KfImgSend(kfimgsendmsg)

					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}



			}else if msg.Event==message.EventScan {
				strpos := strings.Index(msg.EventKey, "qrscene_")

				if strpos !=-1 {
					unique_openid := msg.EventKey[strpos:]
					users :=&class.Users{
						Customerpath: customerpath,
						Unique_openid:  unique_openid,
					}

					usersmap, queryerr := users.QueryByUnique()

					if queryerr!="" {
						text := message.NewText("查询出错了")
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					}

					beego.Info(usersmap)





				}
			}

		}else{
			//回复消息：演示回复用户发送的消息
			content = "欢迎关注该公众号"
		}
		text := message.NewText(content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}

	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}


















