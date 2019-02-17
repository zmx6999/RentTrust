package controllers

import (
		"github.com/astaxie/beego"
	"190216/models"
		"path"
	"time"
	"github.com/astaxie/beego/toolbox"
	"fmt"
	"os"
	"encoding/csv"
	"io"
	"strconv"
	"errors"
	"strings"
)

type AuthController struct {
	BaseController
}

func (this *AuthController) Check()  {
	id:=this.GetString("id")
	name:=this.GetString("name")
	if id=="" || name=="" {
		this.handleResponse(400,"invalid parameter")
		return
	}
	channelId:=beego.AppConfig.String("channel_auth_id")
	user:=beego.AppConfig.String("user_auth")
	chaincodeId:=beego.AppConfig.String("chaincode_auth_id")
	configFile:=beego.AppConfig.String("CORE_AUTH_CONFIG_FILE")
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

func (this *AuthController) Add()  {
	file,head,err:=this.GetFile("auth")
	if err!=nil {
		this.handleResponse(500,err.Error())
		return
	}
	defer file.Close()
	err=this.SaveToFile("auth",path.Join("static/upload",head.Filename))
	if err!=nil {
		this.handleResponse(500,err.Error())
		return
	}
	t:=time.Now().Add(time.Second*5)
	tk:=toolbox.NewTask("auth_task", fmt.Sprintf("%d %d %d * * *", t.Second(), t.Minute(), t.Hour()), func() error {
		defer toolbox.StopTask()
		return authTask(path.Join("static/upload",head.Filename))
	})
	toolbox.AddTask("auth_task",tk)
	toolbox.StartTask()
	this.handleResponse(200,"ok")
}

func authTask(fileName string) error {
	channelId:=beego.AppConfig.String("channel_auth_id")
	user:=beego.AppConfig.String("user_auth")
	chaincodeId:=beego.AppConfig.String("chaincode_auth_id")
	configFile:=beego.AppConfig.String("CORE_AUTH_CONFIG_FILE")
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
		iStr:=strconv.Itoa(i)
		line,err:=reader.Read()
		if err==io.EOF {
			break
		}
		if err!=nil {
			errLines=append(errLines,iStr)
			continue
		}
		if len(line)<2 {
			errLines=append(errLines,iStr)
			continue
		}
		var args [][]byte
		for _,value:=range line{
			args=append(args,[]byte(value))
		}
		_,err=ccs.ChaincodeUpdate("add",args)
		if err!=nil {
			errLines=append(errLines,iStr)
			continue
		}
	}
	if len(errLines)>0 {
		beego.Error("Lines "+strings.Join(errLines,",")+" invalid format")
		return errors.New("Lines "+strings.Join(errLines,",")+" invalid format")
	} else {
		beego.Info("success")
		return nil
	}
}
