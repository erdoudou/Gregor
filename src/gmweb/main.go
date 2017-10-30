package main

import (
	_ "gergorWeb/routers"

	"github.com/astaxie/beego"

	"Public/db"
	"fmt"
	"gergorWeb/box"
	"gergorWeb/item"
)

func main() {
	beego.SetStaticPath("/static/images", "images")
	beego.SetStaticPath("/static/css", "css")
	beego.SetStaticPath("/static/js", "js")
	beego.BConfig.WebConfig.Session.SessionOn = true
	//连接数据库
	db.Init("mysql", "root", "root", "helldivers", "192.168.2.20:3306")

	retBox := box.LoadCfg()
	retProp := item.LoadCfg()
	if retBox == false || retProp == false {
		fmt.Println("配置文件有错误请检查配置文件")
	}

	beego.Run()
}
