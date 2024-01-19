package internet

import (
	"fmt"
	"io"
	"net/http"
	"time"

	. "github.com/rbgayoivoye09/keep-online/src/utils/log"
)

func CheckInternetAccess() bool {
	// 创建一个HTTP客户端
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	// 构建GET请求
	url := "https://www.baidu.com"
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

	// 读取响应的内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	} else {
		Logger.Sugar().Info("Response body:", string(body))
	}

	// 检查响应状态码
	if resp != nil && resp.StatusCode == http.StatusOK {
		Logger.Sugar().Info("Connected to the internet! resp.StatusCode ", resp.StatusCode)
		return true
	} else {
		Logger.Sugar().Info("Failed to connect to the internet. Status code:", resp.StatusCode)
		return false
	}
}
