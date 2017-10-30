package main

import (
	_ "gergorWeb/routers"

	"github.com/astaxie/beego"

	"Public/db"
)

func main() {
	beego.SetStaticPath("/static/img", "images")
	beego.SetStaticPath("/static/css", "css")
	beego.SetStaticPath("/static/js", "js")
	beego.BConfig.WebConfig.Session.SessionOn = true
	//连接数据库
	db.Init("mysql", "root", "root", "gregor", "127.0.0.1:3306")

	beego.Run()
}
