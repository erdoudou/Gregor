package main

import (
	"Public/db"
	"gregor/controll"
	"gregor/roleMgr"
	"log"
	"net/http"
)

func main() {
	log.Println("main")

	db.Init("mysql", "root", "root", "gun", "192.168.2.20:3306")
	roleMgr.GetAllRole()

	http.Handle("/css/", http.FileServer(http.Dir("../../template")))
	http.Handle("/js/", http.FileServer(http.Dir("../../template")))
	http.Handle("/images/", http.FileServer(http.Dir("../../template")))

	http.HandleFunc("/home/", controll.HomeHandler)

	http.HandleFunc("/register/", controll.RegisterHandler)

	http.HandleFunc("/admin/", controll.AdminHandler)

	http.HandleFunc("/login/", controll.LoginHandler)

	http.HandleFunc("/ajax/", controll.AjaxHandler)

	http.HandleFunc("/", controll.NotFoundHandler)

	http.HandleFunc("/test", controll.TestHandler)

	http.ListenAndServe(":8888", nil)
}
