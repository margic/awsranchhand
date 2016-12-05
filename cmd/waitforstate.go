//go:generate gencheck -f=waitforstate.go

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
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// waitforstateCmd represents the waitforstate command
var waitforstateCmd = &cobra.Command{
	Use:   "waitforstate",
	Short: "Waits for a service to be in a specified state.",
	Long: `During a deploy or upgrade services transition between states.
	Operations are not possible if the service is not in the correct state.
	This command will wait for a service to be in a specified state and can
	be used when automating deploy or upgrade steps.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("waitforstate called")

		err := waitforstate()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(waitforstateCmd)
	waitforstateCmd.Flags().DurationP("timeout", "t", 60*time.Second, "How long to wait for service state, default is 60s")
	waitforstateCmd.Flags().StringP("state", "a", "", "State to wait for e.g. example ugraded, active")
	viper.BindPFlags(waitforstateCmd.Flags())
}

func waitforstate() error {
	r := waitrequest{
		Stack:   viper.GetString("stack"),
		Service: viper.GetString("service"),
		Timeout: viper.GetDuration("timeout"),
		State:   viper.GetString("state"),
	}

	err := r.Validate()
	if err != nil {
		log.Error("Invalid options run waitforstate --help for options")
		return err
	}

	rc, err := getRancherClient()

	if err != nil {
		return err
	}

	// channel state found
	cs := make(chan bool)

	go func(cs chan bool) {
		ticker := time.NewTicker(5 * time.Second)
		log.WithFields(log.Fields{
			"stack":   r.Stack,
			"service": r.Service,
			"state":   r.State,
		}).Debug("Checking service state")

		for range ticker.C {
			// this is crappy doesn't call until after first tick of the timer
			svc, err := getServiceByID(rc, r.Stack, r.Service)
			if err != nil {
				log.Debug(err)
				_, ok := err.(ServiceError)
				if ok {
					cs <- false
				}
			} else {
				foundState := svc.State
				log.WithField("currentState", foundState).Debug("Service State")
				if foundState == r.State {
					cs <- true
				}
			}
		}
	}(cs)

	select {
	case found := <-cs:
		if found {
			log.Debug("FoundState")
		} else {
			log.Debug("NotFound")
			return errors.New("State not found")
		}
	case <-time.After(r.Timeout):
		log.WithFields(log.Fields{
			"service": r.Service,
			"state":   r.State,
		}).Debug("State not found")
		return errors.New("State not found")
	}
	return nil
}

type waitrequest struct {
	Stack   string        `valid:"required"`
	Service string        `valid:"required"`
	State   string        `valid:"required"`
	Timeout time.Duration `valid:"required"`
}
