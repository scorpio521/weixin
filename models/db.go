package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"learnbeego/blog/models/class" // 注册模型，需要引入该包
)

/*
使用orm连接数据库步骤：
//告诉orm使用哪一种数据库
1.注册数据库驱动RegisterDriver(driverName, DriverType)
2.注册数据库RegisterDataBase(aliasName, driverName, dataSource, params ...)
3.注册对象模型RegisterModel(models ...)
4.开启同步RunSyncdb(name string, force bool, verbose bool)
*/

// 在init函数中连接数据库，当导入该包的时候便执行此函数
func init(){
	orm.RegisterDriver("mysql", orm.DRMySQL)
	if beego.AppConfig.String("runmode")=="dev" {
		orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("mysqldev::user")+":"+beego.AppConfig.String("mysqldev::pwd")+"@tcp("+beego.AppConfig.String("mysqldev::host")+":"+beego.AppConfig.String("mysqldev::port")+")/"+beego.AppConfig.String("mysqldev::db")+"?charset=utf8")
	}else{
		orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("mysql::user")+":"+beego.AppConfig.String("mysql::pwd")+"@tcp("+beego.AppConfig.String("mysql::host")+":"+beego.AppConfig.String("mysql::port")+")/"+beego.AppConfig.String("mysql::db")+"?charset=utf8")
	}

	orm.RegisterModel(new(class.Config),new(class.Contacts),new(class.Users),new(class.Contact_answers),new(class.Media)) // 注册模型，建立User类型对象，注册模型时，需要引入包
	orm.RunSyncdb("default", false, true)
	orm.Debug = true
}