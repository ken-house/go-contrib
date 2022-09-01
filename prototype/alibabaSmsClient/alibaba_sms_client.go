package alibabaSmsClient

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/pkg/errors"
)

type AlibabaSmsClient interface {
	SendCode(params SendSmsParams) (err error)
}

type alibabaSmsClient struct {
	Client *dysmsapi20170525.Client
}

// ClientConfig 连接客户端配置
type ClientConfig struct {
	Endpoint        string `json:"endpoint" mapstructure:"endpoint"`
	AccessKeyId     string `json:"access_key_id" mapstructure:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret" mapstructure:"access_key_secret"`
}

// SendSmsParams 发送短信参数结构体
type SendSmsParams struct {
	Phone         string `json:"phone" mapstructure:"phone"`
	SignName      string `json:"sign_name" mapstructure:"sign_name"`
	TemplateCode  string `json:"template_code" mapstructure:"template_code"`
	TemplateParam string `json:"template_param" mapstructure:"template_param"`
}

// CreateClient 使用AK&SK初始化账号Client
func CreateClient(cfg ClientConfig) (AlibabaSmsClient, error) {
	config := &openapi.Config{
		// 访问的域名
		Endpoint: &cfg.Endpoint,
		// 您的 AccessKey ID
		AccessKeyId: &cfg.AccessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: &cfg.AccessKeySecret,
	}
	client, err := dysmsapi20170525.NewClient(config)
	return &alibabaSmsClient{
		Client: client,
	}, err
}

// SendCode 发送短信验证码
func (cli *alibabaSmsClient) SendCode(params SendSmsParams) (err error) {
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(params.Phone),
		SignName:      tea.String(params.SignName),
		TemplateCode:  tea.String(params.TemplateCode),
		TemplateParam: tea.String(params.TemplateParam),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (err error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, err = cli.Client.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			return err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		errMsg := util.AssertAsString(error.Message)
		return errors.New(*errMsg)
	}
	return err
}
