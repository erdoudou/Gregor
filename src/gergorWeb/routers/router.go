package routers

import (
	"gergorWeb/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/404", &controllers.BaseController{}, "*:Go404")
	//添加一个路由
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/register", &controllers.RegisterController{})

	beego.Router("/home", &controllers.HomeController{})

	beego.Router("/logout", &controllers.LogoutController{})
}
