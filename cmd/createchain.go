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

var address string

// createchainCmd represents the createchain command
var createchainCmd = &cobra.Command{
	Use:   "createchain",
	Short: "Create a blockchain and send genesis block reward to address",
	Long:  `Create a blockchain and send genesis block reward to address`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createChain called")
		cli.createChain(address)
	},
}

func init() {
	rootCmd.AddCommand(createchainCmd)
	createchainCmd.Flags().StringVarP(&address, "address", "a", "", "address to send genesis block reward to")
}
