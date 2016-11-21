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
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// labelec2instanceCmd represents the labelec2instance command
var labelec2instanceCmd = &cobra.Command{
	Use:   "labelec2instance",
	Short: "Label Rancher host with EC2 Instance ID",
	Long: ` This command uses the aws ec2 meta service to look up the instance
	id for this host and create a label on the rancher host with that value`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Info("labelec2instance called")
		instanceID, err := lookupInstanceID()
		if err != nil {
			log.Error("Unable to lookup instanceID from ec2 meta service")
		}

		req := labelRequest{
			Host:   viper.GetString("host"),
			Key:    "isntance-id",
			Value:  instanceID,
			Add:    viper.GetBool("add"),
			Remove: viper.GetBool("remove"),
		}

		err = labelHost(req)
		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(labelec2instanceCmd)
}
