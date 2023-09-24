package controllers

import (
	"demo/models"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"time"
)

type RegController struct {
	beego.Controller
}

func (this *RegController) ShowReg() {
	this.TplName = "register.html"
}

/*
1、拿浏览器传递的数据
2、数据处理
3、插入数据库
4、返回视图
*/

func (this *RegController) HandleReg() {
	//拿浏览器传递的数据
	name := this.GetString("userName")
	passwd := this.GetString("password")

	//数据处理
	if name == "" || passwd == "" {
		logs.Info("用户名或密码不可为空")
		this.TplName = "register.html"
		return
	}
	//插入数据库
	//获取orm对象、获取插入对象、插入操作、返回
	o := orm.NewOrm()
	user := models.User{}
	user.UserName = name
	user.Passwd = passwd
	_, err := o.Insert(&user)
	if err != nil {
		logs.Info("插入失败")
	}
	//this.TplName="login.html" //会造成状态栏路径错误
	//this.Ctx.WriteString("注册成功")
	//重定向
	this.Redirect("/", 302)
}

type LoginController struct {
	beego.Controller
}

func (this *LoginController) ShowLogin() {
	name := this.Ctx.GetCookie("userName")
	if name != "" {
		this.Data["userName"] = name
		this.Data["check"] = "checked"
	}
	this.TplName = "login.html"
}

func (this *LoginController) HandleLogin() {
	name := this.GetString("userName")
	passwd := this.GetString("password")
	//logs.Info(name, passwd)
	if name == "" || passwd == "" {
		logs.Info("用户名或密码不可为空")
		this.TplName = "login.html"
		return
	}
	//查找数据
	o := orm.NewOrm()
	user := models.User{}
	user.UserName = name
	err := o.Read(&user, "UserName")
	if err != nil {
		logs.Info("用户名失败")

	}
	//判断密码是否正确
	if user.Passwd != passwd {
		logs.Info("密码失败")
		this.TplName = "login.html"
		return
	}
	//返回
	//记住用户名
	check := this.GetString("remember")
	if check == "on" {
		this.Ctx.SetCookie("userName", name, time.Second*3600)
	} else {
		this.Ctx.SetCookie("userName", "name", -1)
	}

	this.SetSession("userName", name)

	//this.Ctx.WriteString("登陆成功")
	this.Redirect("/Article/ShowArticle", 302)

}
