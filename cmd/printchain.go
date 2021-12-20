/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// printchainCmd represents the printChain command
var printchainCmd = &cobra.Command{
	Use:   "printchain",
	Short: "Print all the blocks of the blockchain",
	Long:  `Print all the blocks of the blockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("printChain called")
		cli.printChain(nodeId)
	},
}

func init() {
	rootCmd.AddCommand(printchainCmd)
	printchainCmd.Flags().UintVarP(&nodeId, "node", "n", 0, "node id")
}
