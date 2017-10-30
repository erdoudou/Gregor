package controllers

import (
	"Public/db"
	"fmt"
	"strconv"
	"strings"
)

type LookitemController struct {
	BaseController
}

type item struct {
	Itemuid    string
	Itemtempid int32
}

func (c *LookitemController) Get() {
	if !c.isLogin {
		c.Redirect("/login", 302)
		fmt.Println("没有登录")
		return
	}

	c.TplName = "look/lookitem.html"

}

//添加视图模板变量，指定模板文件
func (c *LookitemController) Post() {
	playername := c.GetString("playername")
	itemuid := c.GetString("itemuid")
	itemtempid := c.GetString("itemtempid")

	fmt.Println("打印出收到的消息数据", itemuid, playername, itemtempid)

	//查出玩家所有的物品
	if "" != playername && "" == itemuid {
		c.lookItem(playername, itemtempid)
		return
	}

	if "" != itemuid && playername != "" {
		fmt.Println("进入删除物品流程dede", itemuid, playername)
		ret := deleteBagItem(itemuid, playername)
		if ret != true {
			c.Data["json"] = map[string]interface{}{"code": 0, "message": "删除背包数据失败"}
			c.ServeJSON()
			return
		}
		c.deleteItem(itemuid)
	}

}

func (c *LookitemController) deleteItem(itemuid string) {

	row, ret := db.QueryNormal("SELECT bdel FROM item  WHERE  itemuid = ?", itemuid)
	if ret == false {
		fmt.Println("从数据库中取出玩家dot数据错误")
		return
	}

	defer row.Close()

	if row == nil || row.Next() == false {
		fmt.Println("玩家物品中没有改数据数据", itemuid)
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "玩家物品中没有改数据数据"}
		c.ServeJSON()
	} else {
		updateCmd := "UPDATE item SET bdel=1 WHERE itemuid = ?"
		idchan1 := make(chan int, 2)
		reschan1 := make(chan bool, 2)

		db.ExecUpdate(idchan1, reschan1, updateCmd, itemuid)

		c.Data["json"] = map[string]interface{}{"code": 1}
		c.ServeJSON()
	}

}

func deleteBagItem(itemUid, playername string) bool {

	//在数据库中进行查找
	row, ret := db.QueryNormal("SELECT itemslots FROM backpack WHERE useruid = (SELECT globaluid FROM accountconvert WHERE username=?)", playername)
	if ret == false || row == nil {
		fmt.Println("查询错误")
		return false
	}

	defer row.Close()

	var itemids string

	if row.Next() == false {
		fmt.Println("没有数据库返回数据111111")
		return false
	}

	err := row.Scan(&itemids)

	if err != nil {
		fmt.Println("解析失败")
		return false
	}

	itemContent := strings.Replace(itemids, itemUid+",", "", -1)

	//更新数据库
	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE backpack SET itemslots=? WHERE useruid= (SELECT globaluid FROM accountconvert WHERE username=?)"

	fmt.Println("更新背包物品============", itemContent)

	db.ExecUpdate(idchan1, reschan1, updateCmd,
		itemContent, playername)

	return true

}

//查看物品流程
func (c *LookitemController) lookItem(playername string, itemtempid string) {
	//在数据库中进行查找
	row, ret := db.QueryNormal("SELECT itemslots FROM backpack WHERE useruid = (SELECT globaluid FROM accountconvert WHERE username=?)", playername)
	if ret == false || row == nil {
		fmt.Println("查询错误")
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
		return
	}

	defer row.Close()

	var itemids string

	if row.Next() == false {
		fmt.Println("没有数据库返回数据111111")
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
		return
	}

	err := row.Scan(&itemids)

	if err != nil {
		fmt.Println("解析失败")
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
		return
	}

	var itemuidSli []int64

	for _, v := range strings.Split(itemids, ",") {
		if v != "" {
			id, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("背包数据的转换错误", err)
				return
			}

			itemuidSli = append(itemuidSli, int64(id))
		}
	}

	var itemSli []item
	htmlIempid, err := strconv.Atoi(itemtempid)
	if err != nil {
		htmlIempid = -1
	}

	for _, v := range itemuidSli {
		tempid := findItemtempid(v)
		if htmlIempid != -1 {
			if int32(htmlIempid) == tempid {
				prop := item{
					Itemuid:    strconv.Itoa(int(v)),
					Itemtempid: tempid,
				}
				itemSli = append(itemSli, prop)
			}
		} else if tempid != -1 {
			prop := item{
				Itemuid:    strconv.Itoa(int(v)),
				Itemtempid: tempid,
			}
			itemSli = append(itemSli, prop)
		}
	}

	fmt.Println("查到该玩家的物品信息deed", itemSli)
	if len(itemSli) == 0 {
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
	}

	c.Data["json"] = itemSli
	c.ServeJSON()

}

//根据一个物品的uid查找这个物品的itemid
func findItemtempid(uid int64) (tempid int32) {
	row, ret := db.QueryNormal("SELECT itemtempid FROM item WHERE  bdel = 0 AND itemuid = ?", uid)
	if ret == false || row == nil {
		fmt.Println("查询错误")
		return -1
	}

	defer row.Close()

	var itemid int32

	if row.Next() == false {
		fmt.Println("没有数据库返回数据22222222", uid)
		return -1
	}

	err := row.Scan(&itemid)

	if err != nil {
		fmt.Println("解析失败")
		return -1
	}

	return itemid
}
