package controllers

import (
	"fmt"
)

type HomeController struct {
	BaseController
}

func (c *HomeController) Get() {
	fmt.Println("跳转到主页")
	if !c.isLogin {
		fmt.Println("没有登录")
		c.Redirect("/login", 302)
		return
	}
	fmt.Println("已经正常登陆")
	c.TplName = "home/home.html"

}
