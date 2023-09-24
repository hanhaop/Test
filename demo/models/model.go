package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

import _ "github.com/go-sql-driver/mysql"

// 放入结构体和初始化语句

type User struct {
	Id       int
	UserName string `orm:"unique"`
	Passwd   string
	Articles []*Article `orm:"reverse(many)"`
}

// 文章表和文章类型表 n：1
type Article struct {
	Id          int          `orm:"pk;auto"`       //主键且自增
	Title       string       `orm:"size(20)"`      //文章标题
	Content     string       `orm:"size(500)"`     //内容
	Img         string       `orm:"size(50);null"` //图片
	Type        string       //类型
	Time        time.Time    `orm:"type(datetima);auto_now_add"` //发布时间
	Count       int          `orm:"default(0)"`                  //阅读量
	ArticleType *ArticleType `orm:"rel(fk)"`                     //外键
	Users       []*User      `orm:"rel(m2m)"`
}

type ArticleType struct {
	Id       int
	TypeName string     `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"` //与fk对应,1:n
}

// step1:在构建多对多关系的.go文件中
func init() {
	//注册表
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	//生成表
	//orm.RunSyncdb("default", false, true)
}
