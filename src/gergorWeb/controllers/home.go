package controllers

import (
	"fmt"
)

type HomeController struct {
	BaseController
}

func (c *HomeController) Get() {
	if !c.isLogin {
		c.Redirect("/login", 302)
		fmt.Println("没有登录")
		return
	}
	c.TplName = "home/home.html"

}
