package cmd

import (
	"time"

	. "github.com/rbgayoivoye09/keep-online/src/utils/log"

	"github.com/rbgayoivoye09/keep-online/src/utils/config"
	"github.com/rbgayoivoye09/keep-online/src/utils/login"
	"github.com/rbgayoivoye09/keep-online/src/utils/ssh"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into a remote server",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.GetConfig()
		t := time.Now()
		s, err := ssh.RemoteFileContent(c.SSH.User, c.SSH.Password, c.SSH.Host, c.SSH.FilePath, c.SSH.Port)
		d := time.Now().Sub(t)
		Logger.Sugar().Infof("SSH took %s", d)
		if err != nil {
			Logger.Sugar().Error(err)
		} else {
			c.User.Password = s
			Logger.Sugar().Info(c.User.Password)
			err = login.AuthenticateVPN(c.Web.LoginURL, c.User.Name, c.User.Password, c.Web.RedirURL)
			if err != nil {
				Logger.Sugar().Error(err)
			}
		}

	},
}
