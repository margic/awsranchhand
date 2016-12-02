//go:generate gencheck -f=upgrade.go

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
	"github.com/rancher/go-rancher/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/Sirupsen/logrus"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade a service",
	Long:  `Upgrade an existing service by name`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("upgrade called")
		err := upgrade()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	log.Info("Command Init")
	RootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().StringP("image", "i", "", "Optional - replace the docker image for upgraded service")
	upgradeCmd.Flags().BoolP("nostartfirst", "n", false, "Start the new service before stopping the old one")
	upgradeCmd.Flags().BoolP("clear", "c", false, "Optional - Clear all existing environment variable for service")
	upgradeCmd.Flags().StringSliceP("env", "e", nil, "Optional - Add or update the current configured environment variables for the service")
	viper.BindPFlags(upgradeCmd.Flags())
}

type upgradeRequest struct {
	Stack        string `valid:"required"`
	Service      string `valid:"required"`
	Image        string
	NoStartFirst bool
	ClearEnv     bool
	Env          []string
}

func upgrade() error {

	r := upgradeRequest{
		Stack:        viper.GetString("stack"),
		Service:      viper.GetString("service"),
		Image:        viper.GetString("image"),
		NoStartFirst: viper.GetBool("nostartfirst"),
		ClearEnv:     viper.GetBool("clear"),
		Env:          viper.GetStringSlice("env"),
	}

	log.WithFields(log.Fields{
		"Stack":        r.Stack,
		"Service":      r.Service,
		"Image":        r.Image,
		"NoStartFirst": r.NoStartFirst,
		"ClearEnv":     r.ClearEnv,
		"Env":          r.Env,
	}).Debug("Upgrade Request")

	err := r.Validate()
	if err != nil {
		log.Info("Invalid update request run upgrade --help for more info")
		return err
	}

	rc, err := getRancherClient()

	if err != nil {
		return err
	}

	svc, err := getServiceByID(rc, r.Stack, r.Service)

	if err != nil {
		return err
	}

	lc := svc.LaunchConfig

	if r.Image != "" {
		// apply new image to lc
		oldImage := lc.ImageUuid
		newImage := "docker:" + r.Image
		lc.ImageUuid = newImage
		log.WithFields(log.Fields{
			"oldImage": oldImage,
			"newImage": newImage,
		}).Info("Docker Image Updated")
	}

	if len(r.Env) > 0 || r.ClearEnv {
		// have to update environment
		lcEnv := lc.Environment
		lcEnv = updateEnvironment(lcEnv, r.Env, r.ClearEnv)
		lc.Environment = lcEnv
	}

	su := client.ServiceUpgrade{
		InServiceStrategy: &client.InServiceUpgradeStrategy{
			LaunchConfig: lc,
			StartFirst:   !r.NoStartFirst,
		},
	}

	svc, err = rc.Service.ActionUpgrade(svc, &su)
	return err
}
