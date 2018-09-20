/*
Uses base classes and Provider interfaces from https://github.com/kubernetes-incubator/custom-metrics-apiserver to build
a metric server for Azure based services.
*/

package main

import (
	"flag"
	"os"
	"runtime"
	"time"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/az-metric-client"
	clientset "github.com/Azure/azure-k8s-metrics-adapter/pkg/client/clientset/versioned"
	informers "github.com/Azure/azure-k8s-metrics-adapter/pkg/client/informers/externalversions"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/controller"
	azureprovider "github.com/Azure/azure-k8s-metrics-adapter/pkg/provider"
	"github.com/golang/glog"
	basecmd "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/cmd"
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

	stopCh := make(chan struct{})
	defer close(stopCh)

	metriccache := metriccache.NewMetricCache()

	// start and run contoller components
	controller, adapterInformerFactory := newController(cmd, metriccache)
	go adapterInformerFactory.Start(stopCh)
	go controller.Run(2, time.Second, stopCh)

	//setup and run metric server
	setupAzureProvider(cmd, metriccache)
	if err := cmd.Run(stopCh); err != nil {
		glog.Fatalf("Unable to run Azure metrics adapter: %v", err)
	}
}

func setupAzureProvider(cmd *basecmd.AdapterBase, metricsCache *metriccache.MetricCache) {
	client, err := cmd.DynamicClient()
	if err != nil {
		glog.Fatalf("unable to construct dynamic client: %v", err)
	}

	mapper, err := cmd.RESTMapper()
	if err != nil {
		glog.Fatalf("unable to construct discovery REST mapper: %v", err)
	}

	azureProvider := azureprovider.NewAzureProvider(client, mapper, azureMetricClient.NewAzureMetricClient(metricsCache))
	cmd.WithCustomMetrics(azureProvider)
	cmd.WithExternalMetrics(azureProvider)
}

func newController(cmd *basecmd.AdapterBase, metricsCache *metriccache.MetricCache) (*controller.Controller, informers.SharedInformerFactory) {
	clientConfig, err := cmd.ClientConfig()
	if err != nil {
		glog.Fatalf("unable to construct client config: %s", err)
	}
	adapterClientSet, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		glog.Fatalf("unable to construct lister client to initialize provider: %v", err)
	}

	adapterInformerFactory := informers.NewSharedInformerFactory(adapterClientSet, time.Second*30)

	handler := controller.NewHandler(adapterInformerFactory.Azure().V1alpha1().ExternalMetrics().Lister(), metricsCache)
	controller := controller.NewController(adapterInformerFactory.Azure().V1alpha1().ExternalMetrics(), handler)

	return controller, adapterInformerFactory
}
