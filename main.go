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

	"github.com/Azure/azure-k8s-metrics-adapter/pkg/aim"
	"github.com/Azure/azure-k8s-metrics-adapter/pkg/metriccache"
	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/go-autorest/autorest/azure/auth"

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

	defaultSubscriptionID := getDefaultSubscriptionID()
	monitorClient := insights.NewMetricsClient(defaultSubscriptionID)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		monitorClient.Authorizer = authorizer
	}

	azureProvider := azureprovider.NewAzureProvider(client, mapper, azureMetricClient.NewAzureMetricClient(defaultSubscriptionID, metricsCache, monitorClient))
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

func getDefaultSubscriptionID() string {
	// if the user explicitly sets we should use that
	subscriptionID := os.Getenv("SUBSCRIPTION_ID")
	if subscriptionID == "" {
		//fallback to trying azure instance meta data
		azureConfig, err := aim.GetAzureConfig()
		if err != nil {
			glog.Errorf("Unable to get azure config from MSI: %v", err)
		}

		subscriptionID = azureConfig.SubscriptionID
	}

	if subscriptionID == "" {
		glog.V(0).Info("Default Azure Subscription is not set.  You must provide subscription id via HPA lables, set an environment variable, or enable MSI.  See docs for more details")
	}

	return subscriptionID
}
