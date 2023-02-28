package cmd

import (
	"go-boilerplate/api"

	"github.com/spf13/cobra"
)

var (
	apiCommand = &cobra.Command{
		Use:   "api",
		Short: "Initializes the API",
		Long:  "Initializes the API.",
		RunE:  apiExecute,
	}
)

func init() {
	RootCmd.AddCommand(apiCommand)
}

func apiExecute(cmd *cobra.Command, args []string) error {
	api.Setup()
	return nil
}
