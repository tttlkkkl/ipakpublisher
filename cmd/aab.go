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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tttlkkkl/ipakpublisher/service"
)

var googleArgs service.GoogleCmdArgs

// aabCmd represents the aab command
var aabCmd = &cobra.Command{
	Use:   "aab",
	Short: "push the aab file",
	Long:  `Use this command when you want to publish an app in AAB format`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aab called")
	},
}

func init() {
	rootCmd.AddCommand(aabCmd)
	initGoogleCmd(aabCmd)
}
func initGoogleCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&googleArgs.FilPath, "file", "f", "", "app file path")
	cmd.PersistentFlags().StringVarP(&googleArgs.PackageName, "package-name", "p", "", "app package name")
	cmd.PersistentFlags().StringSliceVarP(&googleArgs.Track, "track", "t", []string{}, `If this option is set.
	it will be published to these tracks after the application upload is completed.
	You can specify multiple times.
	The optional values are:
	alpha
	beta
	internal
	production`)
	cmd.PersistentFlags().StringVarP(&googleArgs.ReleaseNotes, "notes", "n", "", "This option specifies the update log")
	cmd.PersistentFlags().StringVar(&googleArgs.Timeout, "timeout", "10s", "Set upload timeout")
	cmd.PersistentFlags().StringVar(&googleArgs.ReleaseNotesFilePath, "notes-file", "", `This option loads a text content as the update log.
	The current option overrides the notes option.`)
}
