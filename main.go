package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCLI(c MyConfig) {
	var myconfig MyConfig
	logger.Sugar().Infof("c: %v\n", c)

	if c.UserName == "" || c.UserLoginUrl == "" || c.UserPassword == "" || c.UserRedirul == "" {
		logger.Sugar().Warn("命令行用户名,密码,登录地址,重定向地址 为空")
		logger.Sugar().Warn("使用配置文件")
		// 调用函数读取配置文件
		myconfig, err = readConfigFile()
		if err != nil {
			logger.Sugar().Error(err)
		} else {
			logger.Sugar().Info("读取配置文件成功")
		}
	} else {
		logger.Sugar().Info("使用命令行输入的配置")
		myconfig = c
	}

	// 打印配置项
	logger.Sugar().Info("User Name:", myconfig.UserName)
	logger.Sugar().Info("User Password:", myconfig.UserPassword)
	logger.Sugar().Info("User LoginUrl:", myconfig.UserLoginUrl)
	logger.Sugar().Info("User Redirul:", myconfig.UserRedirul)

	// login()
	err = authenticateVPN(myconfig.UserLoginUrl, myconfig.UserName, myconfig.UserPassword, myconfig.UserRedirul)
	if err != nil {
		logger.Sugar().Error(err)
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

type MyConfig struct {
	UserName     string `mapstructure:"user_name"`
	UserPassword string `mapstructure:"user_password"`
	UserRedirul  string `mapstructure:"user_redirurl"`
	UserLoginUrl string `mapstructure:"user_login_url"`
}

func readConfigFile() (MyConfig, error) {
	filePath := "user.yml"

	// 设置配置文件名和路径
	viper.SetConfigFile(filePath)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return MyConfig{}, fmt.Errorf("Fatal error config file: %s", err)
	}

	// 解析配置文件到结构体
	var config MyConfig
	if err := viper.Unmarshal(&config); err != nil {
		return MyConfig{}, fmt.Errorf("Unable to decode into struct: %s", err)
	}

	if config.UserName == "" || config.UserPassword == "" || config.UserRedirul == "" || config.UserLoginUrl == "" {
		return MyConfig{}, fmt.Errorf("配置文件中缺少必要的配置项")
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
		logger.Sugar().Info("认证成功！")
	} else {
		return fmt.Errorf("认证失败，状态码: %d", response.StatusCode)
	}

	return nil
}

var logger *zap.Logger
var err error

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "A simple CLI tool",
	Long:  `A simple CLI tool built with Cobra`,
	Run: func(cmd *cobra.Command, args []string) {
		user_name, _ := cmd.Flags().GetString("user_name")
		user_password, _ := cmd.Flags().GetString("user_password")
		user_redirurl, _ := cmd.Flags().GetString("user_redirurl")
		user_login_url, _ := cmd.Flags().GetString("user_login_url")

		runCLI(MyConfig{UserName: user_name, UserPassword: user_password, UserRedirul: user_redirurl, UserLoginUrl: user_login_url})

	},
}

func init() {
	rootCmd.Flags().String("user_name", "", "User name")
	rootCmd.Flags().String("user_password", "", "User password")
	rootCmd.Flags().String("user_redirurl", "", "Redir URL")
	rootCmd.Flags().String("user_login_url", "", "Login URL")

	// 配置日志文件的路径和其他相关参数
	logDirectory := "./logs/"
	logFile := logDirectory + "app.log"
	maxSize := 10 // MB
	maxBackups := 5
	maxAge := 7 // days

	// 创建日志目录
	if err := os.MkdirAll(logDirectory, os.ModePerm); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	// 创建一个 lumberjack.Logger，用于处理日志轮换
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	}
	// 创建 zap 的配置
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", logFile}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 配置 zap logger
	logger, err = config.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.ErrorOutput(zapcore.AddSync(lumberjackLogger)),
	)
	if err != nil {
		panic("Failed to initialize zap logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("This is a log message.")

}
