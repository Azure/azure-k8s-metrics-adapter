/*
	Modified from https://github.com/Microsoft/azureimds/blob/master/imdssample.go
	under Apache License Version 2.0. See diff for changes.
*/

package instancemetadata

import (
	"io/ioutil"
	"net/http"

	"k8s.io/klog"
)

type AzureConfig struct {
	SubscriptionID string
}

func GetAzureConfig() (AzureConfig, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "http://169.254.169.254/metadata/instance/compute/subscriptionId", nil)
	req.Header.Add("Metadata", "True")

	q := req.URL.Query()
	q.Add("format", "text")
	q.Add("api-version", "2017-12-01")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		klog.Errorf("unable to get metadata for azure vm: %v", err)
		return AzureConfig{}, err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	subID := string(respBody[:])

	klog.V(2).Infoln("connected to sub:", subID)

	config := AzureConfig{
		SubscriptionID: subID,
	}
	return config, nil
}
