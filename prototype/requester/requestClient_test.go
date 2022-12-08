package requester

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestRequestClient(t *testing.T) {
	responseData := struct {
		Hello string `json:"hello"`
	}{}
	pool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("/Users/zonst/Documents/work/go-contrib/prototype/requester/cert/ca.crt")
	if err != nil {
		log.Fatalln(err)
	}
	pool.AppendCertsFromPEM(ca)

	cliCrt, err := tls.LoadX509KeyPair("/Users/zonst/Documents/work/go-contrib/prototype/requester/cert/client.crt", "/Users/zonst/Documents/work/go-contrib/prototype/requester/cert/client.key")
	if err != nil {
		log.Fatalln("LoadX509KeyPair error:", err.Error())
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 连接超时时间
			KeepAlive: 30 * time.Second, // 长连接保持时间
		}).DialContext,
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
		TLSHandshakeTimeout:   0, // 等待TLS握手。零表示没有超时。
		IdleConnTimeout:       0, // 连接空闲超时时间
		ResponseHeaderTimeout: 0, // 响应头信息超时时间
		ExpectContinueTimeout: 0, // 完全恢复后等待服务器第一个响应标头的时间
	}
	client := &http.Client{
		Transport:     transport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	httpClient := NewRequestClient("https://localhost:10443", nil, client)
	// 开启调试模式
	context := context.WithValue(context.Background(), "debug", true)
	response, err := httpClient.Get(context, "/test", &responseData, nil)
	if err != nil {
		zap.L().Error("请求失败", zap.Error(err))
		assert.Fail(t, err.Error())
		return
	}
	fmt.Println(response)
	fmt.Printf("responseData：%+v\n", responseData)
	assert.Equal(t, err, nil)
}
