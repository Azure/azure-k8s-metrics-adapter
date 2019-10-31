module github.com/Azure/azure-k8s-metrics-adapter

go 1.12

require (
	github.com/Azure/azure-sdk-for-go v30.1.0+incompatible
	github.com/Azure/azure-service-bus-go v0.9.1
	github.com/Azure/go-autorest v12.0.0+incompatible
	github.com/dimchansky/utfbom v1.1.0 // indirect
	github.com/emicklei/go-restful v2.2.1+incompatible // indirect
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170208215640-dcef7f557305 // indirect
	github.com/evanphx/json-patch v4.2.0+incompatible // indirect
	github.com/kubernetes-incubator/custom-metrics-apiserver v0.0.0-20190918110929-3d9be26a50eb
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	k8s.io/api v0.0.0-20190817021128-e14a4b1f5f84
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/apiserver v0.0.0-20190817022445-fd6150da8f40 // indirect
	k8s.io/client-go v0.0.0-20190817021527-637fc595d17a
	k8s.io/code-generator v0.0.0-20190612205613-18da4a14b22b
	k8s.io/component-base v0.0.0-20190817022002-dd0e01d5790f
	k8s.io/klog v0.3.1
	k8s.io/metrics v0.0.0-20190817023635-63ee757b2e8b

)

replace github.com/Sirupsen/logrus v1.4.2 => github.com/sirupsen/logrus v1.4.2
