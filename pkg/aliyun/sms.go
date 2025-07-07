package aliyun

import (
	"errors"
	"password-self-service/pkg/config"
	"password-self-service/pkg/logging"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	teautils "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
)

func SendAliyunSMS(code, mobile, templateCode string) error {
	credConfig := new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(config.Setting.Channel.AliyunSms.AccessKeyId).
		SetAccessKeySecret(config.Setting.Channel.AliyunSms.AccessKeySecret)

	credential, err := credentials.NewCredential(credConfig)
	if err != nil {
		logging.Logger().Sugar().Errorf("初始化凭据错误：%v", err)
		return errors.New("初始化凭据错误")
	}

	clientConfig := &openapi.Config{
		Credential: credential,
		Endpoint:   tea.String("dysmsapi.aliyuncs.com"),
	}

	client := &dysmsapi20170525.Client{}
	client, err = dysmsapi20170525.NewClient(clientConfig)
	if err != nil {
		logging.Logger().Sugar().Errorf("初始化客户端错误：%v", err)
		return errors.New("初始化客户端错误")
	}

	params := &openapi.Params{
		// 接口名称
		Action: tea.String("SendSms"),
		// 接口版本
		Version: tea.String("2017-05-25"),
		// 接口协议
		Protocol: tea.String("HTTPS"),
		// 接口 HTTP 方法
		Method:   tea.String("POST"),
		AuthType: tea.String("AK"),
		Style:    tea.String("RPC"),
		// 接口 PATH
		Pathname: tea.String("/"),
		// 接口请求体内容格式
		ReqBodyType: tea.String("json"),
		// 接口响应体内容格式
		BodyType: tea.String("json"),
	}

	mobile = strings.ReplaceAll(mobile, " ", "")

	queries := map[string]interface{}{}
	queries["PhoneNumbers"] = tea.String(mobile)
	queries["SignName"] = tea.String(config.Setting.Channel.AliyunSms.SignName)
	queries["TemplateCode"] = tea.String(templateCode)
	queries["TemplateParam"] = tea.String("{\"code\":\"" + code + "\"}")

	runtime := &teautils.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	result, err := client.CallApi(params, request, runtime)
	if err != nil {
		logging.Logger().Sugar().Errorf("短信发送失败：%v", err)
		return errors.New("短信发送失败")
	}

	info := result["body"].(map[string]interface{})
	logging.Logger().Sugar().Infof("发送消息请求状态: %v, 消息发送状态: %v, 消息发送回执ID: %v", info["Code"], info["Message"], info["BizId"])

	return nil
}
