//go:generate gencheck -f=rollback.go
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/Sirupsen/logrus"
)

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback a rancher service upgrade",
	Long: `Rollback a rancher service undoing an upgrade.
	Call this command if the verification fails and the service should be rolled back`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("rollback called")
		err := rollback()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(rollbackCmd)
}

func rollback() error {
	r := rollbackRequest{
		Stack:   viper.GetString("stack"),
		Service: viper.GetString("service"),
	}

	err := r.Validate()
	if err != nil {
		log.Info("Invalid rollback request run rollback --help for info")
	}

	rc, err := getRancherClient()

	if err != nil {
		return err
	}

	svc, err := getServiceByID(rc, r.Stack, r.Service)

	if err != nil {
		return err
	}

	_, err = rc.Service.ActionRollback(svc)

	return err
}

type rollbackRequest struct {
	Service string `valid:"required"`
	Stack   string `valid:"required"`
}
