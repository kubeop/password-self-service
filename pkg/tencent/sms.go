package tencent

import (
	"fmt"
	"password-self-service/pkg/config"
	"password-self-service/pkg/logging"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func SendTencentSMS(code, mobile, templateCode string) error {
	credential := common.NewCredential(
		config.Setting.Channel.TencentSms.AccessKeyId,
		config.Setting.Channel.TencentSms.AccessKeySecret,
	)

	// 实例化一个客户端配置对象，可以指定超时时间等配置
	cpf := profile.NewClientProfile()

	// 使用POST方法
	cpf.HttpProfile.ReqMethod = "POST"

	// 指定接入地域域名，默认就近地域接入域名
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	// SDK默认用TC3-HMAC-SHA256进行签名，非必要请不要修改这个字段
	cpf.SignMethod = "HmacSHA1"

	// 实例化要请求产品(以sms为例)的client对象，仅支持ap-beijing、ap-guangzhou、ap-nanjing
	client, _ := sms.NewClient(credential, config.Setting.Channel.TencentSms.Region, cpf)

	// 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	request := sms.NewSendSmsRequest()

	// 短信应用ID
	request.SmsSdkAppId = common.StringPtr(config.Setting.Channel.TencentSms.AppId)

	// 短信签名内容
	request.SignName = common.StringPtr(config.Setting.Channel.TencentSms.SignName)

	// 模板 ID
	request.TemplateId = common.StringPtr(templateCode)

	// 模板参数
	request.TemplateParamSet = common.StringPtrs([]string{code})

	// 发送的手机号
	request.PhoneNumberSet = common.StringPtrs([]string{mobile})

	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.SendSms(request)
	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logging.Logger().Sugar().Errorf("An API error has returned: %s", err)
		return err
	}

	if err != nil {
		return err
	}

	if gconv.String(response.Response.SendStatusSet[0].Code) != "Ok" {
		return fmt.Errorf("%v", gconv.String(response.Response.SendStatusSet[0].Message))
	}

	logging.Logger().Sugar().Infof("发送消息请求状态: %v, 消息发送状态: %v, 消息发送回执ID: %v", gconv.String(response.Response.SendStatusSet[0].Code), gconv.String(response.Response.SendStatusSet[0].Message), gconv.String(response.Response.SendStatusSet[0].SerialNo))
	return nil
}
