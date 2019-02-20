package controllers

import "github.com/astaxie/beego"

type ResponseJSON struct {
	Code int
	Msg string
	Data interface{}
}

type BaseController struct {
	beego.Controller
}

func (this *BaseController) handleResponse(code int,msg string,data interface{})  {
	this.Data["json"]=&ResponseJSON{code,msg,data}
	this.ServeJSON()
}
