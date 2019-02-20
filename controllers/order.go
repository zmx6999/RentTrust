package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"190221/models"
)

type OrderController struct {
	BaseController
}

func (this *OrderController) AddOrderInfo()  {
	var responseBody map[string]interface{}
	json.Unmarshal(this.Ctx.Input.RequestBody,&responseBody)
	rentingID:=responseBody["renting_id"].(string)
	if rentingID=="" {
		this.handleResponse(400,"invalid request",nil)
		return
	}
	docHash:=responseBody["doc_hash"].(string)
	orderID:=responseBody["order_id"].(string)
	renterID:=responseBody["renter_id"].(string)
	rentMoney:=responseBody["rent_money"].(string)
	beginDate:=responseBody["begin_date"].(string)
	endDate:=responseBody["end_date"].(string)
	note:=responseBody["note"].(string)

	channelId:=beego.AppConfig.String("order_channel_id")
	user:=beego.AppConfig.String("order_user")
	chaincodeId:=beego.AppConfig.String("order_chaincode_id")
	configFile:=beego.AppConfig.String("order_config_file")
	ccs,err:=models.Initialize(channelId,user,chaincodeId,configFile)
	if err!=nil {
		this.handleResponse(500,err.Error(),nil)
		return
	}
	defer ccs.Close()
	txId,err:=ccs.ChaincodeUpdate("addOrderInfo",[][]byte{[]byte(rentingID),[]byte(docHash),[]byte(orderID),[]byte(renterID),[]byte(rentMoney),[]byte(beginDate),[]byte(endDate),[]byte(note)})
	if err!=nil {
		this.handleResponse(500,err.Error(),nil)
		return
	}
	this.handleResponse(200,"OK",txId)
}

func (this *OrderController) GetOrderInfo()  {
	rentingID:=this.GetString("renting_id")
	if rentingID=="" {
		this.handleResponse(400,"invalid request",nil)
		return
	}

	channelId:=beego.AppConfig.String("order_channel_id")
	user:=beego.AppConfig.String("order_user")
	chaincodeId:=beego.AppConfig.String("order_chaincode_id")
	configFile:=beego.AppConfig.String("order_config_file")
	ccs,err:=models.Initialize(channelId,user,chaincodeId,configFile)
	if err!=nil {
		this.handleResponse(500,err.Error(),nil)
		return
	}
	defer ccs.Close()
	_data,err:=ccs.ChaincodeQuery("getOrderInfo",[][]byte{[]byte(rentingID)})
	if err!=nil {
		this.handleResponse(500,err.Error(),nil)
		return
	}
	var data []map[string]interface{}
	json.Unmarshal(_data,&data)
	this.handleResponse(200,"OK",data)
}
