package controll

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type loginController struct {
}

func (this *loginController) IndexAction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("加载登录界面")
	t, err := template.ParseFiles("../../template/html/login/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

//oculus源码
