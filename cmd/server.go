package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Run:   serverRun,
}

func serverRun(cmd *cobra.Command, args []string) {

}
