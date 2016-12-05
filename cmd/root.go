// Copyright © 2016 Paul Crofts pmcrofts@margic.com
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

const rKey string = "key"
const rSecret string = "secret"
const rURL string = "url"

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "awsranchhand",
	Short: "Ranch hand helps out on the ranch",
	Long:  `Ranch hand is a set of cli tools for supporting a rancher cluster`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) { }
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
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ranchhand.yml)")
	RootCmd.PersistentFlags().String(rKey, "", "Rancher API key")
	RootCmd.PersistentFlags().String(rSecret, "", "Rancher API secret")
	RootCmd.PersistentFlags().String(rURL, "", "Rancher API url")
	RootCmd.PersistentFlags().StringP("stack", "k", "", "Name of the application stack")
	RootCmd.PersistentFlags().StringP("service", "s", "", "Name of service to upgrade")
	RootCmd.PersistentFlags().StringP("loggingLevel", "l", "INFO", "Logging level DEBUG, INFO etc.")
	viper.BindPFlags(RootCmd.PersistentFlags())
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

	viper.BindEnv(rKey, "CATTLE_ACCESS_KEY") // need to bind name key to env KEY
	viper.BindEnv(rSecret, "CATTLE_SECRET_KEY")
	viper.BindEnv(rURL, "CATTLE_URL")
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

	log.Debug("Rancher Key: ", viper.GetString(rKey))
	rancherSecret := viper.GetString(rSecret)
	// mask the key for debug output
	l := len(rancherSecret) - 10
	if l > 0 {
		rp := strings.Repeat("*", l)
		replacer := strings.NewReplacer(rancherSecret[0:l], rp)
		log.Debug("Rancher Secret: ", replacer.Replace(rancherSecret))
	} else {
		log.Debug("Rancher Secret value is less than 10 chars cannot be valid value: ", rancherSecret)
	}
}
