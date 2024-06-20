/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"os"
	"path"

	"github.com/abibby/jit/git"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jit",
	Short: "Jira + Git integration",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	viper.SetDefault("branch_prefix", "")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		config, err := os.UserConfigDir()
		if err == nil {
			viper.AddConfigPath(path.Join(config, "jit"))
		}
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(path.Join(home, ".config", "jit"))
		}
	}

	viper.SetConfigName("config")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	// viper.ReadInConfig()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	gitRoot, err := git.Root()
	if err == nil {
		viper.AddConfigPath(path.Join(gitRoot, ".jit"))
		if err := viper.MergeInConfig(); err == nil {
			fmt.Println("Adding local config file:", viper.ConfigFileUsed())
		}
	}
}
