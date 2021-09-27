package service

import (
	"../model"
	"errors"
	"time"
)

type ContactService struct{

}

//自动添加好友
func (service *ContactService) AddFriend(
	userid, destid int64) error{
	//第一个限制 不能加自己
	if userid == destid{
		return errors.New("不能添加自己")
	}

	//第二个限制 对方用户不存在
	tmp_user := model.User{}
	DbEngine.Where("id=?",destid).Get(&tmp_user)
	if tmp_user.Id == 0{
		return errors.New("对方用户不存在")
	}

	//第三个限制 判断是否已经添加对方
	tmp := model.Contact{}
	//条件的链式操作
	DbEngine.Where("ownerid=?", userid).
		And("destid=?",destid).
		And("cate=?",model.CONCAT_CATE_USER).
		Get(&tmp)

	//若存在记录则说明不是好友
	if tmp.Id>0{
		return errors.New("该用户已经被添加")
	}

	//使用事务进行添加
	session := DbEngine.NewSession()
	session.Begin()
	//向自己的记录插入
	_,e2 := session.InsertOne(model.Contact{
		Ownerid: userid,
		Dstobj: destid,
		Cate: model.CONCAT_CATE_USER,
		Createat: time.Now(),
	})
	//向对方记录插入
	_,e3 := session.InsertOne(model.Contact{
		Ownerid: destid,
		Dstobj: userid,
		Cate: model.CONCAT_CATE_USER,
		Createat: time.Now(),
	})
	//没有错误
	if e2==nil && e3==nil{
		//可以对事务进行提交
		session.Commit()
		return nil
	}else{
		//有问题存在
		session.Rollback()
		if e2!=nil{
			return e2
		}else{
			return e3
		}
	}
}

//查找有哪些朋友
func (service *ContactService) SearchFriend(userId int64) ([]model.User){
	contacts := make([]model.Contact,0)
	objIds := make([]int64,0)
	//Find 遍历 表
	DbEngine.Where("ownerid=? and cate=?", userId, model.CONCAT_CATE_USER).Find(&contacts)
	for _,v := range contacts{
		objIds = append(objIds, v.Dstobj)
	}
	frids := make([]model.User,0)
	if len(objIds)==0{
		return frids
	}
	DbEngine.In("id",objIds).Find(&frids)
	return frids
}

//建群
func (service *ContactService) CreateCommunity(comm model.Community) (ret model.Community,err error){
	if len(comm.Name)==0{
		err = errors.New("缺少群名称")
		return ret,err
	}
	if comm.Ownerid==0{
		err = errors.New("请先登录")
		return ret,err
	}
	com := model.Community{
		Ownerid:comm.Ownerid,
	}
	num,err := DbEngine.Count(&com)

	if(num>5){
		err = errors.New("一个用户最多只能创见5个群")
		return com,err
	}else{
		comm.Createat=time.Now()
		session := DbEngine.NewSession()
		session.Begin()
		_,err = session.InsertOne(&comm)
		if err!=nil{
			session.Rollback();
			return com,err
		}
		_,err =session.InsertOne(
			model.Contact{
				Ownerid:comm.Ownerid,
				Dstobj:comm.Id,
				Cate:model.CONCAT_CATE_COMUNITY,
				Createat:time.Now(),
			})
		if err!=nil{
			session.Rollback();
		}else{
			session.Commit()
		}
		return com,err
	}
}

//加群
func (service *ContactService) JoinCommunity(userId,comId int64) error{
	cot := model.Contact{
		Ownerid:userId,
		Dstobj:comId,
		Cate:model.CONCAT_CATE_COMUNITY,
	}
	DbEngine.Get(&cot)
	if(cot.Id==0){
		cot.Createat = time.Now()
		_,err := DbEngine.InsertOne(cot)
		return err
	}else{
		return nil
	}
}

func (service *ContactService) SearchComunity(userId int64) ([]model.Community){
	conconts := make([]model.Contact,0)
	comIds :=make([]int64,0)

	DbEngine.Where("ownerid = ? and cate = ?",userId,model.CONCAT_CATE_COMUNITY).Find(&conconts)
	for _,v := range conconts{
		comIds = append(comIds,v.Dstobj);
	}
	coms := make([]model.Community,0)
	if len(comIds)== 0{
		return coms
	}
	DbEngine.In("id",comIds).Find(&coms)
	return coms
}

func (service *ContactService) SearchComunityIds(userId int64) (comIds []int64){
	//todo 获取用户全部群ID
	conconts := make([]model.Contact,0)
	comIds =make([]int64,0)

	DbEngine.Where("ownerid = ? and cate = ?",userId,model.CONCAT_CATE_COMUNITY).Find(&conconts)
	for _,v := range conconts{
		comIds = append(comIds,v.Dstobj);
	}
	return comIds
}
