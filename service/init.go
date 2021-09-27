package service

import (
	"../model"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

var DbEngine *xorm.Engine
func init(){
	driverName := "mysql"
	DsName := "openmp:a0154544956a@(127.0.0.1:3306)/chat?charset=utf8"
	err := errors.New("")
	DbEngine, err = xorm.NewEngine(driverName, DsName)
	if err != nil && err.Error()!= ""{
		log.Fatal(err.Error())
	}
	//是否显示SQL语句
	DbEngine.ShowSQL(true)
	//设置数据库打开的最大连接数
	DbEngine.SetMaxOpenConns(2)
	//自动建表 将model中的结构体与数据库中的表结构建立起一一映射的关系
	DbEngine.Sync2(new(model.User), new(model.Contact), new(model.Community))
	fmt.Println("init database ok")
}

