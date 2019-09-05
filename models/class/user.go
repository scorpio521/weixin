package class

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

// 完成User类型定义
type User struct {
	Id       int `orm:"pk"` // 设置为主键，字段Id, Password首字母必须大写
	Password string
}

func (u *User) ReadDB() (err error) {
	o := orm.NewOrm()
	err = o.Read(u)
	return err
}

func (u *User) Create() (err error) {
	o := orm.NewOrm()
	fmt.Println("Create success!")
	_, _ = o.Insert(u)
	return err
}

func (u *User) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(u)
	return err
}