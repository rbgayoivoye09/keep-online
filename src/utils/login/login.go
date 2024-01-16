package login

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	. "github.com/rbgayoivoye09/keep-online/src/utils/log"
)

// checkInternetAccess 检测当前环境是否可以接入互联网
func checkInternetAccess() bool {
	urt := "www.baidu.com:80"
	// urt := "www.google.com:80"
	_, err := net.DialTimeout("tcp", urt, 5*time.Second)
	if err != nil {
		Logger.Sugar().Error("No internet access ", urt, err.Error())
		return false
	}
	Logger.Sugar().Info("Internet access available ", urt)
	return true
}

func AuthenticateVPN(loginUrl, authUser, authPass, redirectUrl string) error {
	if !checkInternetAccess() {
		Logger.Sugar().Warn("无网络访问，执行 VPN 认证")
		return _authenticateVPN(loginUrl, authUser, authPass, redirectUrl)
	}
	Logger.Sugar().Info("接入网络，跳过 VPN 认证")
	return nil
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
