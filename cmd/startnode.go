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

// startnodeCmd represents the startnode command
var startnodeCmd = &cobra.Command{
	Use:   "startnode",
	Short: "Start the node running on given port",
	Long:  `Start the node running on given port`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO start the node
		fmt.Println("startnode called on node", nodeId)
	},
}

func init() {
	rootCmd.AddCommand(startnodeCmd)
	startnodeCmd.Flags().UintVarP(&nodeId, "node", "n", 0, "node id, idenfical to the port which the node runs on")
	startnodeCmd.MarkFlagRequired("node")
}
