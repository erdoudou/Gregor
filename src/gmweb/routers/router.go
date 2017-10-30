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

	beego.Router("/lookdot", &controllers.LookdotController{})

	beego.Router("/lookitem", &controllers.LookitemController{})
	beego.Router("/additem", &controllers.AdditemController{})

	beego.Router("/logout", &controllers.LogoutController{})
}
