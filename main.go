/*
Uses base classes and Provider interfaces from https://github.com/kubernetes-incubator/custom-metrics-apiserver to build
a metric server for Azure based services.
*/

package main

import (
	"flag"
	"os"
	"runtime"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/az-metric-client"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/provider"
	"github.com/golang/glog"
	basecmd "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/cmd"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/util/logs"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	cmd := &basecmd.AdapterBase{}
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	cmd.Flags().Parse(os.Args)

	client, err := cmd.DynamicClient()
	if err != nil {
		glog.Fatalf("unable to construct dynamic client: %v", err)
	}

	mapper, err := cmd.RESTMapper()
	if err != nil {
		glog.Fatalf("unable to construct discovery REST mapper: %v", err)
	}

	azureProvider := provider.NewAzureProvider(client, mapper, azureMetricClient.NewAzureMetricClient())
	cmd.WithCustomMetrics(azureProvider)
	cmd.WithExternalMetrics(azureProvider)

	if err := cmd.Run(wait.NeverStop); err != nil {
		glog.Fatalf("Unable to run Azure metrics adapter: %v", err)
	}
}
