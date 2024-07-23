package cmd

import (
	"github.com/rbgayoivoye09/keep-online/src/utils/config"
	"github.com/rbgayoivoye09/keep-online/src/utils/internet"
	"github.com/rbgayoivoye09/keep-online/src/utils/log"
	"github.com/rbgayoivoye09/keep-online/src/utils/login"

	"github.com/spf13/cobra"
)

var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Execute a custom command",
	Run: func(cmd *cobra.Command, args []string) {

		if internet.CheckInternetAccess() {
			return
		}

		userName, _ := cmd.Flags().GetString("user_name")
		userPassword, _ := cmd.Flags().GetString("user_password")
		userRedirURL, _ := cmd.Flags().GetString("user_redirurl")
		userLoginURL, _ := cmd.Flags().GetString("user_login_url")

		log.Logger.Sugar().Infof("User Name: %s\n", userName)
		log.Logger.Sugar().Infof("User RedirURL: %s\n", userRedirURL)
		log.Logger.Sugar().Infof("User LoginURL: %s\n", userLoginURL)

		u := config.User{
			Name:     userName,
			Password: userPassword,
		}
		runCLI(config.Config{User: u})
	},
}

// init initializes the necessary flags for the cmdCmd function.
//
// It adds the required parameters:
// - user_name: User name (required)
// - user_password: User password (required)
//
// It also adds the optional parameters:
// - user_redirurl: User redirurl (optional)
// - user_login_url: User login URL (optional)
func init() {
	// 添加必选参数
	cmdCmd.Flags().StringP("user_name", "u", "", "User name (required)")
	err := cmdCmd.MarkFlagRequired("user_name")
	if err != nil {
		log.Logger.Sugar().Error(err.Error())
	}
	cmdCmd.Flags().StringP("user_password", "p", "", "User password (required)")
	err = cmdCmd.MarkFlagRequired("user_password")
	if err != nil {
		log.Logger.Sugar().Error(err.Error())
	}
	// 添加可选参数
	cmdCmd.Flags().String("user_redirurl", "", "User redirurl (optional)")
	cmdCmd.Flags().String("user_login_url", "", "User login URL (optional)")
}

func runCLI(c config.Config) {
	var myconfig config.Config
	var err error
	log.Logger.Sugar().Infof("c: %v\n", c)
	myconfig = c

	if c.Web.LoginURL == "" || c.Web.RedirURL == "" {
		log.Logger.Sugar().Warn("登录地址,重定向地址 为空")
		log.Logger.Sugar().Warn("使用配置文件")

		// 调用函数读取配置文件
		rconfig := config.GetConfig(inputConfigFilePath)
		myconfig.Web.LoginURL = rconfig.Web.LoginURL
		myconfig.Web.RedirURL = rconfig.Web.RedirURL

		log.Logger.Sugar().Infof("User RedirURL: %s\n", myconfig.Web.LoginURL)
		log.Logger.Sugar().Infof("User LoginURL: %s\n", myconfig.Web.RedirURL)
	} else {
		log.Logger.Sugar().Info("使用命令行输入的配置")
		myconfig = c
	}

	// 打印配置项
	log.Logger.Sugar().Info("User Name:", myconfig.User.Name)
	log.Logger.Sugar().Info("User LoginUrl:", myconfig.Web.LoginURL)
	log.Logger.Sugar().Info("User Redirul:", myconfig.Web.RedirURL)

	// login()
	err = login.AuthenticateVPN(myconfig.Web.LoginURL, myconfig.User.Name, myconfig.User.Password, myconfig.Web.RedirURL)
	if err != nil {
		log.Logger.Sugar().Error(err)
	}

}
