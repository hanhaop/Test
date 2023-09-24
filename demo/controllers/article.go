package controllers

import (
	"demo/models"
	"fmt"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"math"
	"path"
	"strconv"
	"time"
)

type ArticleController struct {
	beego.Controller
}

func (this *ArticleController) ShowArticleList() {

	//删除该session过滤后可直接通过URL进入系统内容
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/", 302)
		return
	}

	this.TplName = "index.html"
	o := orm.NewOrm()
	qs := o.QueryTable("Article") //文章查询器，只创建但不查询
	var articles []models.Article
	//qs.All(&articles)
	//获取部分数据
	pageIndex := this.GetString("pageIndex")
	pageIndex1, err := strconv.Atoi(pageIndex)
	if err != nil {
		pageIndex1 = 1
	}
	count, err := qs.RelatedSel("ArticleType").Count()
	if err != nil {
		logs.Info("获取记录数错误：", err)
		return
	}
	pageSize := 2
	start := pageSize * (pageIndex1 - 1)
	qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)

	pageCount := float64(count) / float64(pageSize)
	pageCount1 := math.Ceil(pageCount)
	if err != nil {
		logs.Info("查询错误")
		return
	}

	//首页末页数据处理
	FirstPage := false
	if pageIndex1 == 1 {
		FirstPage = true
	}
	LastPage := false
	if float64(pageIndex1) == pageCount1 {
		LastPage = true
	}

	//获取类型数据
	var types []models.ArticleType
	var articles_Type []models.Article

	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types
	typeName := this.GetString("select")
	logs.Info(typeName + "+++++++++++++++++++")
	if typeName == "" {
		logs.Info("下拉框传递数据失败")
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles_Type)
	} else {
		count, err = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
		if err != nil {
			fmt.Println("获取记录数错误：", err)
			return
		}
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles_Type)
	}

	this.Data["userName"] = userName
	logs.Info("count=", count)
	this.Data["typeName"] = typeName
	this.Data["FirstPage"] = FirstPage
	this.Data["LastPage"] = LastPage
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount1
	this.Data["pageIndex"] = pageIndex1
	this.Data["articles"] = articles_Type

	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["contentHead"] = "ShowArticleHead.html"
	this.TplName = "index.html"
}

func (this *ArticleController) HandleSelect() {
	typeName := this.GetString("select")
	var articles []models.Article
	if typeName == "" {
		logs.Info("下拉框传递数据失败")
		return
	}
	o := orm.NewOrm()
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)
	this.Data["articles"] = articles

	articleTypes := []models.ArticleType{}
	o.QueryTable("ArticleType").All(articleTypes)
	this.Data["articleTypes"] = articleTypes
	this.TplName = "index.html"

}

func (this *ArticleController) ShowAddArticle() {
	//获取类型数据
	var types []models.ArticleType
	o := orm.NewOrm()
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types
	this.TplName = "add.html"
}

func (this *ArticleController) HandleAddArticle() {
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	file, header, err := this.GetFile("uploadname")
	defer file.Close()
	//判断文件格式
	ext := path.Ext(header.Filename)
	logs.Info(ext)
	if ext != ".jpg" && ext != ".png" && ext != ".jepg" {
		logs.Info("文件格式不正确")
		return
	}
	//判断文件大小
	if header.Size > 50000000 {
		logs.Info("文件过大，无法上传")
		return
	}
	//是否重名
	fileName := time.Now().Format("2006-01-02 15-04-05")
	this.SaveToFile("uploadname", "./static/img/"+fileName+ext)

	if err != nil {
		logs.Info("文件上传失败")
		return
	}
	//logs.Info(fileName + ext)
	//logs.Info(articleName, content)

	//数据处理
	o := orm.NewOrm()
	Article := models.Article{}
	Article.Title = articleName
	Article.Content = content
	Article.Img = "./static/img/" + fileName + ext

	//获取下拉框传递的数据
	typeName := this.GetString("select")
	if typeName == "" {
		logs.Info("下拉框数据获取失败" + "--------------------")
		return
	}
	var artiType models.ArticleType
	artiType.TypeName = typeName
	err = o.Read(&artiType, "TypeName")
	if err != nil {
		logs.Info("获取类型错误")
		return
	}
	Article.ArticleType = &artiType

	_, err = o.Insert(&Article)
	if err != nil {
		logs.Info("插入失败", err)
		return
	}

	this.Redirect("/Article/ShowArticle", 302)
}

func (this *ArticleController) ShowContent() {
	//获取文章ID
	id, err := this.GetInt("id")
	//logs.Info("id is ", id)
	if err != nil {
		logs.Info("获取文章ID错误")
		return
	}
	//查询数据库获取数据
	o := orm.NewOrm()
	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		logs.Info("查询错误", err)
		return
	}
	arti.Count += 1
	//多对多插入读者
	//1.获取操作对象
	//article := models.Article{Id: id}
	//2.获取多对多操作对象
	m2m := o.QueryM2M(&arti, "Users")
	userName := this.GetSession("userName")
	user := models.User{}
	user.UserName = userName.(string)
	o.Read(&user, "UserName")
	//4.多对多插入
	_, err = m2m.Add(&user)
	if err != nil {
		logs.Info("插入失败")
		return
	}

	o.Update(&arti)

	//o.LoadRelated(&arti, "Users")//会有重复

	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id", id).Distinct().All(&users)
	logs.Info(arti)
	this.Data["users"] = users
	this.Data["article"] = arti
	logs.Info(arti.Img)
	//传递数据给视图
	this.Layout = "layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["contentHead"] = "head.html"
	this.TplName = "content.html"
}

func (this *ArticleController) HandleDelete() {
	id, err := this.GetInt("id")
	if err != nil {
		logs.Info("获取文章id错误")
		return
	}
	o := orm.NewOrm()
	arti := models.Article{Id: id}
	o.Delete(&arti)
	this.Redirect("/Article/ShowArticle", 302)
}

func (this *ArticleController) ShowUpdate() {
	//获取文章ID
	id, err := this.GetInt("id")
	//logs.Info("id is ", id)
	if err != nil {
		logs.Info("获取文章ID错误")
		return
	}
	//查询数据库获取数据
	o := orm.NewOrm()
	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		logs.Info("查询错误", err)
		return
	}
	this.Data["article"] = arti
	this.Layout = "layout.html"
	this.TplName = "update.html"
}

func (this *ArticleController) HandleUpdate() {

	var filename string

	id, errr := this.GetInt("id")
	if errr != nil {
		logs.Info("获取文章id失败")
		return
	}
	name := this.GetString("articleName")
	content := this.GetString("content")
	if name == "" || content == "" {
		logs.Info("更新数据失败")
		return
	}
	f, h, err := this.GetFile("uploadname")
	if err != nil {
		logs.Info("上传文件失败")
	} else {
		ext := path.Ext(h.Filename)
		if ext != ".jpg" && ext != ".jepg" && ext != ".png" {
			logs.Info("文件格式错误，无法上传")
			return
		}
		if h.Size > 5000000 {
			logs.Info("照片过大，无法上传")
			return
		}
		filename = time.Now().Format("2006-01-02 15-04-05") + ext
		err = this.SaveToFile("uploadname", "./static/img/"+filename)
		if err != nil {
			logs.Info("文件保存失败", err)
			return
		}
		defer f.Close()
	}

	o := orm.NewOrm()
	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		logs.Info("文章不存在")
		return
	}
	arti.Title = name
	arti.Content = content
	if filename != "" {
		arti.Img = "./static/img/" + filename
	}

	_, err = o.Update(&arti)
	if err != nil {
		logs.Info("更新失败")
		return
	}
	this.Redirect("/Article/ShowArticle", 302)
}

func (this *ArticleController) ShowAddType() {
	o := orm.NewOrm()
	var artiTypes []models.ArticleType
	_, err := o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil {
		logs.Info("查询错误")
	}
	this.Data["types"] = artiTypes
	this.TplName = "addType.html"
}
func (this *ArticleController) HandleAddType() {

	typename := this.GetString("typeName")
	if typename == "" {
		logs.Info("添加类型数据为空")
		return
	}
	o := orm.NewOrm()
	var artiType models.ArticleType
	artiType.TypeName = typename
	_, err := o.Insert(&artiType)
	if err != nil {
		logs.Info("插入失败")
		return
	}
	this.Redirect("/Article/AddArticleType", 302)
}

// 退出登录
func (this *ArticleController) Logout() {
	//删除登录状态
	this.DelSession("userName")
	//跳转登录界面
	this.Redirect("/", 302)
}
