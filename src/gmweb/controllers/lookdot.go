package controllers

import (
	"Public/db"
	"fmt"
	"strconv"
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
	playerdot := c.GetString("playerdot")

	if "" != playername && playerdot == "" {
		c.lookDot(playername)
		return
	}

	if "" != playername && playerdot != "" {
		dot, err := strconv.Atoi(playerdot)
		if err != nil {
			fmt.Println("转换数据失败", err, playerdot)
			c.Data["json"] = map[string]interface{}{"code": 0, "message": "修改玩家科技点失败"}
			c.ServeJSON()
			return
		}
		ret := changeDot(playername, int64(dot))
		if ret {
			c.Data["json"] = map[string]interface{}{"code": 1, "message": "修改玩家科技点成功"}
			c.ServeJSON()
			return
		}

		c.Data["json"] = map[string]interface{}{"code": 0, "message": "操作失败"}
		c.ServeJSON()
		return
	}
	fmt.Println("网页传入的参数出现问题", playername, playerdot)
	c.Data["json"] = map[string]interface{}{"code": 0, "message": "传入的参数有问题"}
	c.ServeJSON()
}

func changeDot(playername string, playerdot int64) bool {
	//在数据库中进行查找
	row, ret := db.QueryNormal("SELECT sciencedot FROM playerdot WHERE globaluid = (SELECT globaluid FROM accountconvert WHERE username=?)", playername)
	if ret == false || row == nil {
		fmt.Println("查询错误")
		return false
	}

	defer row.Close()

	var dot int64

	if row.Next() == false {
		fmt.Println("没有数据库返回数据111111")
		return false
	}

	err := row.Scan(&dot)

	if err != nil {
		fmt.Println("解析失败")
		return false
	}

	//更新数据库
	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE playerdot SET sciencedot=? WHERE globaluid= (SELECT globaluid FROM accountconvert WHERE username=?)"

	db.ExecUpdate(idchan1, reschan1, updateCmd,
		playerdot, playername)

	return true
}

func (c *LookdotController) lookDot(playername string) {
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
