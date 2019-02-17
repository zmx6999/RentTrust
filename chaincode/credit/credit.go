package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"errors"
	"fmt"
)

func main()  {
	shim.Start(new(CreditChaincode))
}

type CreditChaincode struct {

}

func (this *CreditChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *CreditChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn,args:=stub.GetFunctionAndParameters()
	if fn=="check" {
		return this.check(stub,args)
	} else if fn=="add" {
		return this.add(stub,args)
	}
	return shim.Error("Method doesn't exist")
}

func (this *CreditChaincode) check(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,1)
	if err!=nil {
		return shim.Error(err.Error())
	}
	id:=args[0]
	data,err:=stub.GetState(id)
	if err!=nil {
		return shim.Error(err.Error())
	}
	if data==nil {
		return shim.Success([]byte("false"))
	}
	return shim.Success(data)
}

func (this *CreditChaincode) add(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	id:=args[0]
	credit:=args[1]
	err=stub.PutState(id,[]byte(credit))
	if err!=nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func checkArgsNum(args []string,n int) error {
	if len(args)!=n {
		return errors.New(fmt.Sprintf("%d parameter(s) required",n))
	}
	return nil
}
