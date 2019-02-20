package models

import (
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
)

type ChaincodeSpec struct {
	client apitxn.ChannelClient
	chaincodeId string
}

func Initialize(channelId string,user string,chaincodeId string,configFile string) (*ChaincodeSpec,error) {
	sdk,err:=fabapi.NewSDK(fabapi.Options{ConfigFile:configFile})
	if err!=nil {
		return nil,err
	}
	client,err:=sdk.NewChannelClient(channelId,user)
	if err!=nil {
		return nil,err
	}
	return &ChaincodeSpec{client,chaincodeId}, nil
}

func (this *ChaincodeSpec) ChaincodeUpdate(function string,args [][]byte) ([]byte,error) {
	request:=apitxn.ExecuteTxRequest{ChaincodeID:this.chaincodeId,Fcn:function,Args:args}
	id,err:=this.client.ExecuteTx(request)
	return []byte(id.ID),err
}

func (this *ChaincodeSpec) ChaincodeQuery(function string,args [][]byte) ([]byte,error)  {
	request:=apitxn.QueryRequest{ChaincodeID:this.chaincodeId,Fcn:function,Args:args}
	return this.client.Query(request)
}

func (this *ChaincodeSpec) Close()  {
	this.client.Close()
}
