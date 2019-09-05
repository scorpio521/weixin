package main

import (
	_ "learnbeego/blog/models"
	"github.com/astaxie/beego/logs"
	_ "learnbeego/blog/routers"
	"github.com/astaxie/beego"
	 "encoding/json"
    "fmt"
)

func main() {
	initLogger()
	beego.Run()
}

func initLogger()(err error) {

    config := make(map[string]interface{})
    config["filename"] = beego.AppConfig.String("log_path")

    // map 转 json
    configStr, err := json.Marshal(config)
    if err != nil {
        fmt.Println("initLogger failed, marshal err:", err)
        return
    }
    // log 的配置
    beego.SetLogger(logs.AdapterFile, string(configStr))
    // log打印文件名和行数
    beego.SetLogFuncCall(true)
    fmt.Println(string(configStr))
    return
}
