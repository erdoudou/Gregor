package controllers

import (
	"Public/db"
	"fmt"
	"strconv"
	"strings"
)

type LookBoxController struct {
	BaseController
}

func (c *LookBoxController) Get() {
	if !c.isLogin {
		c.Redirect("/login", 302)
		fmt.Println("没有登录")
		return
	}

	c.TplName = "look/lookbox.html"

}

type box struct {
	Boxuid    string
	Boxtempid int32
}

//添加视图模板变量，指定模板文件
func (c *LookBoxController) Post() {
	playername := c.GetString("playername")
	boxuid := c.GetString("boxuid")

	fmt.Println("打印出收到的消息数据", playername)

	//查出玩家所有的宝箱
	if "" != playername && "" == boxuid {
		c.lookBox(playername)
		return
	}

	if "" != boxuid && playername != "" {
		fmt.Println("进入删除物品流程dede", boxuid, playername)
		ret := deleteBoxhouse(boxuid, playername)
		if ret != true {
			c.Data["json"] = map[string]interface{}{"code": 0, "message": "删除背包数据失败"}
			c.ServeJSON()
			return
		}
		c.deleteBox(boxuid)
	}

}

func (c *LookBoxController) deleteBox(boxuid string) {

	row, ret := db.QueryNormal("SELECT bdel FROM treasurebox  WHERE  treasureuid = ?", boxuid)
	if ret == false {
		fmt.Println("从数据库中取出玩家box数据错误")
		return
	}

	defer row.Close()

	if row == nil || row.Next() == false {
		fmt.Println("玩家物品中没有改数据数据", boxuid)
		c.Data["json"] = map[string]interface{}{"code": 0, "message": "玩家物品中没有改数据数据"}
		c.ServeJSON()
	} else {
		updateCmd := "UPDATE treasurebox SET bdel=1 WHERE treasureuid = ?"
		idchan1 := make(chan int, 2)
		reschan1 := make(chan bool, 2)

		db.ExecUpdate(idchan1, reschan1, updateCmd, boxuid)

		c.Data["json"] = map[string]interface{}{"code": 1}
		c.ServeJSON()
	}

}

func deleteBoxhouse(boxUid, playername string) bool {

	//在数据库中进行查找
	row, ret := db.QueryNormal("SELECT treasureboxs FROM treasurehouse WHERE useruid = (SELECT globaluid FROM accountconvert WHERE username=?)", playername)
	if ret == false || row == nil {
		fmt.Println("查询错误")
		return false
	}

	defer row.Close()

	var boxids string

	if row.Next() == false {
		fmt.Println("没有数据库返回数据111111")
		return false
	}

	err := row.Scan(&boxids)

	if err != nil {
		fmt.Println("解析失败")
		return false
	}

	boxContent := strings.Replace(boxids, boxUid+",", "", -1)

	//更新数据库
	idchan1 := make(chan int, 2)
	reschan1 := make(chan bool, 2)

	updateCmd := "UPDATE treasurehouse SET treasureboxs=? WHERE useruid= (SELECT globaluid FROM accountconvert WHERE username=?)"

	fmt.Println("更新背包物品============", boxContent)

	db.ExecUpdate(idchan1, reschan1, updateCmd,
		boxContent, playername)

	return true

}

//查看物品流程
func (c *LookBoxController) lookBox(playername string) {
	//在数据库中进行查找
	row, ret := db.QueryNormal("SELECT treasureboxs FROM treasurehouse WHERE useruid = (SELECT globaluid FROM accountconvert WHERE username=?)", playername)
	if ret == false || row == nil {
		fmt.Println("查询错误")
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
		return
	}

	defer row.Close()

	var boxids string

	if row.Next() == false {
		fmt.Println("没有数据库返回数据111111")
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
		return
	}

	err := row.Scan(&boxids)

	if err != nil {
		fmt.Println("解析失败")
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
		return
	}

	var boxuidSli []int64

	for _, v := range strings.Split(boxids, ",") {
		if v != "" {
			id, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println("背包数据的转换错误", err)
				return
			}

			boxuidSli = append(boxuidSli, int64(id))
		}
	}

	if len(boxuidSli) == 0 {
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
	}

	var boxSli []box

	for _, v := range boxuidSli {
		tempid := findBoxtempid(v)
		if tempid != -1 {
			prop := box{
				Boxuid:    strconv.Itoa(int(v)),
				Boxtempid: tempid,
			}
			boxSli = append(boxSli, prop)
		}
	}

	fmt.Println("查到该玩家的物品信息deed", boxSli)
	if len(boxSli) == 0 {
		c.Data["json"] = map[string]interface{}{"code": 0}
		c.ServeJSON()
	}

	c.Data["json"] = boxSli
	c.ServeJSON()
}

//根据一个物品的uid查找这个物品的itemid
func findBoxtempid(uid int64) (tempid int32) {
	row, ret := db.QueryNormal("SELECT treasuretempid FROM treasurebox WHERE  bdel = 0 AND treasureuid = ?", uid)
	if ret == false || row == nil {
		fmt.Println("查询错误")
		return -1
	}

	defer row.Close()

	var boxid int32

	if row.Next() == false {
		fmt.Println("没有数据库返回数据22222222", uid)
		return -1
	}

	err := row.Scan(&boxid)

	if err != nil {
		fmt.Println("解析失败")
		return -1
	}

	return boxid
}
