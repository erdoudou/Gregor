package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "127.0.0.1:9090/login"
	c.Data["Email"] = "1485669171@qq.com"
	c.Data["EmailName"] = "erdoudou"
	c.TplName = "404.html"
}
