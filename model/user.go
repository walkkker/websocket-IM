package model

import "time"

const(
	SEX_WOMEN="W"
	SEX_MAN="M"
	SEX_UNKNWON="U"
)

//model.SEX_WOMEN
type User struct {
	//用户ID
	Id         int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	//手机号码
	Mobile   string 		`xorm:"varchar(20)" form:"mobile" json:"mobile"`
	//用户密码=f(plainwd+salt) f加密函数-MD5
	Passwd       string	`xorm:"varchar(40)" form:"passwd" json:"-"`
	//头像
	Avatar	   string 		`xorm:"varchar(150)" form:"avatar" json:"avatar"`
	Sex        string	`xorm:"varchar(2)" form:"sex" json:"sex"`
	Nickname    string	`xorm:"varchar(20)" form:"nickname" json:"nickname"`
	//加盐随机字符串6
	Salt       string	`xorm:"varchar(10)" form:"salt" json:"-"`
	Online     int	`xorm:"int(10)" form:"online" json:"online"`
	//前端鉴权因子, /chat?id=1&token=x
	Token      string	`xorm:"varchar(40)" form:"token" json:"token"`
	Memo      string	`xorm:"varchar(140)" form:"memo" json:"memo"`
	//用于统计每天用户增量
	Createat   time.Time	`xorm:"datetime" form:"createatTime" json:"createat"`
}


