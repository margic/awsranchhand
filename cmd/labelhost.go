//go:generate gencheck -f=labelhost.go

// Package cmd provides rancher commands for supporting rancher in aws
// Copyright © 2016 pmcrofts@margic.com
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
	"strings"

	"bytes"
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// labelhost should invoke put on host resource on the api.
// This curl works
// curl -u key:secret -X PUT -H 'Content-Type: application/json' \
// --data-ascii '{"labels":{"io.rancher.host.docker_version":"1.12","io.rancher.host.linux_kernel_version":"4.4","mykey":"myval"}}' \
//  https://server/v1/projects/1a22/hosts/1h173

// if the host links are retrieved at account level the urls will be in the
// form https://server/v1/hosts/1h173
// rather than https://server/v1/projects/1a22/hosts/1h173
// We can't update unless we hit the project url see code below where account is used to create project link

// labelhostCmd represents the labelhost command
var labelhostCmd = &cobra.Command{
	Use:   "labelhost",
	Short: "Apply a label to a rancher host",
	Long: `Used to label a specific rancher host with a supplied label.
				host id should be the rancher host id. Label is supplied
				as key and value flags`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("labelhost called")
		req := labelRequest{
			Host:  viper.GetString("host"),
			Key:   viper.GetString("key"),
			Value: viper.GetString("value"),
		}
		err := req.Validate()
		if err != nil {
			log.Error(err)
		}
		labelHost(req)
	},
}

func init() {
	RootCmd.AddCommand(labelhostCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// labelhostCmd.PersistentFlags().String("foo", "", "A help for foo")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	labelhostCmd.Flags().String("host", "", "Rancher host id")
	labelhostCmd.Flags().String("key", "", "Rancher host label key")
	labelhostCmd.Flags().String("value", "", "Rancher host label value")
	viper.BindPFlags(labelhostCmd.Flags())
}

func labelHost(r labelRequest) {
	rc, err := GetRancherClient()
	if err != nil {
		log.Error("Error getting client: ", err)
		return
	}

	host, err := rc.Host.ById(r.Host)
	if err != nil {
		log.Error("Error checking if host exists: ", err)
		return
	}

	// get the host link for update later
	hostLinks := host.Links
	// use this link to update the host with new labels
	selfLink := hostLinks["self"]

	log.WithFields(log.Fields{
		"ID":     r.Host,
		"State":  host.State,
		"Labels": host.Labels,
		"Link":   selfLink,
	}).Debug("Found host")

	// If this code breaks this is likely where you should be looking!
	if !strings.Contains(selfLink, "projects") {
		// this is the account link not the project link.
		selfLink = strings.Replace(selfLink, "hosts", "projects/"+host.AccountId+"/hosts", 1)
	}

	labels := host.Labels
	labels[r.Key] = r.Value
	log.Debug(labels)

	l := make(map[string]string)

	for key, value := range labels {
		l[key] = value.(string)
	}

	lu := labelUpdate{
		Labels: l,
	}

	ub, _ := json.Marshal(lu)
	log.Debug(string(ub))

	err = doRancherPut(selfLink, bytes.NewReader(ub))
	if err != nil {
		log.Error("Error updating host with labels: ", err)
	}
}

type labelRequest struct {
	Host  string `valid:"required"`
	Key   string `valid:"required"`
	Value string `valid:"required"`
}

type labelUpdate struct {
	Labels map[string]string `json:"labels"`
}
