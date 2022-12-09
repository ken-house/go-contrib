package tools

import (
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// IsContain 检查字符串是否在slice
func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// GetClientIp 获取真实客户端IP
func GetClientIp(r *http.Request) string {
	// 尝试从 X-Forwarded-For 中获取
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	// 如果 X-Forwarded-For 有值，取第一个非unknown的ip
	clientIp := ""
	ipArr := strings.Split(xForwardedFor, ",")
	for _, ip := range ipArr {
		ip = strings.TrimSpace(ip)
		if ip != "" && strings.ToLower(ip) != "unknown" {
			clientIp = ip
			break
		}
	}
	if clientIp == "" {
		// 尝试从 X-Real-Ip 中获取
		clientIp = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
		if clientIp == "" {
			// 直接从 Remote Addr 中获取
			remoteIp, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
			if err != nil {
				clientIp = ""
			} else {
				clientIp = remoteIp
			}
		}
	}
	return clientIp
}

// GetOutBoundIp 获取本机出口IP
func GetOutBoundIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// FileNotExistAndCreate 判断文件是否存在，不存在则创建
func FileNotExistAndCreate(filePath string) (f *os.File, err error) {
	dirPath := filepath.Dir(filePath)
	_, err = os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				return nil, err
			}
			return os.Create(filePath)
		}
	}
	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

// GenerateRandStr 随机产生指定长度字符串
// flag指定类型 0 所有字符 1 全数字 2 全字母 3 数字+字母
func GenerateRandStr(length int, flag int) string {
	arr := make([]string, 0, 16)
	rand.Seed(time.Now().UnixNano())
	var randNum int
	for {
		switch flag {
		case 1:
			randNum = rand.Intn(10) + 48
		case 2:
			rd := rand.Intn(2)
			if rd == 0 {
				randNum = rand.Intn(26) + 65
			} else {
				randNum = rand.Intn(26) + 97
			}
		case 3:
			rd := rand.Intn(3)
			if rd == 0 {
				randNum = rand.Intn(26) + 65
			} else if rd == 1 {
				randNum = rand.Intn(26) + 97
			} else {
				randNum = rand.Intn(10) + 48
			}
		default:
			randNum = rand.Intn(127)
		}

		r := rune(randNum)
		arr = append(arr, string(r))
		if len(arr) == length {
			break
		}
	}
	return strings.Join(arr, "")
}
