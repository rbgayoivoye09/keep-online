package cmd

import (
	"github.com/spf13/cobra"
)

var TrootCmd = &cobra.Command{Use: "keep-online", Short: "Keep online commands"}

func init() {
	TrootCmd.AddCommand(cfgCmd, sshCmd, cmdCmd, mailCmd)
}
