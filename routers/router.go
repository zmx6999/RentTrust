package routers

import (
	"190221/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/house/add",&controllers.HouseController{},"post:AddHouseInfo")
	beego.Router("/house/info",&controllers.HouseController{},"get:GetHouseInfo")

	beego.Router("/area/add",&controllers.AreaController{},"post:AddAreaInfo")
	beego.Router("/area/info",&controllers.AreaController{},"get:GetAreaInfo")

	beego.Router("/order/add",&controllers.OrderController{},"post:AddOrderInfo")
	beego.Router("/order/list",&controllers.OrderController{},"get:GetOrderInfo")
}
