package cmd

import (
	"github.com/spf13/cobra"
	"go-ccsp/components/mta"
)

func init() {
	mtaCmd.Flags().StringP("dml-config", "c", "/usr/ccsp/mta/mta_dml_config.json", "DML configuration file path")
}

var mtaCmd = &cobra.Command{
	Use:   "mta",
	Short: "Run the MTA subsystem",
	Long:  `Run the MTA subsystem of the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		mtaAgent := mta.New(&mta.MTAAgentConfig{
			DMLConfigPath:        cmd.Flag("dml-config").Value.String(),
			EthernetWANEnabled:   true,
			ErouterDHCPOptionMTA: true,
		})
		err := mtaAgent.Run("eRT")
		if err != nil {
			panic(err)
		}
	},
}
