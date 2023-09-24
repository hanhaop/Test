package routers

import (
	"demo/controllers"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {

	beego.InsertFilter("/Article/*", beego.BeforeRouter, FilterFunc)

	//beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.RegController{}, "get:ShowReg;Post:HandleReg")
	beego.Router("/", &controllers.LoginController{}, "get:ShowLogin;Post:HandleLogin")
	beego.Router("/Article/ShowArticle", &controllers.ArticleController{}, "get:ShowArticleList;post:HandleSelect")
	beego.Router("/Article/AddArticle", &controllers.ArticleController{}, "get:ShowAddArticle;Post:HandleAddArticle")
	beego.Router("ShowContent", &controllers.ArticleController{}, "get:ShowContent")
	beego.Router("/Article/DeleteArticle", &controllers.ArticleController{}, "get:HandleDelete")
	beego.Router("/Article/UpdateArticle", &controllers.ArticleController{}, "get:ShowUpdate;Post:HandleUpdate")
	beego.Router("/Article/AddArticleType", &controllers.ArticleController{}, "get:ShowAddType;Post:HandleAddType")
	beego.Router("/Article/Logout", &controllers.ArticleController{}, "get:Logout")
}

var FilterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/")
	}
}
