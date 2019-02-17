package controllers

import "github.com/astaxie/beego"

type BaseController struct {
	beego.Controller
}

type ResponseJSON struct {
	Code int
	Msg interface{}
}

func (this *BaseController) handleResponse(code int,msg interface{})  {
	this.Data["json"]=&ResponseJSON{Code:code,Msg:msg}
	this.ServeJSON()
}

