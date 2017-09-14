package controll

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	UserName string
}

type adminController struct {
}

func (this *adminController) IndexAction(w http.ResponseWriter, r *http.Request, user string) {
	fmt.Println("加载后台管理系统界面")
	t, err := template.ParseFiles("../../template/html/admin/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, &User{user})
}
