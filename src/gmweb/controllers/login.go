package controllers

import (
	"Public/db"
	"fmt"

	//"github.com/astaxie/beego"
	//"github.com/astaxie/beego/utils/pagination"
)

type user struct {
	Name     string `form:"username"`
	Password string `form:"password"`
}

type LoginController struct {
	BaseController
}

func (c *LoginController) Get() {
	c.TplName = "login/login.html"

}

//添加视图模板变量，指定模板文件
func (c *LoginController) Post() {
	u := user{}
	if err := c.ParseForm(&u); err != nil {
		fmt.Println("handle error")
		//c.Ctx.WriteString("0")
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
	}

	if "" == u.Name || "" == u.Password {
		//c.Ctx.WriteString("0")
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
		return
	}

	row, ret := db.QueryNormal("select password from user where username = ?", u.Name)
	if ret == false || row == nil {
		fmt.Println("查询错误")
		//c.Ctx.WriteString("0")
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
		return
	}

	defer row.Close()

	var pwd string

	if row.Next() == false {
		fmt.Println("没有数据库返回数据")
		//c.Ctx.WriteString("0")
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
		return
	}

	err := row.Scan(&pwd)

	if err != nil {
		fmt.Println("解析失败")
		//c.Ctx.WriteString("0")
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
		return
	}

	if pwd != u.Password {
		fmt.Println("传入的密码和数据库密码对不上", pwd, u.Password)
		//c.Ctx.WriteString("0")
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
		return
	}

	//c.Ctx.WriteString("1")
	c.Data["json"] = map[string]interface{}{"code": 1}
	c.ServeJSON()
	c.SetSession("userLogin", "1")
	//c.Redirect("/home", 200)

	fmt.Println("解析数据成功", u.Password, pwd)
}
