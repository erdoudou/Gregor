package controll

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type registerController struct {
}

func (this *registerController) RegisterAction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("加载登录注册界面")
	t, err := template.ParseFiles("../../template/html/register/register.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}
