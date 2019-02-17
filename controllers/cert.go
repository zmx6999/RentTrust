package controllers

import (
	"github.com/astaxie/beego"
	"190216/models"
	"path"
	"os"
	"encoding/csv"
	"strconv"
	"io"
	"errors"
	"strings"
	"time"
	"github.com/astaxie/beego/toolbox"
	"fmt"
)

type CertController struct {
	BaseController
}

func (this *CertController) Check()  {
	id:=this.GetString("id")
	name:=this.GetString("name")
	if id=="" || name=="" {
		this.handleResponse(400,"invalid parameter")
		return
	}
	channelId:=beego.AppConfig.String("channel_cert_id")
	user:=beego.AppConfig.String("user_cert")
	chaincodeId:=beego.AppConfig.String("chaincode_cert_id")
	configFile:=beego.AppConfig.String("CORE_CERT_CONFIG_FILE")
	ccs,err:=models.Initialize(channelId,user,chaincodeId,configFile)
	if err!=nil {
		this.handleResponse(500,err.Error())
		return
	}
	defer ccs.Close()
	data,err:=ccs.ChaincodeQuery("check",[][]byte{[]byte(id),[]byte(name)})
	if err!=nil {
		this.handleResponse(500,err.Error())
		return
	}
	this.handleResponse(200,string(data))
}

func (this *CertController) Add()  {
	file,head,err:=this.GetFile("cert")
	if err!=nil {
		this.handleResponse(500,err.Error())
		return
	}
	defer file.Close()
	err=this.SaveToFile("cert",path.Join("static/upload",head.Filename))
	if err!=nil {
		this.handleResponse(500,err.Error())
		return
	}
	t:=time.Now().Add(time.Second*5)
	tk:=toolbox.NewTask("cert_task", fmt.Sprintf("%d %d %d * * *", t.Second(), t.Minute(), t.Hour()), func() error {
		defer toolbox.StopTask()
		return certTask(path.Join("static/upload",head.Filename))
	})
	toolbox.AddTask("cert_task",tk)
	toolbox.StartTask()
	this.handleResponse(200,"ok")
}

func certTask(fileName string) error {
	channelId:=beego.AppConfig.String("channel_cert_id")
	user:=beego.AppConfig.String("user_cert")
	chaincodeId:=beego.AppConfig.String("chaincode_cert_id")
	configFile:=beego.AppConfig.String("CORE_CERT_CONFIG_FILE")
	ccs,err:=models.Initialize(channelId,user,chaincodeId,configFile)
	if err!=nil {
		return err
	}
	defer ccs.Close()
	file,err:=os.Open(fileName)
	if err!=nil {
		return err
	}
	defer file.Close()
	reader:=csv.NewReader(file)
	i:=0
	var errLines []string
	for  {
		i+=1
		line,err:=reader.Read()
		if err==io.EOF {
			break
		}
		if err!=nil {
			errLines=append(errLines,strconv.Itoa(i))
			continue
		}
		if len(line)<2 {
			errLines=append(errLines,strconv.Itoa(i))
			continue
		}
		var args [][]byte
		for _,value:=range line{
			args=append(args,[]byte(value))
		}
		_,err=ccs.ChaincodeUpdate("add",args)
		if err!=nil {
			errLines=append(errLines,strconv.Itoa(i))
			continue
		}
	}
	if len(errLines)>0 {
		beego.Error("Line "+strings.Join(errLines,",")+" invalid format")
		return errors.New("Line "+strings.Join(errLines,",")+" invalid format")
	} else {
		beego.Info("success")
		return nil
	}
}
