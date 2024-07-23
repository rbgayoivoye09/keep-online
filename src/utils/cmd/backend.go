package cmd

import (
	"time"

	"github.com/rbgayoivoye09/keep-online/src/utils/config"
	"github.com/rbgayoivoye09/keep-online/src/utils/log"
	"github.com/rbgayoivoye09/keep-online/src/utils/login"
	"github.com/rbgayoivoye09/keep-online/src/utils/ssh"
	"github.com/spf13/cobra"
)

var backendCmd = &cobra.Command{
	Use:   "backend",
	Short: "Backend commands",
	Run: func(cmd *cobra.Command, args []string) {

		go func() {
			for {
				log.Logger.Sugar().Info("Backend commands")

				c := config.GetConfig(inputConfigFilePath)

				t := time.Now()
				s, err := ssh.RemoteFileContent(c.SSH.User, c.SSH.Password, c.SSH.Host, c.SSH.FilePath, c.SSH.Port)
				d := time.Since(t)

				log.Logger.Sugar().Infof("SSH took %s", d)
				if err != nil {
					log.Logger.Sugar().Error(err)
				} else {
					c.User.Password = s
					err = login.AuthenticateVPN(c.Web.LoginURL, c.User.Name, c.User.Password, c.Web.RedirURL)
					if err != nil {
						log.Logger.Sugar().Error(err)
					}
				}

				time.Sleep(10 * time.Second)
			}

		}()

		select {}
	},
}
