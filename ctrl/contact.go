package ctrl

import (
	"../args"
	"../model"
	"../service"
	"../util"
	"net/http"
)

var contactService service.ContactService

func AddFriend(w http.ResponseWriter, r *http.Request){
	var arg args.ContactArg
	util.Bind(r,&arg)
	//调用service
	err := contactService.AddFriend(arg.Userid, arg.Dstid)
	if err!=nil{
		util.RespFail(w,err.Error())
	}else{
		util.RespOK(w,nil,"好友添加成功")
	}
}

//查找有哪些朋友 然后加载入页面中
func LoadFriend(w http.ResponseWriter, r *http.Request){
	var arg args.ContactArg
	util.Bind(r,&arg)
	users:=contactService.SearchFriend(arg.Userid)
	util.RespOKList(w,users,len(users))
}

func LoadCommunity(w http.ResponseWriter, req *http.Request){
	var arg args.ContactArg
	//如果这个用的上,那么可以直接
	util.Bind(req,&arg)
	communities := contactService.SearchComunity(arg.Userid)
	util.RespOKList(w,communities,len(communities))
}
func JoinCommunity(w http.ResponseWriter, req *http.Request){
	var arg args.ContactArg
	//如果这个用的上,那么可以直接
	util.Bind(req,&arg)
	err := contactService.JoinCommunity(arg.Userid,arg.Dstid);
	//todo 刷新用户的群组信息
	AddGroupId(arg.Userid,arg.Dstid)

	if err!=nil{
		util.RespFail(w,err.Error())
	}else {
		util.RespOK(w,nil,"")
	}
}
func CreateCommunity(w http.ResponseWriter, req *http.Request){
	var arg model.Community
	//如果这个用的上,那么可以直接
	util.Bind(req,&arg)
	com,err := contactService.CreateCommunity(arg);
	if err!=nil{
		util.RespFail(w,err.Error())
	}else {
		util.RespOK(w,com,"")
	}
}
