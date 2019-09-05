package class

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"learnbeego/blog/tools"
)



// 完成User类型定义
type Config struct {
	Id       int64 `json:"id" orm:"pk;auto;column(id)" description:"配置id"`  //设置主键自增长 列名设为id// 设置为主键，字段Id, Password首字母必须大写
	Customerpath string `json:"customerpath" orm:"size(20);column(customerpath)" description:"customer 路径"`
	Appid 	string `json:"appid" orm:"size(20);column(appid)"`
	Appsec 	string `json:"appsec" orm:"size(20);column(appsec)"`
}
func (c *Config) QueryByPath() (maps []Config, err string) {
	o := orm.NewOrm()
	_, sqlerr :=o.Raw("SELECT * FROM config WHERE customerpath = ? ", c.Customerpath).QueryRows(&maps)
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
func (c *Config) ReadDB() (err error) {
	o := orm.NewOrm()
	err = o.Read(c)
	return err
}

func (c *Config) Create() (id int64,err error) {
	o := orm.NewOrm()
	fmt.Println("Create success!")
	id, err = o.Insert(c)
	return id,err
}

func (c *Config) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(c)
	return err
}