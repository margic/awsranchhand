package cmd

import (
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

func lookupInstanceID() (instanceID string, err error) {
	metaServiceURL := viper.GetString("metaServiceURL")
	log.WithField("instanceMetaUrl", metaServiceURL).Debug("Url")
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
