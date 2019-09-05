package routers

import (
	"learnbeego/blog/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/wx",&controllers.WXController{})
	beego.Router("/login", &controllers.UserController{}, `get:PageLogin`)
	beego.Router("/register", &controllers.UserController{}, `post:Register`)
	beego.Router("/reallogin", &controllers.UserController{}, `post:Reallogin`)
	beego.Router("/wxframe", &controllers.WXframeController{}, `get:Hello;post:Hello`)
}
