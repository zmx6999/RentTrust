package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

// 订单信息
type OrderInfo struct {
	DocHash   string `json:"doc_hash"`    // 电子合同Hash
	OrderID   string `json:"order_id"`    // 订单编号
	RenterID  string `json:"renter_id"` // 承租人信息
	RentMoney string `json:"rent_money"`  // 租金
	BeginDate string `json:"begin_date"`  // 开始日期
	EndDate   string `json:"end_date"`    // 结束日期
	Note      string `json:"note"`       // 备注
}

// 房源信息
type RentingHouseInfo struct {
	RentingID        string    `json:"renting_id"`         // 统一编码
	RentingOrderInfo OrderInfo `json:"renting_order_info"` //订单信息
}

func checkArgsNum(args []string,n int) error {
	if len(args)<n {
		return errors.New(fmt.Sprintf("%d parameter(s) required",n))
	}
	return nil
}

type OrderChaincode struct {

}

func (this *OrderChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *OrderChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn,args:=stub.GetFunctionAndParameters()
	if fn=="addOrderInfo" {
		return this.addOrderInfo(stub,args)
	} else if fn=="getOrderInfo" {
		return this.getOrderInfo(stub,args)
	}
	return shim.Error("invalid request")
}

func (this *OrderChaincode) addOrderInfo(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,8)
	if err!=nil {
		return shim.Error(err.Error())
	}
	var rentingHouseInfo RentingHouseInfo
	rentingHouseInfo.RentingID=args[0]
	if rentingHouseInfo.RentingID=="" {
		return shim.Error("RentingID required")
	}
	rentingHouseInfo.RentingOrderInfo.DocHash=args[1]
	rentingHouseInfo.RentingOrderInfo.OrderID=args[2]
	rentingHouseInfo.RentingOrderInfo.RenterID=args[3]
	rentingHouseInfo.RentingOrderInfo.RentMoney=args[4]
	rentingHouseInfo.RentingOrderInfo.BeginDate=args[5]
	rentingHouseInfo.RentingOrderInfo.EndDate=args[6]
	rentingHouseInfo.RentingOrderInfo.Note=args[7]
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

func (this *OrderChaincode) getOrderInfo(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}
	var orderInfoList []OrderInfo
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
		if rentingHouseInfo.RentingOrderInfo.DocHash!="" {
			orderInfoList=append(orderInfoList,rentingHouseInfo.RentingOrderInfo)
		}
	}
	data,err:=json.Marshal(orderInfoList)
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(data)
}

func main()  {
	shim.Start(new(OrderChaincode))
}
