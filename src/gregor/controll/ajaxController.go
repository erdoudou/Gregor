package controll

import (
	"encoding/json"
	"fmt"
	//"log"
	"net/http"

	"Public/db"

	"database/sql"

	"gregor/roleMgr"

	//"html/template"
)

type Result struct {
	Ret    int
	Reason string
	Data   interface{}
}

type ajaxController struct {
}

func (this *ajaxController) LoginAction(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")

	err := r.ParseForm()

	if err != nil {
		OutputJson(w, 0, "error canshucuowu", nil)
		return
	}

	admin_name := r.FormValue("admin_name")
	admin_password := r.FormValue("admin_password")

	if admin_name == "" || admin_password == "" {
		OutputJson(w, 0, "error uesr or password err", nil)
		return
	}

	//判断一个玩家是否存在数据库中

	row1 := make(chan *sql.Rows, 1)
	db.Query(row1, "select username,password from account where username = ?", admin_name)
	in1 := <-row1

	if in1 == nil {

		return
	}

	defer in1.Close()

	var name, password string

	if in1.Next() != false {

		err := in1.Scan(&name, &password)

		if err != nil {
			fmt.Println("4444444444444444", err)
			return
		}

	}

	admin_password_db := password

	if admin_password_db != admin_password {
		OutputJson(w, 0, "密码输入错误", nil)
		return
	}

	// 存入cookie,使用cookie存储
	cookie := http.Cookie{Name: "admin_name", Value: name, Path: "/"}
	http.SetCookie(w, &cookie)
	fmt.Println("上传网页")

	OutputJson(w, 2, "successfull!", nil)

	return
}

//注册登录流程
func (this *ajaxController) RegisterAction(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")

	err := r.ParseForm()

	if err != nil {
		OutputJson(w, 0, "error canshucuowu", nil)
		return
	}

	admin_name := r.FormValue("admin_name")
	admin_password := r.FormValue("admin_password")

	if admin_name == "" || admin_password == "" {
		OutputJson(w, 0, "error uesr or password err", nil)
		return
	}

	isRole := roleMgr.IsRole(admin_name)
	if isRole == true {
		//玩家在数据库中验证密码
		ispwd := roleMgr.IsPwd(admin_name, admin_password)
		if ispwd == false {
			OutputJson(w, 0, "error password err", nil)
			return
		}
	} else {
		//玩家不存在数据库中，将玩家数据存入数据库
		roleMgr.SetRole(admin_name, admin_password)
	}

	// 存入cookie,使用cookie存储
	cookie := http.Cookie{Name: "admin_name", Value: admin_name, Path: "/"}
	http.SetCookie(w, &cookie)
	fmt.Println("上传网页")

	OutputJson(w, 2, "successfull!", nil)

	return
}

func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &Result{ret, reason, i}
	fmt.Println(out)

	b, err := json.Marshal(out)
	if err != nil {
		return
	}

	w.Write(b)
}
