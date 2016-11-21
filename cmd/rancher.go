//go:generate gencheck -f=rancher.go

package cmd

import (
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"
	"github.com/spf13/viper"
)

// GetRancherClient Returns a rancher client based on viper config
func GetRancherClient() (rc *client.RancherClient, err error) {
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
