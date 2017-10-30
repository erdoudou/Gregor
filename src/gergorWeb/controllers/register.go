package controllers

import (
	"Public/db"
	"fmt"
)

type RegisterController struct {
	BaseController
}

func (c *RegisterController) Get() {
	if !c.isLogin {
		c.Redirect("/login", 302)
		fmt.Println("没有登录")
		return
	}

	c.TplName = "login/register.html"

}

//添加视图模板变量，指定模板文件
func (c *RegisterController) Post() {
	u := user{}
	if err := c.ParseForm(&u); err != nil {
		fmt.Println("handle error")
		//c.Ctx.WriteString("0")
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "解析错误"}
		c.ServeJSON()
	}

	if "" == u.Name || "" == u.Password {
		//c.Ctx.WriteString("0")
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "输入的用户名或者密码错误"}
		c.ServeJSON()
		return
	}

	//将数据直接插入数据库中
	idchan := make(chan int, 2)
	reschan := make(chan bool, 2)
	db.ExecInsert(idchan, reschan, "INSERT  user(username,password,power) value (?,?,?)",
		u.Name, u.Password, 0)

	c.Data["json"] = map[string]interface{}{"code": 1}
	c.ServeJSON()

	fmt.Println("解析数据成功", u)
}
