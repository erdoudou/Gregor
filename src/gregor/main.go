package main

import (
	"Public/db"
	"gregor/controll"
	"gregor/roleMgr"
	"log"
	"net/http"
)

//首先访问的网址为：http://localhost:8888/register/register
func main() {
	log.Println("main")

	db.Init("mysql", "root", "root", "gregor", "127.0.0.1:3306")

	roleMgr.GetAllRole()

	loadStaticFile()
	loadHtmlFile()

	http.ListenAndServe(":8888", nil)
}

//加载静态文件
func loadStaticFile() {
	http.Handle("/css/", http.FileServer(http.Dir("../../template")))
	http.Handle("/js/", http.FileServer(http.Dir("../../template")))
	http.Handle("/images/", http.FileServer(http.Dir("../../template")))
}

//加载htmll页面
func loadHtmlFile() {
	http.HandleFunc("/home/", controll.HomeHandler)

	http.HandleFunc("/register/", controll.RegisterHandler)

	http.HandleFunc("/admin/", controll.AdminHandler)

	http.HandleFunc("/login/", controll.LoginHandler)

	http.HandleFunc("/ajax/", controll.AjaxHandler)

	http.HandleFunc("/", controll.NotFoundHandler)

	http.HandleFunc("/test", controll.TestHandler)
}
