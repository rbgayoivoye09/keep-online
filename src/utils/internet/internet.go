package internet

import (
	"net/http"

	. "github.com/rbgayoivoye09/keep-online/src/utils/log"
)

func CheckInternetAccess() bool {
	// 创建一个HTTP客户端
	client := &http.Client{}

	// 构建GET请求
	url := "http://www.baidu.com"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Logger.Sugar().Info("Error creating request:", err)
		return false
	}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		Logger.Sugar().Info("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode == http.StatusOK {
		Logger.Sugar().Info("Connected to the internet!")
		return true
	} else {
		Logger.Sugar().Info("Failed to connect to the internet. Status code:", resp.StatusCode)
		return false
	}
}
