package controllers

import (
	"Public/db"
	"fmt"
)

type LookdotController struct {
	BaseController
}

func (c *LookdotController) Get() {
	if !c.isLogin {
		c.Redirect("/login", 302)
		fmt.Println("没有登录")
		return
	}
	c.TplName = "look/lookdot.html"

}

//添加视图模板变量，指定模板文件
func (c *LookdotController) Post() {
	playername := c.GetString("playername")

	if "" == playername {
		c.Ctx.WriteString("0")
		return
	}

	//在数据库中进行查找
	row, ret := db.QueryNormal("SELECT sciencedot FROM playerdot WHERE globaluid = (SELECT globaluid FROM accountconvert WHERE username=?)", playername)
	if ret == false || row == nil {
		fmt.Println("查询错误")
		c.Ctx.WriteString("0")
		return
	}

	defer row.Close()

	var dot int32

	if row.Next() == false {
		fmt.Println("没有数据库返回数据")
		c.Ctx.WriteString("0")
		return
	}

	err := row.Scan(&dot)

	if err != nil {
		fmt.Println("解析失败")
		c.Ctx.WriteString("0")
		return
	}

	fmt.Println("查到该玩家的科技点", dot)
	c.Data["json"] = map[string]interface{}{"code": 1, "dot": dot}
	c.ServeJSON()
}
