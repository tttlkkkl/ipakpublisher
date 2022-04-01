/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"os"

	"github.com/spf13/cobra"
	"github.com/tttlkkkl/ipakpublisher/com"
	"github.com/tttlkkkl/ipakpublisher/service"
)

var appleArgs service.AppleCmdArgs

// initCmd represents the init command
var ipaCmd = &cobra.Command{
	Use:   "ipa",
	Short: "Perform IPA related operations",
	Long:  `Perform IPA related operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := service.NewAppleService(&cmdLineArgs, &appleArgs)
		if err != nil {
			com.Log.Error(err)
			os.Exit(1)
		}
		if err := svc.Ipa(); err != nil {
			com.Log.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(ipaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	ipaCmd.Flags().BoolVarP(&appleArgs.Init, "init", "i", false, "If you specify this option, metadata is downloaded from the API and an attempt is made to overwrite the local record.When this operation is completed, the program terminates.")
	ipaCmd.Flags().StringVarP(&appleArgs.BundleID, "bundle-id", "b", "", "Specify the ipa bundleid.")
	ipaCmd.Flags().StringVarP(&appleArgs.Platform, "platform", "p", "IOS", "Specify the platform.")
}
