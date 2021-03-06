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

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use: "connect",
	Short: "Connect a node to another",
	Long: `Connect a node to another`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("connect called with args (from = %s, to = %s)\n", from, to)
		cli.connectNodes(from, to)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringVarP(&from, "from", "f", "", "sender node host address")
	connectCmd.Flags().StringVarP(&to, "to", "t", "", "receiver node host address")
}