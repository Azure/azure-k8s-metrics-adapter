package controller

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	informers "github.com/Azure/azure-k8s-metrics-adapter/pkg/client/informers/externalversions/metrics/v1alpha1"
)

// Controller will do the work of syncing the external metrics the metric adapter knows about.
type Controller struct {
	metricQueue          workqueue.RateLimitingInterface
	externalMetricSynced cache.InformerSynced
	customMetricSynced   cache.InformerSynced
	enqueuer             func(obj interface{})
	metricHandler        ContollerHandler
}

// NewController returns a new controller for handling external and custom metric types
func NewController(externalMetricInformer informers.ExternalMetricInformer, customMetricInformer informers.CustomMetricInformer, metricHandler ContollerHandler) *Controller {
	controller := &Controller{
		externalMetricSynced: externalMetricInformer.Informer().HasSynced,
		metricQueue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "metrics"),
		metricHandler:        metricHandler,
		customMetricSynced:   customMetricInformer.Informer().HasSynced,
	}

	// wire up enque step.  This provides a hook for testing enqueue step
	controller.enqueuer = controller.enqueueExternalMetric

	glog.Info("Setting up external metric event handlers")
	externalMetricInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueuer,
		UpdateFunc: func(old, new interface{}) {
			// Watches and Informers will “sync”.
			// Periodically, they will deliver every matching object in the cluster to your Update method.
			// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
			controller.enqueuer(new)
		},
		DeleteFunc: controller.enqueuer,
	})

	glog.Info("Setting up custom metric event handlers")
	customMetricInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueuer,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueuer(new)
		},
		DeleteFunc: controller.enqueuer,
	})

	return controller
}

// Run is the main path of execution for the controller loop
func (c *Controller) Run(numberOfWorkers int, interval time.Duration, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.metricQueue.ShutDown()

	glog.V(2).Info("initializing controller")

	// do the initial synchronization (one time) to populate resources
	if !cache.WaitForCacheSync(stopCh, c.externalMetricSynced, c.customMetricSynced) {
		runtime.HandleError(fmt.Errorf("Error syncing controller cache"))
		return
	}

	glog.V(2).Infof("starting %d workers with %d interval", numberOfWorkers, interval)
	for i := 0; i < numberOfWorkers; i++ {
		go wait.Until(c.runWorker, interval, stopCh)
	}

	<-stopCh
	glog.Info("Shutting down workers")
	return
}

func (c *Controller) runWorker() {
	glog.V(2).Info("Worker starting")

	for c.processNextItem() {
		glog.V(2).Info("processing next item")
	}

	glog.V(2).Info("worker completed")
}

func (c *Controller) processNextItem() bool {
	glog.V(2).Info("processing item")

	key, quit := c.metricQueue.Get()
	if quit {
		glog.V(2).Info("recieved quit signal")
		return false
	}

	defer c.metricQueue.Done(key)

	var namespaceNameKey string
	var ok bool
	if namespaceNameKey, ok = key.(string); !ok {
		// not valid key do not put back on queue
		c.metricQueue.Forget(key)
		runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", key))
		return true
	}

	err := c.metricHandler.Process(namespaceNameKey)
	if err != nil {
		retrys := c.metricQueue.NumRequeues(key)
		if retrys < 5 {
			glog.Errorf("Transient error with %d retrys for key %s: %s", retrys, key, err)
			c.metricQueue.AddRateLimited(key)
			return true
		}

		// something was wrong with the item on queue
		glog.Errorf("Max retries hit for key %s: %s", key, err)
		c.metricQueue.Forget(key)
		utilruntime.HandleError(err)
		return true
	}

	//if here success for get item
	glog.V(2).Infof("succesfully proccessed item '%s'", namespaceNameKey)
	c.metricQueue.Forget(key)
	return true
}

func (c *Controller) enqueueExternalMetric(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}

	glog.V(2).Infof("adding item to queue for '%s'", key)
	c.metricQueue.AddRateLimited(key)
}
