package main

import (
	"os"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"society2chain/app"
)

func main() {
	// Initialize with hardware binding for Society2
	rootCmd := app.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "SOCIETY2", app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}