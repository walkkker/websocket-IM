package main

import (
	"./ctrl"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"text/template"
)

func RegisterView(){
	tpl, err := template.ParseGlob("view/**/*")
	if err != nil{
		log.Fatal(err.Error())
	}
	for _, v := range tpl.Templates(){
		tplName := v.Name()
		http.HandleFunc(tplName, func(writer http.ResponseWriter, request *http.Request) {
			tpl.ExecuteTemplate(writer, tplName, nil)
		})
	}
}

func main() {
	// 将前端请求与后端处理函数相绑定
	http.HandleFunc("/user/login", ctrl.UserLogin)
	http.HandleFunc("/user/register", ctrl.UserRegister) //注意只要函数名称，不要加括号
	http.HandleFunc("/contact/addfriend", ctrl.AddFriend)
	http.HandleFunc("/contact/loadfriend",ctrl.LoadFriend)
	http.HandleFunc("/contact/loadcommunity", ctrl.LoadCommunity)
	http.HandleFunc("/contact/joincommunity", ctrl.JoinCommunity)
	http.HandleFunc("/contact/createcommunity", ctrl.CreateCommunity)
	http.HandleFunc("/user/find", ctrl.UserFind)
	http.HandleFunc("/chat",ctrl.Chat)
	http.HandleFunc("/attach/upload", ctrl.Upload)

	//提供静态资源目录支持 这里要注意asset的前后都要有，尤其是后面一定要有/, 前面有因为需要匹配，后面有因为要表示后面还有url，如果不加的话，是无法访问的
	http.Handle("/asset/", http.FileServer(http.Dir(".")))
	http.Handle("/mnt/", http.FileServer(http.Dir(".")))
	http.Handle("/chat/", http.FileServer(http.Dir(".")))

	/*
	http.HandleFunc("/user/login.shtml", func(writer http.ResponseWriter, request *http.Request) {
		//解析 template，得到模板的指针
		tpl, err := template.ParseFiles("view/user/login.html")
		if err != nil{
			log.Fatal(err.Error())
		}
		tpl.ExecuteTemplate(writer, "/user/login.shtml",nil)
	}) */

	RegisterView()
	// start the web server
	http.ListenAndServe(":80",nil)
}
