package emailClient

import (
	"log"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

type EmailClient interface {
	Send(receiveUserList []EmailUser, message EmailMessage) error
	SendSpecial(receiveUserList []EmailUser) error
}

type emailClient struct {
	cfg EmailConf
}

type EmailConf struct {
	SendEmail string `json:"send_email" mapstructure:"send_email"`
	Host      string `json:"host" mapstructure:"host"`
	Port      int    `json:"port" mapstructure:"port"`
	Account   string `json:"account" mapstructure:"account"`
	Password  string `json:"password" mapstructure:"password"`
}

func NewEmailClient(cfg EmailConf) (EmailClient, func(), error) {
	return &emailClient{
		cfg: cfg,
	}, nil, nil
}

// EmailUser 用户邮箱
type EmailUser struct {
	Name         string
	EmailAddress string
	EmailMessage
}

type EmailMessage struct {
	Subject            string
	ContentType        string
	Body               string
	AttachFilePathList []string
}

// 设置邮箱发送方
func (cli *emailClient) setSender(m *gomail.Message) error {
	if cli.cfg.SendEmail == "" {
		return errors.New("发送地址")
	}
	m.SetHeader("From", cli.cfg.SendEmail)
	return nil
}

// 设置邮件标题
func (cli *emailClient) setSubject(m *gomail.Message, subject string) {
	m.SetHeader("Subject", subject)
}

// 设置邮件内容
func (cli *emailClient) setBody(m *gomail.Message, contentType string, body string) {
	m.SetBody(contentType, body)
}

// 设置邮件附件
func (cli *emailClient) setAttach(m *gomail.Message, attachFilePathList []string) {
	if len(attachFilePathList) == 0 {
		return
	}
	for _, filePath := range attachFilePathList {
		m.Attach(filePath)
	}
}

// Send 发送邮箱，邮件内容相同
func (cli *emailClient) Send(receiveUserList []EmailUser, message EmailMessage) error {
	m := gomail.NewMessage()
	if err := cli.setSender(m); err != nil {
		return err
	}

	if len(receiveUserList) == 0 {
		return errors.New("接收用户列表为空")
	}
	// 接收方
	addressList := make([]string, 0, len(receiveUserList))
	for _, user := range receiveUserList {
		addressList = append(addressList, user.EmailAddress)
	}
	m.SetHeader("To", addressList...)

	// 邮件标题
	cli.setSubject(m, message.Subject)
	// 邮件内容
	cli.setBody(m, message.ContentType, message.Body)
	// 邮件附件
	cli.setAttach(m, message.AttachFilePathList)

	// 连接并发送邮箱
	dialer := gomail.NewDialer(cli.cfg.Host, cli.cfg.Port, cli.cfg.Account, cli.cfg.Password)
	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// SendSpecial 给对应用户发送专属邮件，每个邮件内容不同
func (cli *emailClient) SendSpecial(receiveUserList []EmailUser) error {
	m := gomail.NewMessage()
	if err := cli.setSender(m); err != nil {
		return err
	}

	dialer := gomail.NewDialer(cli.cfg.Host, cli.cfg.Port, cli.cfg.Account, cli.cfg.Password)
	dial, err := dialer.Dial()
	if err != nil {
		return err
	}
	for _, user := range receiveUserList {
		// 设置发送者
		if err := cli.setSender(m); err != nil {
			return err
		}
		// 设置接收者，指定邮件接收者名称，若使用SetHeader则默认使用邮箱名
		m.SetAddressHeader("To", user.EmailAddress, user.Name)
		// 邮件标题
		cli.setSubject(m, user.Subject)
		// 邮件内容
		cli.setBody(m, user.ContentType, user.Body)
		// 邮件附件
		cli.setAttach(m, user.AttachFilePathList)

		if err := gomail.Send(dial, m); err != nil {
			log.Printf("Could not send email to %q: %v", user.EmailAddress, err)
		}
		m.Reset()
	}
	return nil
}
