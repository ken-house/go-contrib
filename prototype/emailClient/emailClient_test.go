package emailClient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailClient_Send(t *testing.T) {
	emailConf := EmailConf{
		SendEmail: "*********@qq.com",
		Host:      "smtp.qq.com",
		Port:      465,
		Account:   "*********@qq.com",
		Password:  "*********", // 授权码
	}
	emailClient, _, err := NewEmailClient(emailConf)
	if err != nil {
		assert.Failf(t, err.Error(), "创建失败")
	}
	userList := []EmailUser{
		{
			Name:         "Ken",
			EmailAddress: "xudengtang@zonst.cn",
		},
	}

	message := EmailMessage{
		Subject:     "测试邮件",
		ContentType: "text/html",
		Body:        "<h1 style='color:#F00'>hello world</h1>",
		AttachFilePathList: []string{
			"/Users/zonst/Downloads/WechatIMG51.jpeg",
			"/Users/zonst/Downloads/download.zip",
		},
	}
	if err := emailClient.Send(userList, message); err != nil {
		assert.Failf(t, err.Error(), "发送失败")
	}
	assert.True(t, true)
}

func TestEmailClient_SendSpecial(t *testing.T) {
	emailConf := EmailConf{
		SendEmail: "*********@qq.com",
		Host:      "smtp.qq.com",
		Port:      465,
		Account:   "*********@qq.com",
		Password:  "*********",
	}
	emailClient, _, err := NewEmailClient(emailConf)
	if err != nil {
		assert.Failf(t, err.Error(), "创建失败")
	}
	userList := []EmailUser{
		{
			Name:         "Ken",
			EmailAddress: "xudengtang@zonst.cn",
			EmailMessage: EmailMessage{
				Subject:     "测试邮件ken",
				ContentType: "text/html",
				Body:        "<h1 style='color:#F00'>hello world</h1>",
				AttachFilePathList: []string{
					"/Users/zonst/Downloads/WechatIMG51.jpeg",
					"/Users/zonst/Downloads/download.zip",
				},
			},
		},
		{
			Name:         "Ken",
			EmailAddress: "xudengtang@zonst.cn",
			EmailMessage: EmailMessage{
				Subject:     "测试邮件ken2",
				ContentType: "text/html",
				Body:        "<h1 style='color:#F00'>hello gomail</h1>",
				AttachFilePathList: []string{
					"/Users/zonst/Downloads/WechatIMG51.jpeg",
					"/Users/zonst/Downloads/download.zip",
				},
			},
		},
	}

	if err := emailClient.SendSpecial(userList); err != nil {
		assert.Failf(t, err.Error(), "发送失败")
	}
	assert.True(t, true)
}
