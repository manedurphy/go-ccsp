package cmd

import (
	"github.com/spf13/cobra"
	"go-ccsp/components/mta"
)

var mtaCmd = &cobra.Command{
	Use:   "mta",
	Short: "Run the MTA subsystem",
	Long:  `Run the MTA subsystem of the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := mta.Run("eRT")
		if err != nil {
			panic(err)
		}
	},
}
