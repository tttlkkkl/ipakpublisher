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

// apkCmd represents the apk command
var apkCmd = &cobra.Command{
	Use:   "apk",
	Short: "push the apk file",
	Long:  `Use this command when you want to publish an app in APK format`,
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := service.NewGoogleService(&cmdLineArgs, &googleArgs)
		if err != nil {
			com.Log.Error(err)
			os.Exit(1)
		}
		if err := svc.Do("apk"); err != nil {
			com.Log.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(apkCmd)
	initGoogleCmd(apkCmd)
}
