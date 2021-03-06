package controllers

import (
	"github.com/astaxie/beego"
	"github.com/ulricqin/goutils/paginator"
	"github.com/vincentruan/beego-blog/g"
	"github.com/vincentruan/beego-blog/models/admin"
	"strconv"
)

type Checker interface {
	CheckLogin()
}

type BaseController struct {
	beego.Controller
	IsAdmin bool
}

func (this *BaseController) Prepare() {
	this.Data["BlogLogo"] = g.BlogLogo
	this.Data["BlogTitle"] = g.BlogTitle
	this.Data["BlogResume"] = g.BlogResume

	this.AssignIsAdmin()
	if app, ok := this.AppController.(Checker); ok {
		app.CheckLogin()
	}
}

func (this *BaseController) AssignIsAdmin() {
	if sess := this.GetSession("username"); sess != nil {
		beego.BeeLogger.Debug("Retrive user name [%s] from session", sess.(string))
		this.IsAdmin = true
	} else {
		//从cookie中获取
		bb_name := this.Ctx.GetCookie("bb_name")
		bb_password := this.Ctx.GetCookie("bb_password")
		if bb_name == "" || bb_password == "" {
			this.IsAdmin = false
		} else {
			this.IsAdmin = admin.CheckAdminUser(bb_name, bb_password, "host")
			if this.IsAdmin {
				this.SetSession("username", bb_name)
				beego.BeeLogger.Debug("Put user name [%s] into session", bb_name)
			}
		}
	}
	beego.BeeLogger.Debug("Is Admin [%t]", this.IsAdmin)
	this.Data["IsAdmin"] = this.IsAdmin
}

func (this *BaseController) SetPaginator(per int, nums int64) *paginator.Paginator {
	p := paginator.NewPaginator(this.Ctx.Request, per, nums)
	this.Data["paginator"] = p
	return p
}

func (this *BaseController) GetIntWithDefault(paramKey string, defaultVal int) int {
	valStr := this.GetString(paramKey)
	var val int
	if valStr == "" {
		val = defaultVal
	} else {
		var err error
		val, err = strconv.Atoi(valStr)
		if err != nil {
			val = defaultVal
		}
	}
	return val
}

func (this *BaseController) JsStorage(action, key string, values ...string) {
	value := action + ":::" + key
	if len(values) > 0 {
		value += ":::" + values[0]
	}
	this.Ctx.SetCookie("JsStorage", value, 1<<31-1, "/")
}

//返回分页参数
func (this *BaseController) GetPaginationParam() (limit, offset int) {
	var err error
	limit, err = this.GetInt("limit")
	if err != nil {
		beego.Info("Unable to format limit from parameters, use -1 as default!", err)
		limit = -1
	}
	offset, err = this.GetInt("offset")
	if err != nil {
		beego.Info("Unable to format offset from parameters, use 0 as default!", err)
		offset = 0
	}

	if limit == 0 {
		beego.Info("limit 0 is illegal, change to -1 as default!")
		limit = -1
	}

	return
}
