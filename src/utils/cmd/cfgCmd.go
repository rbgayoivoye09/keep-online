package cmd

import (
	"github.com/rbgayoivoye09/keep-online/src/utils/config"
	"github.com/rbgayoivoye09/keep-online/src/utils/internet"
	"github.com/rbgayoivoye09/keep-online/src/utils/log"
	"github.com/rbgayoivoye09/keep-online/src/utils/login"

	"github.com/spf13/cobra"
)

var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "Configure keep-online settings",
	Run: func(cmd *cobra.Command, args []string) {

		if internet.CheckInternetAccess() {
			return
		}

		c := config.GetConfig(inputConfigFilePath)

		err := login.AuthenticateVPN(c.Web.LoginURL, c.User.Name, c.User.Password, c.Web.RedirURL)
		if err != nil {
			log.Logger.Sugar().Error(err)
		}
	},
}
