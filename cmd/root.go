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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cmdLineArgs service.CmdLineArgs

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ipakpublisher",
	Short: "Publish your app to app store or google play",
	Long: `When you use the upload tool provided by IOS to upload the app to the app store.
For example, Xcode, altool, transporter.
You can then use this program to complete the subsequent review steps.
At the same time, it also uses the Google play developer API to upload and distribute Android apps.
It is divided into two steps: uploading application files and publishing applications.
In this program, you can choose whether to publish after uploading APK.
For more information, see command usage and parameters.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		com.Log.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cmdLineArgs.ConfigFile, "config", "c", "", "config file for app store connect api auth (default is $HOME/.ipakpublisher.toml)")
	rootCmd.PersistentFlags().StringVarP(&cmdLineArgs.WorkDir, "workdir", "w", "", "workdir path (default is $(pwd))")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cmdLineArgs.ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cmdLineArgs.ConfigFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			com.Log.Error(err)
			os.Exit(1)
		}
		// Search config in home directory with name ".ipakpublisher" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".ipakpublisher")
	}
	viper.SetEnvPrefix("at")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		com.Log.Info("Using config file:", viper.ConfigFileUsed())
	}
	// init config
	if err := service.InitConfig(&cmdLineArgs); err != nil {
		com.Log.Error("init config error.", err)
		os.Exit(1)
	}
}
