package main

import (
	_ "gmweb/routers"

	"github.com/astaxie/beego"

	"Public/db"
	"fmt"
	"gmweb/box"
	"gmweb/item"
	"runtime"
)

var servCfg ServerCfg

type ServerCfg struct {
	UserName  string
	Pwd       string
	IpAndPort string
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	item.Load("config/servercfg.txt", &servCfg)
	retBox := box.LoadCfg()
	retProp := item.LoadCfg()
	if retBox == false || retProp == false {
		fmt.Println("配置文件有错误请检查配置文件")
	}
}

func main() {
	beego.SetStaticPath("/static/images", "images")
	beego.SetStaticPath("/static/css", "css")
	beego.SetStaticPath("/static/js", "js")
	beego.BConfig.WebConfig.Session.SessionOn = true
	//连接数据库
	db.Init("mysql", servCfg.UserName, servCfg.Pwd, "helldivers", servCfg.IpAndPort)

	beego.Run()
}
