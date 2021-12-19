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

var port uint

// newnodeCmd represents the newnode command
var newnodeCmd = &cobra.Command{
	Use:   "newnode",
	Short: "Create a new node with given port",
	Long:  `Create a new node with given port`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("newnode called with args (port = %d)\n", port)
		cli.newNode(port)
	},
}

func init() {
	rootCmd.AddCommand(newnodeCmd)
	newnodeCmd.Flags().UintVarP(&port, "port", "p", 0, "port to run node on, which is also the node id")
	newnodeCmd.MarkFlagRequired("port")
}
