package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"190221/models"
)

type AreaController struct {
	BaseController
}

func (this *AreaController) AddAreaInfo()  {
	var responseBody map[string]interface{}
	json.Unmarshal(this.Ctx.Input.RequestBody,&responseBody)
	rentingID:=responseBody["renting_id"].(string)
	if rentingID=="" {
		this.handleResponse(400,"invalid request",nil)
		return
	}
	areaID:=responseBody["area_id"].(string)
	areaAddress:=responseBody["area_address"].(string)
	basicNetWork:=responseBody["basic_net_work"].(string)
	cPoliceName:=responseBody["c_police_name"].(string)
	cPoliceNum:=responseBody["c_police_num"].(string)

	channelId:=beego.AppConfig.String("area_channel_id")
	user:=beego.AppConfig.String("area_user")
	chaincodeId:=beego.AppConfig.String("area_chaincode_id")
	configFile:=beego.AppConfig.String("area_config_file")
	ccs,err:=models.Initialize(channelId,user,chaincodeId,configFile)
	if err!=nil {
		this.handleResponse(500,err.Error(),nil)
		return
	}
	defer ccs.Close()
	txId,err:=ccs.ChaincodeUpdate("addAreaInfo",[][]byte{[]byte(rentingID),[]byte(areaID),[]byte(areaAddress),[]byte(basicNetWork),[]byte(cPoliceName),[]byte(cPoliceNum)})
	if err!=nil {
		this.handleResponse(500,err.Error(),nil)
		return
	}
	this.handleResponse(200,"OK",txId)
}

func (this *AreaController) GetAreaInfo()  {
	rentingID:=this.GetString("renting_id")
	if rentingID=="" {
		this.handleResponse(400,"invalid request",nil)
		return
	}

	channelId:=beego.AppConfig.String("area_channel_id")
	user:=beego.AppConfig.String("area_user")
	chaincodeId:=beego.AppConfig.String("area_chaincode_id")
	configFile:=beego.AppConfig.String("area_config_file")
	ccs,err:=models.Initialize(channelId,user,chaincodeId,configFile)
	if err!=nil {
		this.handleResponse(500,err.Error(),nil)
		return
	}
	defer ccs.Close()
	_data,err:=ccs.ChaincodeQuery("getAreaInfo",[][]byte{[]byte(rentingID)})
	if err!=nil {
		this.handleResponse(500,err.Error(),nil)
		return
	}
	var data map[string]interface{}
	json.Unmarshal(_data,&data)
	this.handleResponse(200,"OK",data)
}
