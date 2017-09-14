package controll

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type homeController struct {
}

func (this *homeController) HomeAction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("加载主页面界面")
	t, err := template.ParseFiles("../../template/html/home/home.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}
