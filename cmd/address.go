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

// addressCmd represents the address command
var addressCmd = &cobra.Command{
	Use:   "address",
	Short: "Display the address of the given node",
	Long:  `Display the address of the given node`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO query node address
		fmt.Println("address called for node", nodeId)
	},
}

func init() {
	rootCmd.AddCommand(addressCmd)
	addressCmd.Flags().UintVarP(&nodeId, "node", "n", 0, "node id")
	addressCmd.MarkFlagRequired("node")
}
