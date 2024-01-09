package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/viper"
)

func main() {
	// P()
	if CheckInternetConnection() {
		fmt.Println("已成功连接到互联网!")

		// 调用函数读取配置文件
		config, err := readConfigFile()
		if err != nil {
			panic(err)
		}

		// 打印配置项
		fmt.Printf("User Name: %s\n", config.UserName)
		fmt.Printf("User Password: %s\n", config.UserPassword)
		fmt.Printf("User LoginUrl: %s\n", config.UserLoginUrl)
		fmt.Printf("User Redirul: %s\n", config.UserRedirul)

		// login()
		err = authenticateVPN(config.UserLoginUrl, config.UserName, config.UserPassword, config.UserRedirul)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		fmt.Println("无法连接到互联网.")

		// 调用函数读取配置文件
		config, err := readConfigFile()
		if err != nil {
			panic(err)
		}

		// 打印配置项
		fmt.Printf("User Name: %s\n", config.UserName)
		fmt.Printf("User Password: %s\n", config.UserPassword)
		fmt.Printf("User LoginUrl: %s\n", config.UserLoginUrl)
		fmt.Printf("User Redirul: %s\n", config.UserRedirul)

		// login()
		err = authenticateVPN(config.UserLoginUrl, config.UserName, config.UserPassword, config.UserRedirul)
		if err != nil {
			fmt.Println(err)
		}
	}
}

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

type Config struct {
	UserName     string `mapstructure:"user_name"`
	UserPassword string `mapstructure:"user_password"`
	UserRedirul  string `mapstructure:"user_redirurl"`
	UserLoginUrl string `mapstructure:"user_login_url"`
}

func readConfigFile() (Config, error) {
	filePath := "user.yml"

	// 设置配置文件名和路径
	viper.SetConfigFile(filePath)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("Fatal error config file: %s", err)
	}

	// 解析配置文件到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("Unable to decode into struct: %s", err)
	}

	return config, nil
}

func authenticateVPN(loginUrl, authUser, authPass, redirectUrl string) error {
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
		fmt.Println("认证成功！")
	} else {
		return fmt.Errorf("认证失败，状态码: %d", response.StatusCode)
	}

	return nil
}

func P() {
	// 设置Ping的目标地址
	target := "www.baidu.com"

	// 设置超时时间
	timeout := 2 * time.Second

	// 执行Ping操作
	err := ping(target, timeout)
	if err != nil {
		fmt.Println("Ping失败:", err)
		os.Exit(1)
	}

	fmt.Println("Ping成功，已连接到互联网！")
}

func ping(target string, timeout time.Duration) error {
	// 创建一个icmp类型的网络连接
	conn, err := net.DialTimeout("ip4:icmp", target, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 获取本机IP地址
	localAddr := conn.LocalAddr().(*net.IPAddr)

	// 构建ICMP消息
	icmpMsg := []byte{8, 0, 0, 0, 0, 13, 0, 37, byte(os.Getpid() & 0xff), byte(os.Getpid() >> 8)}

	// 发送Ping消息
	_, err = conn.Write(icmpMsg)
	if err != nil {
		return err
	}

	// 设置读取超时时间
	conn.SetReadDeadline(time.Now().Add(timeout))

	// 接收Ping响应
	receive := make([]byte, 28+len(icmpMsg))
	_, err = conn.Read(receive)
	if err != nil {
		return err
	}

	fmt.Printf("接收到的数据: %v\n", receive)
	fmt.Printf("Ping成功，来自 %s 的回复。\n", localAddr.IP.String())

	return nil
}
