package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

// 房屋信息
type HouseInfo struct {
	HouseID    string `json:"house_id"`    // 房产证编号
	HouseOwner string `json:"house_owner"` // 房主
	RegDate    string `json:"reg_date"`    // 登记日期
	HouseArea  string `json:"house_area"`  // 住房面积
	HouseUsed  string `json:"house_used"`  // 房屋设计用途
	IsMortgage string `json:"is_mortgage"` // 是否抵押
}

// 房源信息
type RentingHouseInfo struct {
	RentingID        string    `json:"renting_id"`         // 统一编码
	RentingHouseInfo HouseInfo `json:"renting_house_info"` // 房屋信息
}

func checkArgsNum(args []string,n int) error {
	if len(args)<n {
		return errors.New(fmt.Sprintf("%d parameter(s) required",n))
	}
	return nil
}

type HouseChaincode struct {

}

func (this *HouseChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *HouseChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn,args:=stub.GetFunctionAndParameters()
	if fn=="addHouseInfo" {
		return this.addHouseInfo(stub,args)
	} else if fn=="getHouseInfo" {
		return this.getHouseInfo(stub,args)
	}
	return shim.Error("invalid request")
}

func (this *HouseChaincode) addHouseInfo(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,7)
	if err!=nil {
		return shim.Error(err.Error())
	}
	var rentingHouseInfo RentingHouseInfo
	rentingHouseInfo.RentingID=args[0]
	if rentingHouseInfo.RentingID=="" {
		return shim.Error("RentingID required")
	}
	rentingHouseInfo.RentingHouseInfo.HouseID=args[1]
	rentingHouseInfo.RentingHouseInfo.HouseOwner=args[2]
	rentingHouseInfo.RentingHouseInfo.RegDate=args[3]
	rentingHouseInfo.RentingHouseInfo.HouseArea=args[4]
	rentingHouseInfo.RentingHouseInfo.HouseUsed=args[5]
	rentingHouseInfo.RentingHouseInfo.IsMortgage=args[6]
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

func (this *HouseChaincode) getHouseInfo(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}
	var houseInfo HouseInfo
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
		if rentingHouseInfo.RentingHouseInfo.HouseOwner!="" {
			houseInfo=rentingHouseInfo.RentingHouseInfo
		}
	}
	data,err:=json.Marshal(houseInfo)
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(data)
}

func main()  {
	shim.Start(new(HouseChaincode))
}
