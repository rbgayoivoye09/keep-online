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

// AuthenticateVPN authenticates the user for VPN access.
//
// Parameters:
// - loginUrl: the URL of the login page.
// - authUser: the username for authentication.
// - authPass: the password for authentication.
// - redirectUrl: the URL to redirect to after successful authentication.
//
// Returns:
// - error: an error if authentication fails.
func AuthenticateVPN(loginUrl, authUser, authPass, redirectUrl string) error {
	return _authenticateVPN(loginUrl, authUser, authPass, redirectUrl)
}

// _authenticateVPN is an internal function that actually authenticates the user for VPN access.
func _authenticateVPN(loginUrl, authUser, authPass, redirectUrl string) error {
	Logger.Sugar().Infof("开始认证VPN... %s %s %s %s", loginUrl, authUser, authPass, redirectUrl)
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
		checkInternetAccess()
	} else {
		return fmt.Errorf("认证失败，状态码: %d", response.StatusCode)
	}

	return nil
}
