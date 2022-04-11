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

var googleArgs service.GoogleCmdArgs

// aabCmd represents the aab command
var aabCmd = &cobra.Command{
	Use:   "aab",
	Short: "push the aab file",
	Long:  `Use this command when you want to publish an app in AAB format`,
	Run: func(cmd *cobra.Command, args []string) {
		svc, err := service.NewGoogleService(&cmdLineArgs, &googleArgs)
		if err != nil {
			com.Log.Error(err)
			os.Exit(1)
		}
		if err := svc.Do("aab"); err != nil {
			com.Log.Error(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(aabCmd)
	initGoogleCmd(aabCmd)
}
func initGoogleCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&googleArgs.FilPath, "file", "f", "", "the app file path.")
	cmd.PersistentFlags().StringVarP(&googleArgs.Platform, "platform", "p", "android", "Specify the platform.")
	cmd.PersistentFlags().BoolVarP(&googleArgs.Init, "init", "i", false, "If you specify this option, A multilingual change log configuration template will be generated locally. Due to API limitations, this method will not pull data from the remote.")
	cmd.PersistentFlags().BoolVarP(&googleArgs.Sync, "sync", "s", false, "If you specify this option, To simultaneously transfer the local data to the remote location through API.")
	cmd.PersistentFlags().StringVarP(&googleArgs.PackageName, "package-name", "n", "", "app package name")
	cmd.PersistentFlags().StringSliceVarP(&googleArgs.Track, "track", "t", []string{}, `If this option is set.
	it will be published to these tracks after the application upload is completed.
	You can specify multiple times.
	The optional values are:
	alpha
	beta
	internal
	production`)
	cmd.PersistentFlags().StringVar(&googleArgs.Timeout, "timeout", "2m", "Set upload timeout")
	cmd.PersistentFlags().StringVar(&googleArgs.ReleaseNotesFilePath, "notes-file", "", `This option loads a text content as the update log.`)
}
