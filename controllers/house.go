package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"190221/models"
)

type HouseController struct {
	BaseController
}

func (this *HouseController) AddHouseInfo()  {
	var responseBody map[string]interface{}
	json.Unmarshal(this.Ctx.Input.RequestBody,&responseBody)
	rentingID:=responseBody["renting_id"].(string)
	if rentingID=="" {
		this.handleResponse(400,"invalid request",nil)
		return
	}
	houseID:=responseBody["house_id"].(string)
	houseOwner:=responseBody["house_owner"].(string)
	regDate:=responseBody["reg_date"].(string)
	houseArea:=responseBody["house_area"].(string)
	houseUsed:=responseBody["house_used"].(string)
	isMortgage:=responseBody["is_mortgage"].(string)

	channelId:=beego.AppConfig.String("house_channel_id")
	user:=beego.AppConfig.String("house_user")
	chaincodeId:=beego.AppConfig.String("house_chaincode_id")
	configFile:=beego.AppConfig.String("house_config_file")
	ccs,err:=models.Initialize(channelId,user,chaincodeId,configFile)
	if err!=nil {
		this.handleResponse(500,err.Error(),nil)
		return
	}
	defer ccs.Close()
	txId,err:=ccs.ChaincodeUpdate("addHouseInfo",[][]byte{[]byte(rentingID),[]byte(houseID),[]byte(houseOwner),[]byte(regDate),[]byte(houseArea),[]byte(houseUsed),[]byte(isMortgage)})
	if err!=nil {
		this.handleResponse(500,err.Error(),nil)
		return
	}
	this.handleResponse(200,"OK",txId)
}

func (this *HouseController) GetHouseInfo()  {
	rentingID:=this.GetString("renting_id")
	if rentingID=="" {
		this.handleResponse(400,"invalid request",nil)
		return
	}

	channelId:=beego.AppConfig.String("house_channel_id")
	user:=beego.AppConfig.String("house_user")
	chaincodeId:=beego.AppConfig.String("house_chaincode_id")
	configFile:=beego.AppConfig.String("house_config_file")
	ccs,err:=models.Initialize(channelId,user,chaincodeId,configFile)
	if err!=nil {
		this.handleResponse(500,err.Error(),nil)
		return
	}
	defer ccs.Close()
	_data,err:=ccs.ChaincodeQuery("getHouseInfo",[][]byte{[]byte(rentingID)})
	if err!=nil {
		this.handleResponse(500,err.Error(),nil)
		return
	}
	var data map[string]interface{}
	json.Unmarshal(_data,&data)
	this.handleResponse(200,"OK",data)
}
