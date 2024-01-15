package login

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	. "github.com/rbgayoivoye09/keep-online/src/utils/log"

	"time"
)

// CheckInternetConnection 检测当前环境是否可以接入互联网
func CheckInternetConnection() bool {
	// 设置超时时间为2秒
	timeout := time.Second * 2

	// 尝试连接到一个已知的互联网地址（例如Google的公共DNS服务器）
	conn, err := net.DialTimeout("tcp", "114.114.114.114:80", timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

func AuthenticateVPN(loginUrl, authUser, authPass, redirectUrl string) error {
	return _authenticateVPN(loginUrl, authUser, authPass, redirectUrl)
}

func _authenticateVPN(loginUrl, authUser, authPass, redirectUrl string) error {
	// 发送GET请求获取网页

	response, err := http.Get(loginUrl)
	if err != nil {
		return fmt.Errorf("获取网页失败: %v", err)
	}
	defer response.Body.Close()

	// 构造认证信息并提交表单
	postURL := loginUrl
	payload := url.Values{
		"auth_user": {authUser},
		"auth_pass": {authPass},
		"redirurl":  {redirectUrl}, // 你可能需要根据实际情况修改这些值
		"accept":    {"登录"},
	}

	response, err = http.PostForm(postURL, payload)
	if err != nil {
		return fmt.Errorf("提交表单失败: %v", err)
	}
	defer response.Body.Close()

	// 检查认证是否成功
	if response.StatusCode == http.StatusOK {
		Logger.Sugar().Info("认证成功！")
	} else {
		return fmt.Errorf("认证失败，状态码: %d", response.StatusCode)
	}

	return nil
}
