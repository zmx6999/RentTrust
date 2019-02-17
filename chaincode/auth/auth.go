package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"errors"
	"fmt"
)

func main()  {
	shim.Start(new(AuthChaincode))
}

type AuthChaincode struct {

}

func (this *AuthChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (this *AuthChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn,args:=stub.GetFunctionAndParameters()
	if fn=="check" {
		return this.check(stub,args)
	} else if fn=="add" {
		return this.add(stub,args)
	}
	return shim.Error("Method doesn't exist")
}

func (this *AuthChaincode) check(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	id:=args[0]
	name:=args[1]
	data,err:=stub.GetState(id)
	if err!=nil {
		return shim.Error(err.Error())
	}
	if data==nil {
		return shim.Success([]byte("false"))
	}
	if string(data)==name {
		return shim.Success([]byte("true"))
	} else {
		return shim.Success([]byte("false"))
	}
}

func (this *AuthChaincode) add(stub shim.ChaincodeStubInterface,args []string) peer.Response {
	err:=checkArgsNum(args,2)
	if err!=nil {
		return shim.Error(err.Error())
	}
	id:=args[0]
	name:=args[1]
	err=stub.PutState(id,[]byte(name))
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
