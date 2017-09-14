package controll

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	// 获取cookie
	fmt.Println("进入后台管理界面")
	cookie, err := r.Cookie("admin_name")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login/index", http.StatusFound)
	}

	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	admin := &adminController{}
	controller := reflect.ValueOf(admin)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		fmt.Println("3333333333333")
		method = controller.MethodByName(strings.Title("index") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	userValue := reflect.ValueOf(cookie.Value)
	method.Call([]reflect.Value{responseValue, requestValue, userValue})
}

func AjaxHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("进入页面跳转流程")
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	ajax := &ajaxController{}
	controller := reflect.ValueOf(ajax)
	method := controller.MethodByName(action)
	fmt.Println("打印出action的值:", action)
	if !method.IsValid() {
		fmt.Println("未发现应该跳转的页面", strings.Title("index")+"Action")
		method = controller.MethodByName(strings.Title("index") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("跳转到登录界面")

	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	login := &loginController{}

	controller := reflect.ValueOf(login)

	method := controller.MethodByName(action)
	if !method.IsValid() {
		fmt.Println("登录页面的if判断句中")
		method = controller.MethodByName(strings.Title("index") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)

	method.Call([]reflect.Value{responseValue, requestValue})
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("进入未发现页面的登录流程")
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/login/index", http.StatusFound)
	}

	t, err := template.ParseFiles("../../template/html/404.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("跳转到测试界面")

	t, err := template.ParseFiles("../../template/html/test.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("跳转到注册页面")

	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	register := &registerController{}

	controller := reflect.ValueOf(register)

	method := controller.MethodByName(action)
	if !method.IsValid() {
		fmt.Println("111111111111")

		method = controller.MethodByName(strings.Title("register") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)

	method.Call([]reflect.Value{responseValue, requestValue})
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("跳转到主页面")

	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}

	home := &homeController{}

	controller := reflect.ValueOf(home)

	method := controller.MethodByName(action)
	if !method.IsValid() {
		fmt.Println("111111111111")

		method = controller.MethodByName(strings.Title("home") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)

	method.Call([]reflect.Value{responseValue, requestValue})
}
