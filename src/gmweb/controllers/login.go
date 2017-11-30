package controllers

import (
	"Public/db"
	"fmt"
)

type user struct {
	Name     string `form:"username"`
	Password string `form:"password"`
}

type LoginController struct {
	BaseController
}

func (c *LoginController) Get() {
	fmt.Println("进入登录页面")
	c.TplName = "login/login.html"

}

//添加视图模板变量，指定模板文件
func (c *LoginController) Post() {
	u := user{}
	if err := c.ParseForm(&u); err != nil {
		fmt.Println("handle error")
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "登录失败"}
		c.ServeJSON()
	}

	if "" == u.Name || "" == u.Password {
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "登陆失败，你输入的账号或密码为空"}
		c.ServeJSON()
		return
	}

	row, ret := db.QueryNormal("select password from user where username = ?", u.Name)
	if ret == false || row == nil {
		fmt.Println("查询错误")
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "登陆失败，数据库错误，请联系管理员"}
		c.ServeJSON()
		return
	}

	defer row.Close()

	var pwd string

	if row.Next() == false {
		fmt.Println("没有数据库返回数据")
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "登陆失败，你没有账号请联系管理员进行添加"}
		c.ServeJSON()
		return
	}

	err := row.Scan(&pwd)

	if err != nil {
		fmt.Println("解析失败")
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "登陆失败，解析数据错误"}
		c.ServeJSON()
		return
	}

	if pwd != u.Password {
		fmt.Println("传入的密码和数据库密码对不上", pwd, u.Password)
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "登录失败，你输入的密码不正确"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{}{"code": 1}
	c.ServeJSON()
	c.SetSession("userLogin", "1")
	fmt.Println("登录成功")
}
