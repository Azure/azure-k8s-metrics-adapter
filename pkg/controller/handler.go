package controller

import (
	listers "github.com/Azure/azure-k8s-metrics-adapter/pkg/client/listers/externalmetric/v1alpha1"
	"github.com/golang/glog"
)

type Handler struct {
	externalmetricLister listers.ExternalMetricLister
}

func NewHandler(externalmetricLister listers.ExternalMetricLister) Handler {
	return Handler{
		externalmetricLister: externalmetricLister,
	}
}

func (handler *Handler) Process(namespace string, name string) error {
	glog.Infof("processing item '%s' in namespace '%s'", namespace, name)
	return nil
}
