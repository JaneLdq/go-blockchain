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

// initchainCmd represents the initchain command
var initchainCmd = &cobra.Command{
	Use:   "initchain",
	Short: "Init the blockchain on given node",
	Long:  `Init the blockchain on given node`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO init blockchain on the node
		fmt.Println("initchain called on node", nodeId)
	},
}

func init() {
	rootCmd.AddCommand(initchainCmd)
	initchainCmd.Flags().UintVarP(&nodeId, "node", "n", 0, "node id, idenfical to the port which the node runs on")
	initchainCmd.MarkFlagRequired("node")
}
