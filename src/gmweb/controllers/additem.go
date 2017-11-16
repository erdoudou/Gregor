package controllers

import (
	"Public/db"
	"fmt"
	//"strconv"
	//"strings"
	"gmweb/orderform"
)

type AdditemController struct {
	BaseController
}

func (c *AdditemController) Get() {
	if !c.isLogin {
		c.Redirect("/login", 302)
		fmt.Println("没有登录")
		return
	}
	c.TplName = "look/additem.html"

}

func (c *AdditemController) Post() {
	playername := c.GetString("playername")
	itemtempid := c.GetString("itemtempid")
	fmt.Println("打印出网页传过来的数据", playername, itemtempid)

	if "" == playername || "" == itemtempid {
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "你输入的信息有误，请重新输入"}
		c.ServeJSON()
		return
	}

	//首先到steam上面进行添加物品
	//itemarray 物品的模板id
	//从数据库中查找平台uid
	row, ret := db.QueryNormal("SELECT platformuid FROM accountconvert WHERE username=?", playername)
	if ret == false || row == nil {
		fmt.Println("查询错误")
		return
	}

	defer row.Close()

	var platformuid string

	if row.Next() == false {
		fmt.Println("没有数据库返回数据111111")
		return
	}

	err := row.Scan(&platformuid)

	if err != nil {
		fmt.Println("解析失败")
		return
	}

	isAdd := orderMgr.TryToAddItemOnSteam(platformuid, itemtempid, 0, 0)
	if isAdd == false {
		fmt.Println("添加物品不成功")
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "在数据库或者steam上面添加物品失败"}
		c.ServeJSON()
		return
	}

	fmt.Println("添加物品成功")
	c.Data["json"] = map[string]interface{}{"code": 1}
	c.ServeJSON()

}
