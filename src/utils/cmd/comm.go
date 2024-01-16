package cmd

import (
	"net"
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
