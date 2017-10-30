package controllers

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	isLogin bool
}

func (c *BaseController) Prepare() {
	userLogin := c.GetSession("userLogin")
	if userLogin == nil {
		c.isLogin = false
	} else {
		c.isLogin = true
	}
	c.Data["isLogin"] = c.isLogin
}

func (c *BaseController) Go404() {
	c.Data["Website"] = "127.0.0.1:9090/login"
	c.Data["Email"] = "1485669171@qq.com"
	c.Data["EmailName"] = "erdoudou"
	c.TplName = "404.html"
}
