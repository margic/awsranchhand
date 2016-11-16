// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ranchhand",
	Short: "Ranch hand helps out on the ranch",
	Long:  `Ranch hand is a set of cli tools for supporting a rancher cluster`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	log.Debug("Adding Flags")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ranchhand.yml)")
	RootCmd.PersistentFlags().StringP("rancherKey", "k", "", "Rancher API key")
	RootCmd.PersistentFlags().StringP("rancherSecret", "s", "", "Rancher api secret")
	viper.BindPFlag("rancherKey", RootCmd.Flags().Lookup("rancherKey"))
	viper.BindPFlag("rancherSecret", RootCmd.Flags().Lookup("rancherSecret"))
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// log.Debug("Rancher Key: ", viper.GetString("rancherKey"))
	// rancherSecret := viper.GetString("rancherSecret")
	// // mask the key for debug output
	// l := len(rancherSecret) - 10
	// if l > 0 {
	// 	rp := strings.Repeat("*", l)
	// 	replacer := strings.NewReplacer(rancherSecret[0:l], rp)
	// 	log.Debug("Rancher Secret: ", replacer.Replace(rancherSecret))
	// } else {
	// 	log.Debug("Rancher Secret value is less than 10 chars cannot be valid value: ", rancherSecret)
	// }
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}
	viper.SetConfigType("yaml")
	viper.SetConfigName("ranchhand")
	viper.AddConfigPath("/etc/ranchhand/")
	viper.AddConfigPath("$HOME/.ranchhand")
	viper.AddConfigPath("./")

	viper.BindEnv("rancherKey", "RANCHER_KEY") // need to bind name key to env KEY
	viper.BindEnv("rancherSecret", "RANCHER_SECRET")
	viper.BindEnv("url", "RANCHER_URL")
	viper.BindEnv("logginglevel", "LOGGING_LEVEL")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config")
	}

	level, err := log.ParseLevel(strings.ToLower(viper.GetString("loggingLevel")))
	if err != nil {
		log.Error("Invalid logging level")
	} else {
		log.SetLevel(level)
		log.Debug("Logging level: ", level)
		log.Debug("Using config file: ", viper.ConfigFileUsed())
	}
}
