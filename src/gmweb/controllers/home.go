package controllers

import (
	"fmt"
)

type HomeController struct {
	BaseController
}

func (c *HomeController) Get() {
	fmt.Println("进入主页跳转流程")
	if !c.isLogin {
		fmt.Println("没有登录")
		c.Redirect("/login", 302)
		return
	}
	fmt.Println("跳转到home页面")
	c.TplName = "home/home.html"
}
