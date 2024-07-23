package cmd

import (
	"github.com/rbgayoivoye09/keep-online/src/utils/log"
	"github.com/spf13/cobra"
)

var inputConfigFilePath string

var TrootCmd = &cobra.Command{
	Use: "keep-online", Short: "Keep online commands",
	Run: func(cmd *cobra.Command, args []string) {
		if s, err := cmd.Flags().GetString("config"); err != nil {
			log.Logger.Sugar().Error(err.Error())
		} else {
			inputConfigFilePath = s
		}
		log.Logger.Sugar().Infof("config file path: %s", inputConfigFilePath)
	},
}

func init() {
	TrootCmd.PersistentFlags().StringVarP(&inputConfigFilePath, "config", "c", "", "custom config file path")
	TrootCmd.AddCommand(cfgCmd, sshCmd, cmdCmd, mailCmd, backendCmd)
}
