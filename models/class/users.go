package class

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"learnbeego/blog/tools"
	"time"
)



// 完成User类型定义
type Users struct {
	Id       int64 `json:"id" orm:"pk;auto;column(id)" description:"配置id"`  //设置主键自增长 列名设为id// 设置为主键，字段Id, Password首字母必须大写
	Customerpath string `json:"customerpath" orm:"size(20);column(customerpath)" description:"customer 路径"`
	Contactid string `json:"contactid" orm:"column(contactid)" description:"客户id"`
	Unique_openid 	string `json:"unique_openid" orm:"size(50);column(unique_openid)"`
	Created_on 	time.Time `json:"created_on" orm:"column(created_on)"`
}
func (c *Users) QueryByUnique() (maps []Users, err string) {
	o := orm.NewOrm()
	_, sqlerr :=o.Raw("SELECT * FROM users WHERE customerpath = ? and unique_openid = ? ", c.Customerpath,c.Unique_openid).QueryRows(&maps)
	if sqlerr == nil {// 这个是查询的sql语句是正确的，没有出错，则return 出maps 数据，无论是空的还是有数据的，都是从这里出
		return
	}
	// 这里return 出去的是error错误处理,比如sql 语句出错了，有点类似try catch 的用法
	de := &tools.DivideError{
		Code: 1,
		Msg: "the sql sentence is wrong,please check",
	}
	err = de.Error()
	return
}
func (c *Users) ReadDB() (err error) {
	o := orm.NewOrm()
	err = o.Read(c)
	return err
}

func (c *Users) Create() (id int64,err error) {
	o := orm.NewOrm()
	fmt.Println("Create success!")
	id, err = o.Insert(c)
	return id,err
}

func (c *Users) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(c)
	return err
}