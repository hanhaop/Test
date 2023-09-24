package main

import (
	_ "demo/models"
	_ "demo/routers"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	"strconv"
)

// step2；在main函数中增加
func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//连接数据库
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/newsWeb?charset=utf8")
}
func main() {

	beego.AddFuncMap("ShowPrePage", HandlePrepage)
	beego.AddFuncMap("ShowNextPage", HandleNextpage)
	beego.Run()
}

func HandlePrepage(data int) string {
	pageIndex := data - 1
	if pageIndex < 1 {

	}
	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}
func HandleNextpage(data int) string {
	pageIndex := data + 1
	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}
