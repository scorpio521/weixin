package class

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"learnbeego/blog/tools"
	"time"
)

// 完成Contact类型定义
type Contacts struct {
	Id       int64 `json:"id" orm:"pk;auto;column(id)" description:"配置id"`  //设置主键自增长 列名设为id// 设置为主键，字段Id, Password首字母必须大写
	Customerpath string `json:"customerpath" orm:"size(20);column(customerpath)" description:"customer 路径"`
	Openid 	string `json:"openid" orm:"size(50);column(openid)"`
	Created_on 	time.Time `json:"created_on" orm:"column(created_on)"`
}

// 完成Contact类型定义
type Contact_answers struct {
	Id       int64 `json:"id" orm:"pk;auto;column(id)" description:"配置id"`  //设置主键自增长 列名设为id// 设置为主键，字段Id, Password首字母必须大写
	Customerpath string `json:"customerpath" orm:"size(20);column(customerpath)" description:"customer 路径"`
	Contactid 	int64 `json:"contactid" orm:"size(10);column(contactid)"`
	Aid int64 `json:"aid" orm:"column(aid)" description:"answer id"`
	Created_on 	time.Time `json:"created_on" orm:"column(created_on)"`
}

// 完成answers类型定义
type Answers struct {
	Id       int64 `json:"id" orm:"pk;auto;column(id)" description:"配置id"`  //设置主键自增长 列名设为id// 设置为主键，字段Id, Password首字母必须大写
	Customerpath string `json:"customerpath" orm:"size(20);column(customerpath)" description:"customer 路径"`
	Qid int `json:"qid" orm:"column(qid)" description:"question id"`
	t 	string `json:"t" orm:"size(500);column(t)"`
	Created_on 	time.Time `json:"created_on" orm:"column(created_on)"`
}

// 完成media类型定义
type Media struct {
	Id       int64 `json:"id" orm:"pk;auto;column(id)" description:"配置id"`  //设置主键自增长 列名设为id// 设置为主键，字段Id, Password首字母必须大写
	Customerpath string `json:"customerpath" orm:"size(20);column(customerpath)" description:"customer 路径"`
	Openid string `json:"openid" orm:"size(50);column(openid)"`
	Mediaid 	string `json:"mediaid" orm:"size(100);column(mediaid)"`
	Expiretime 	string `json:"expiretime" orm:"column(expiretime)"`
}



type resParams struct{
	Aid int64	`json:"aid" orm:"column(aid)"`
	Qid int64	`json:"qid" orm:"column(qid)"`
	Mediaid string `json:"mediaid" orm:"column(mediaid)"`
}

func GetMediaid(customerpath string,openid string) (mediaid string, err string) {//interface 接收任何形式的转换
	o := orm.NewOrm()
	res := new(resParams)
	//这个查出来就一条数据
	_, sqlerror :=o.Raw("select mediaid as value,'mediaid' as name from media where customerpath = ?  and openid = ? and expiretime>now()",customerpath,openid).RowsToStruct(res,"name","value")
	beego.Info("please test media again",sqlerror)

	if sqlerror == nil  { // 这个是查询的sql语句是正确的，没有出错，则return 出maps 数据，无论是空的还是有数据的，都是从这里出
		//beego.Info(maps[0]["num"])
		mediaid = res.Mediaid
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

func (c *Contacts) QueryByIdPath() (maps [] Contacts, err string) {
	o := orm.NewOrm()
	//这个查出来就一条数据
	_, sqlerror :=o.Raw("SELECT id FROM contacts WHERE customerpath = ? and openid = ? ", c.Customerpath,c.Openid).QueryRows(&maps)
	if sqlerror == nil { // 这个是查询的sql语句是正确的，没有出错，则return 出maps 数据，无论是空的还是有数据的，都是从这里出
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


func GetAidNum(contactid int64,customerpath string) (num interface{}, err string) {//interface 接收任何形式的转换
	o := orm.NewOrm()
	var maps []orm.Params
	//这个查出来就一条数据
	no, sqlerror :=o.Raw("select count(1) as num from (select distinct id from answers where id in (select aid from contact_answers ca where contactid = ? and customerpath = ?)) b",contactid,customerpath).Values(&maps)
	if sqlerror == nil && no>0{ // 这个是查询的sql语句是正确的，没有出错，则return 出maps 数据，无论是空的还是有数据的，都是从这里出
		//beego.Info(maps[0]["num"])
		num =maps[0]["num"]
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

func GetQidNum(customerpath string) (num interface{}, err string) {//interface 接收任何形式的转换
	o := orm.NewOrm()
	var maps []orm.Params
	//这个查出来就一条数据
	no, sqlerror :=o.Raw("select count(1) as num from questions where customerpath=? limit 1",customerpath).Values(&maps)
	if sqlerror == nil && no>0{ // 这个是查询的sql语句是正确的，没有出错，则return 出maps 数据，无论是空的还是有数据的，都是从这里出
		//beego.Info(maps[0]["num"])
		num =maps[0]["num"]
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

func GetOneQuestion(customerpath string) (maps map[string]interface{}, err string) {//interface 接收任何形式的转换
	o := orm.NewOrm()
	//这个查出来就一条数据
	//这个数据就变成了map[1:qwdawdadwa]
	//_, sqlerror :=o.Raw("select id,t from questions where customerpath= ? limit 1 ",customerpath).RowsToMap(&res)
	var res []orm.Params

	// [map[t:qwdawdadwa id:1]]
	no, sqlerror :=o.Raw("select qid as id,t from questions where customerpath= ? limit 1 ",customerpath).Values(&res)
	//io writer 还是不明白
	//f, xerr := os.Create("/tmp/ormlog")
	//if xerr != nil {
	//	fmt.Println(xerr)
	//	f.Close()
	//	return
	//}
	//d := []string{"Welcome to the world of Go1.", "Go is a compiled language.",
	//	"It is easy to learn Go."}
	//
	//for _, v := range d {
	//	fmt.Fprintln(f, v)
	//	if xerr != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//}
	//xerr = f.Close()
	//if xerr != nil {
	//	fmt.Println(xerr)
	//	return
	//}
	//fmt.Println("file written successfully")
	//
	//orm.DebugLog = orm.NewLog(f)

	if sqlerror == nil && no>0{
		// 这个是查询的sql语句是正确的，没有出错，则return 出maps 数据，有数据的才从这里出，
		// 如果是空数据的话不会从这里出，因为res 是orm.Params
		beego.Info(res[0])
		maps = res[0]
		//maps["id"] = res[0]
		//maps["t"] = res["t"]
		//num =maps[0]["num"]
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
func GetAllQuestion(customerpath string) (maps []orm.Params, err string) {//interface 接收任何形式的转换
	o := orm.NewOrm()
	//这个查出来就一条数据
	//这个数据就变成了map[1:qwdawdadwa]
	//_, sqlerror :=o.Raw("select id,t from questions where customerpath= ? limit 1 ",customerpath).RowsToMap(&res)

	// [map[t:qwdawdadwa id:1]]
	no, sqlerror :=o.Raw("select qid as id,t from questions where customerpath= ? ",customerpath).Values(&maps)

	if sqlerror == nil && no>0{
		// 这个是查询的sql语句是正确的，没有出错，则return 出maps 数据，有数据的才从这里出，
		// 如果是空数据的话不会从这里出，因为res 是orm.Params
		beego.Info(maps)
		//maps["id"] = res[0]
		//maps["t"] = res["t"]
		//num =maps[0]["num"]
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

func GetAnswerQuestionid(aid int64) (qid int64, err string) {//interface 接收任何形式的转换
	o := orm.NewOrm()
	res := new(resParams)
	//这个查出来就一条数据
	no, sqlerror :=o.Raw("select distinct a.qid as value,'qid' as name from contact_answers ca left join answers a on (ca.customerpath=a.customerpath and ca.aid=a.id) where a.id = ? ",aid).RowsToStruct(res,"name","value")
	beego.Info("please test again",no,sqlerror)

	if sqlerror == nil { // 这个是查询的sql语句是正确的，没有出错，则return 出maps 数据，无论是空的还是有数据的，都是从这里出
		//beego.Info(maps[0]["num"])
		qid = res.Qid
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

func GetObjQuestionid(contactid int64,customerpath string) (maps []orm.Params, err string) {//interface 接收任何形式的转换
	o := orm.NewOrm()
	//这个查出来就一条数据
	_, sqlerror :=o.Raw("select distinct a.qid from contact_answers ca left join answers a on (ca.customerpath=a.customerpath and ca.aid=a.id)where ca.contactid = ? and ca.customerpath= ? " ,contactid,customerpath).Values(&maps)
	if sqlerror == nil { // 这个是查询的sql语句是正确的，没有出错，则return 出maps 数据，无论是空的还是有数据的，都是从这里出
		//beego.Info(maps[0]["num"])
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

func GetOneAnswer(customerpath string,qid int) (maps []orm.Params,err string) {//interface 接收任何形式的转换
	beego.Info("qid in contacts model is :",qid)
	o := orm.NewOrm()
	//这个查出来就一条数据
	//这个数据就变成了map[1:qwdawdadwa]
	//_, sqlerror :=o.Raw("select id,t from questions where customerpath= ? limit 1 ",customerpath).RowsToMap(&res)
	// [map[t:qwdawdadwa id:1]]
	no, sqlerror :=o.Raw("select id,t from answers where customerpath= ? and qid = ? ", customerpath,qid).Values(&maps)
	if sqlerror == nil && no>0{
		// 这个是查询的sql语句是正确的，没有出错，则return 出maps 数据，有数据的才从这里出，
		// 如果是空数据的话不会从这里出，因为res 是orm.Params
		beego.Info(maps)
		//maps = res[0]
		//maps["id"] = res[0]
		//maps["t"] = res["t"]
		//num =maps[0]["num"]
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

func (c *Contacts) ReadDB() (err error) {
	o := orm.NewOrm()
	err = o.Read(c)
	return err
}

func (c *Contacts) Create() (id int64,err error) {
	o := orm.NewOrm()
	beego.Info("create success!")
	beego.Info(c)
	id, err  = o.Insert(c)
	return id,err
}


func (c *Contacts) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(c)
	return err
}

func (c *Contact_answers) Create() (id int64,err error) {
	o := orm.NewOrm()
	beego.Info("create success!")
	beego.Info(c)
	id, err  = o.Insert(c)
	return id,err
}
func (c *Media) Create() (id int64,err error) {
	o := orm.NewOrm()
	beego.Info("create success!")
	beego.Info(c)
	id, err  = o.Insert(c)
	return id,err
}