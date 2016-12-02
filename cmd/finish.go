//go:generate gencheck -f=finish.go

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
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// finishCmd represents the finish command
var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "Complete a rancher service upgrade",
	Long: `Complete an upgrade and remove the old containers.
		This should be run after an upgrade has been verified.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("finish called")
		err := finish()
		if err != nil {
			log.Fatal("Finish Upgrade failed", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(finishCmd)
}

func finish() error {
	r := finishRequest{
		Stack:   viper.GetString("stack"),
		Service: viper.GetString("service"),
	}

	err := r.Validate()
	if err != nil {
		log.Info("Invalid finish request run finish --help")
		return err
	}

	if r.Service == "" {
		return errors.New("Upgrade command missing service name type ranchdeploy upgrade --help")
	}

	rc, err := getRancherClient()

	if err != nil {
		return err
	}

	svc, err := getServiceByID(rc, r.Stack, r.Service)

	if err != nil {
		return err
	}

	_, err = rc.Service.ActionFinishupgrade(svc)

	return err
}

type finishRequest struct {
	Stack   string `valid:"required"`
	Service string `valid:"required"`
}
