package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var address string

// getbalanceCmd represents the getbalance command
var getbalanceCmd = &cobra.Command{
	Use:   "getbalance",
	Short: "Get the balance of the given address",
	Long:  `Get the balance of the given address`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getbalance called for address", address)
		cli.getBalance(address, nodeId)
	},
}

func init() {
	rootCmd.AddCommand(getbalanceCmd)
	getbalanceCmd.Flags().StringVarP(&address, "address", "a", "", "address")
	getbalanceCmd.Flags().UintVarP(&nodeId, "node", "n", 0, "node id")
	getbalanceCmd.MarkFlagRequired("address")
	getbalanceCmd.MarkFlagRequired("node")
}
