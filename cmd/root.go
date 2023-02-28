package cmd

import (
	"os"
	"time"

	"go-boilerplate/common"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
)

// RootCmd holds reference to root cmd to be used by children commands
var RootCmd = &cobra.Command{
	Use:   "go-boilerplate",
	Short: "go Boilerplate API",
	Long:  "go Boilerplate API",
}

// Execute executes root cmd
func Execute() {
	defer func() {
		err := recover()
		if err != nil {
			casted, ok := err.(error)
			if ok {
				common.HandleError("unexpected error while executing command", casted)
			}
		}
	}()

	err := RootCmd.Execute()
	if err != nil {
		common.HandleError("error while executing command", err)
		sentry.Flush(time.Second * 2)
		os.Exit(-1)
	}
}
