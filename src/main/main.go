package main

import (
	"os"

	. "github.com/rbgayoivoye09/keep-online/src/utils/log"

	"github.com/rbgayoivoye09/keep-online/src/utils/cmd"
)

func main() {
	if err := cmd.TrootCmd.Execute(); err != nil {
		Logger.Sugar().Error(err)
		os.Exit(1)
	}
}
