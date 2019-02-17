package routers

import (
	"190216/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/auth/check",&controllers.AuthController{},"get:Check")
	beego.Router("/auth/add",&controllers.AuthController{},"post:Add")

	beego.Router("/cert/check",&controllers.CertController{},"get:Check")
	beego.Router("/cert/add",&controllers.CertController{},"post:Add")

	beego.Router("/credit/check",&controllers.CreditController{},"get:Check")
	beego.Router("/credit/add",&controllers.CreditController{},"post:Add")
}
