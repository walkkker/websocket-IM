package service

import (
	"../model"
	"../util"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

//模块化设计
type UserService struct{
}

func (s *UserService)Register(
	mobile,
	plainpwd, //明文密码
	nickname,
	avatar,
	sex string) (user model.User, err error){

	//检测手机号是否存在，如存在，则提示已经注册；否则拼接数据，插入数据库并返回新用户信息
	tmp := model.User{}
	_, err = DbEngine.Where("mobile=?",mobile).Get(&tmp)
	if err != nil{
		return tmp,err
	}
	if tmp.Id>0{
		return tmp, errors.New("该手机号已经注册")
	}
	tmp.Mobile=mobile
	tmp.Avatar=avatar
	tmp.Nickname=nickname
	tmp.Sex = sex
	tmp.Salt=fmt.Sprintf("%06d",rand.Int31n(10000))
	tmp.Passwd = util.MakePasswd(plainpwd, tmp.Salt)
	tmp.Createat = time.Now()
	//token可以是一个随机数
	tmp.Token = fmt.Sprintf("%08d",rand.Int31())

	//插入数据
	_,err = DbEngine.InsertOne(&tmp)
	//前端恶意插入特殊字符
	//数据库连接操作失败

	return tmp, err
}


func (s *UserService)Login(mobile, plainpwd string)(user model.User, err error){
	//首先通过手机号查询用户
	//查询比对密码
	//比对是否正确
	//若正确，则刷新token

	tmp := model.User{}
	DbEngine.Where("mobile=?",mobile).Get(&tmp)
	//if not find
	if tmp.Id == 0{
		return tmp, errors.New("该用户不存在")
	}
	//查询到了比对代码
	if !util.ValidatePasswd(plainpwd, tmp.Salt, tmp.Passwd){
		return tmp, errors.New("密码输入错误")
	}
	//为确保安全，刷新token
	str := fmt.Sprintf("%d", time.Now().Unix())
	token := util.MD5Encode(str)
	tmp.Token = token
	//返回数据
	DbEngine.ID(tmp.Id).Cols("token").Update(&tmp)
	return tmp, nil
}

//查询查找某个用户
func (s *UserService)Find(userId int64)(user model.User){
	tmp := model.User{}
	DbEngine.ID(userId).Get(&tmp)
	return tmp
}
