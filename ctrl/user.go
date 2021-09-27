package ctrl

import (
	"../model"
	"../service"
	"../util"
	"fmt"
	"math/rand"
	"net/http"
)

var userService service.UserService
func UserRegister(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	plainpwd := request.PostForm.Get("passwd")
	nickname := fmt.Sprintf("user%06d", rand.Int31())
	avatar := ""
	sex := model.SEX_UNKNWON

	user, err := userService.Register(mobile, plainpwd, nickname, avatar, sex)
	if err != nil{
		util.RespFail(writer, err.Error())
	}else{
		util.RespOK(writer, user, "")
	}
}

func UserLogin(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	mobile := r.PostForm.Get("mobile")
	plainpwd := r.PostForm.Get("passwd")

	user, err := userService.Login(mobile, plainpwd)

	if err!=nil{
		util.RespFail(w, err.Error())
	}else{
		util.RespOK(w, user,"")
	}
}

func UserFind(w http.ResponseWriter, r *http.Request){
	/*query := r.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	userId, _ := strconv.ParseInt(id, 10, 64)
	user := userService.Find(userId)
	if (user.Token==token){
		util.RespOK(w, user, "")
	}*/
}
