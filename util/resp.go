package util

import (
	"encoding/json"
	"log"
	"net/http"
)

// 定义结构体 为Resp函数使用
type H struct{
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
	Rows interface{} `json:"rows,omitempty"`
	Total interface{} `json:"total,omitempty"`
}

// 返回JSON语句
func Resp(w http.ResponseWriter, code int, data interface{}, msg string){
	//设置header为JSON 默认的text/html,所以特别指出返回的为application/json
	w.Header().Set("Content-Type", "application/json")
	//设置200状态
	w.WriteHeader(http.StatusOK)

	// 定义一个结构体,将接收的参数放入定义的结构体内
	h := H{
		Code:code,
		Msg:msg,
		Data:data,
	}
	// 将结构体转化为JSON字符串 非文本格式到文本格式
	ret,err := json.Marshal(h)
	// 如果有错误的话 err的值不为空
	if err != nil{
		log.Println(err.Error())
	}
	// 输出
	w.Write(ret)
}

func RespFail(w http.ResponseWriter, msg string){
	Resp(w, -1, nil, msg)
}

func RespOK(w http.ResponseWriter, data interface{}, msg string){
	Resp(w, 0, data, msg)
}

func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}){
	w.Header().Set("Content-Type","application/json")
	//设置200状态
	w.WriteHeader(http.StatusOK)
	/*
	输出 定义一个结构体 满足某一条件的全部记录数目
	*/
	h := H{
		Code:code,
		Rows:data,
		Total:total,
	}
	//将结构体转化为json字符串
	ret, err := json.Marshal(h)
	if err != nil{
		log.Println(err.Error())
	}
	w.Write(ret)
}

func RespOKList(w http.ResponseWriter, lists interface{}, total interface{}){
	//分页数目
	RespList(w,0,lists,total)
}