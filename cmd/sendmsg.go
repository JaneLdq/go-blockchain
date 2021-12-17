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

// sendmsgCmd represents the sendmsg command
var sendmsgCmd = &cobra.Command{
	Use:   "sendmsg",
	Short: "Send a message from a node to another",
	Long:  `Send a message from a node to another`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO trigger a trasaction
		fmt.Println("sendmsg called from ", from, " to", to, "message ", msg)
		cli.sendmsg(from, to, msg)
	},
}

func init() {
	rootCmd.AddCommand(sendmsgCmd)
	sendmsgCmd.Flags().StringVarP(&from, "from", "f", "", "sender node address")
	sendmsgCmd.Flags().StringVarP(&to, "to", "t", "", "receiver node address")
	sendmsgCmd.Flags().StringVarP(&msg, "message", "m", "", "message")
}
