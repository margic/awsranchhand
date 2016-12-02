//go:generate gencheck -f=rancher.go

package cmd

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"
	"github.com/spf13/viper"
)

// GetRancherClient Returns a rancher client based on viper config
func getRancherClient() (rc *client.RancherClient, err error) {
	ro := rancherOpts{
		URL:       viper.GetString(rURL),
		AccessKey: viper.GetString(rKey),
		SecretKey: viper.GetString(rSecret),
	}

	err = ro.Validate()
	if err != nil {
		return nil, err
	}

	return client.NewRancherClient(&client.ClientOpts{
		Url:       ro.URL,
		AccessKey: ro.AccessKey,
		SecretKey: ro.SecretKey,
	})
}

type rancherOpts struct {
	URL       string `valid:"required"`
	AccessKey string `valid:"required"`
	SecretKey string `valid:"required"`
}

func doRancherPut(url string, body io.Reader) error {
	log.Debug(url)
	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, body)
	request.Header.Add("Content-Type", "application/json")
	request.SetBasicAuth(viper.GetString(rKey), viper.GetString(rSecret))

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"Status":  response.StatusCode,
		"Content": string(contents),
	}).Debug("Response")
	return nil
}

func lookupRancherHostID() (instanceID string, err error) {
	metaServiceURL := viper.GetString("rancherHostMetaServiceURL")
	log.WithField("rancherHostMetaServiceURL", metaServiceURL).Debug("Url")
	resp, err := http.Get(metaServiceURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func getServiceByID(rc *client.RancherClient, stackName string, serviceName string) (*client.Service, error) {
	log.WithField("serviceName", serviceName).Debug("Getting service ID")

	// environment ID or Application Stack name as referred to in rancher ui
	eID := ""

	if stackName != "" {
		// find the stack (environmnet id)
		eList, err := rc.Environment.List(&client.ListOpts{
			Filters: map[string]interface{}{
				"name": stackName,
			},
		})
		if err != nil {
			return nil, err
		}
		eData := eList.Data
		if len(eData) != 1 {
			return nil, NewServiceError("Application stack name not found or is not unique")
		}
		eID = eData[0].Id
		log.WithFields(log.Fields{
			"stackName": stackName,
			"stackId":   eID,
		}).Debug("Stack Found")
	}

	sFilter := map[string]interface{}{
		"name": serviceName,
	}
	if eID != "" {
		sFilter["environmentId"] = eID
	}

	sList, err := rc.Service.List(&client.ListOpts{
		Filters: sFilter,
	})

	if err != nil {
		return nil, err
	}
	sData := sList.Data
	if len(sData) == 0 {
		return nil, NewServiceError("Service not found")
	}
	if len(sData) > 1 {
		return nil, NewServiceError("Serice name not unique try running with stack and service name flags")
	}

	svcID := sData[0].Id
	log.WithFields(log.Fields{
		"serviceName": serviceName,
		"serviceId":   svcID,
	}).Debug("Service Found")
	// find the service by name
	svc, err := rc.Service.ById(svcID)

	return svc, err
}

// ServiceError common error that can occur when looking up services
type ServiceError struct {
	s string
}

func (e ServiceError) Error() string {
	return e.s
}

// NewServiceError create a new service error
func NewServiceError(msg string) ServiceError {
	return ServiceError{
		s: msg,
	}
}

// updates a laucnch config environment map with flags provided
func updateEnvironment(currentEnv map[string]interface{}, envFlags []string, clear bool) (newEnv map[string]interface{}) {
	if clear {
		currentEnv = make(map[string]interface{})
	}
	for _, e := range envFlags {
		log.Debug(e)
		key, value := splitEnv(e)
		currentEnv[key] = value
	}
	return currentEnv
}

func splitEnv(envString string) (key string, value string) {
	envs := strings.Split(envString, "=")
	return envs[0], envs[1]
}
