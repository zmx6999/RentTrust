package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

// 社区信息
type AreaInfo struct {
	AreaID       string `json:"area_id"`       // 社区编号
	AreaAddress  string `json:"area_address"`  // 房源所在区域
	BasicNetWork string `json:"basic_net_work"` // 区域基础网络编号
	CPoliceName  string `json:"c_police_name"`  // 社区民警姓名
	CPoliceNum   string `json:"c_police_num"`   // 社区民警工号
}

// 房源信息
type RentingHouseInfo struct {
	RentingID        string    `json:"renting_id"`         // 统一编码
	RentingAreaInfo  AreaInfo  `json:"renting_area_info"`  //区域信息
}

func checkArgsNum(args []string,n int) error {
	if len(args)<n {
		return errors.New(fmt.Sprintf("%d parameter(s) required",n))
	}
	return nil
}

type AreaChaincode struct {

}

func (this *AreaChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *AreaChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn,args:=stub.GetFunctionAndParameters()
	if fn=="addAreaInfo" {
		return this.addAreaInfo(stub,args)
	} else if fn=="getAreaInfo" {
		return this.getAreaInfo(stub,args)
	}
	return shim.Error("invalid request")
}

func (this *AreaChaincode) addAreaInfo(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,6)
	if err!=nil {
		return shim.Error(err.Error())
	}
	var rentingHouseInfo RentingHouseInfo
	rentingHouseInfo.RentingID=args[0]
	if rentingHouseInfo.RentingID=="" {
		return shim.Error("RentingID required")
	}
	rentingHouseInfo.RentingAreaInfo.AreaID=args[1]
	rentingHouseInfo.RentingAreaInfo.AreaAddress=args[2]
	rentingHouseInfo.RentingAreaInfo.BasicNetWork=args[3]
	rentingHouseInfo.RentingAreaInfo.CPoliceName=args[4]
	rentingHouseInfo.RentingAreaInfo.CPoliceNum=args[5]
	data,err:=json.Marshal(rentingHouseInfo)
	if err!=nil {
		return shim.Error(err.Error())
	}
	err=stub.PutState(rentingHouseInfo.RentingID,data)
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (this *AreaChaincode) getAreaInfo(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}
	var areaInfo AreaInfo
	var rentingHouseInfo RentingHouseInfo
	rentingHouseInfo.RentingID=args[0]
	if rentingHouseInfo.RentingID=="" {
		return shim.Error("RentingID required")
	}
	history,err:=stub.GetHistoryForKey(rentingHouseInfo.RentingID)
	if err!=nil {
		return shim.Error(err.Error())
	}
	defer history.Close()
	for history.HasNext() {
		item,err:=history.Next()
		if err!=nil {
			return shim.Error(err.Error())
		}
		err=json.Unmarshal(item.Value,&rentingHouseInfo)
		if err!=nil {
			return shim.Error(err.Error())
		}
		if rentingHouseInfo.RentingAreaInfo.AreaID!="" {
			areaInfo=rentingHouseInfo.RentingAreaInfo
		}
	}
	data,err:=json.Marshal(areaInfo)
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(data)
}

func main()  {
	shim.Start(new(AreaChaincode))
}
