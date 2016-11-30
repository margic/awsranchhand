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
	"errors"
	"strings"

	"bytes"
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"
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
		log.Info("labelhost called")
		req := labelRequest{
			Host:   viper.GetString("host"),
			Key:    viper.GetString("key"),
			Value:  viper.GetString("value"),
			Add:    viper.GetBool("add"),
			Remove: viper.GetBool("remove"),
		}
		err := req.Validate()
		if err != nil {
			log.Error(err)
		}
		if req.Add && req.Remove {
			log.Error("Cannot add and remove a label at same time")
		} else {
			err = labelHost(req)
			if err != nil {
				log.Error(err)
			}
		}
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
	labelhostCmd.Flags().BoolP("add", "a", false, "Add a label to a host")
	labelhostCmd.Flags().BoolP("remove", "r", false, "Remove a label from a host")

	viper.BindPFlags(labelhostCmd.Flags())
}

func lookupHost(r labelRequest) (host *client.Host, err error) {
	rc, err := getRancherClient()
	if err != nil {
		return nil, err
	}

	host, err = rc.Host.ById(r.Host)
	if err != nil {
		return nil, err
	}

	if host == nil {
		return nil, errors.New("Host " + r.Host + " not found")
	}

	log.WithFields(log.Fields{
		"ID":     r.Host,
		"State":  host.State,
		"Labels": host.Labels,
	}).Debug("Found host")
	return host, nil
}

func listLabels(host *client.Host) map[string]interface{} {
	labels := host.Labels

	if len(labels) > 0 {
		log.WithField("Labels", labels).Info("Current Labels")
	} else {
		log.Info("Host " + host.Id + " has no labels")
	}
	return labels
}

func labelHost(r labelRequest) error {
	host, err := lookupHost(r)
	if err != nil {
		return err
	}
	labels := listLabels(host)

	if !r.Add && !r.Remove {
		// only a  list required just return after listing
		log.Info("Neither -a or -r flags passed. Listing labels only. Exiting")
		return nil
	}

	if r.Add {
		if len(r.Key) == 0 || len(r.Value) == 0 {
			return errors.New("To add a key and value must be specified")
		}
		labels[r.Key] = r.Value
	}

	if r.Remove {
		if len(r.Key) == 0 {
			return errors.New("To remove a key must be specified")
		}
		if labels[r.Key] == nil {
			return errors.New("Key " + r.Key + " not a current label for host " + r.Host)
		}
		delete(labels, r.Key)
	}

	// make the request map to get the json from for update request
	l := make(map[string]string)
	for key, value := range labels {
		l[key] = value.(string)
	}

	lu := labelUpdate{
		Labels: l,
	}

	ub, _ := json.Marshal(lu)
	log.Debug(string(ub))

	// get the host link for update later
	hostLinks := host.Links
	// use this link to update the host with new labels
	selfLink := hostLinks["self"]
	// If this code breaks this is likely where you should be looking!
	if !strings.Contains(selfLink, "projects") {
		// this is the account link not the project link.
		selfLink = strings.Replace(selfLink, "hosts", "projects/"+host.AccountId+"/hosts", 1)
	}

	err = doRancherPut(selfLink, bytes.NewReader(ub))
	if err != nil {
		log.Error("Error updating host with labels: ", err)
	}
	return nil
}

type labelRequest struct {
	Host   string `valid:"required"`
	Key    string
	Value  string
	Add    bool
	Remove bool
}

type labelUpdate struct {
	Labels map[string]string `json:"labels"`
}
