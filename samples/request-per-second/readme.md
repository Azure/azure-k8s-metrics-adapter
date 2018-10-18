# Requests per Second Custom Metric Scaling

This is an example on using custom metric from Application insights to scale a deployment.

- [Requests per Second Custom Metric Scaling](#requests-per-second-custom-metric-scaling)
    - [Walkthrough](#walkthrough)
    - [Configure Application Insights](#configure-application-insights)
        - [Create Application Insights](#create-application-insights)
        - [Get your instrumentation key](#get-your-instrumentation-key)
        - [Get your appid and api key](#get-your-appid-and-api-key)
            - [Using Azure Application Insights API Key](#using-azure-application-insights-api-key)
            - [Using Azure AD Pod Identity](#using-azure-ad-pod-identity)
    - [Deploy the app that will be scaled](#deploy-the-app-that-will-be-scaled)
- [there is probably a better way to get at that array](#there-is-probably-a-better-way-to-get-at-that-array)
    - [Scale on Requests per Second (RPS)](#scale-on-requests-per-second-rps)
        - [Deploy the Custom Metric Configuration](#deploy-the-custom-metric-configuration)
        - [Deploy the HPA](#deploy-the-hpa)
        - [Put it under load and scale by RPS](#put-it-under-load-and-scale-by-rps)
- [100000 requests at 100 RPS](#100000-requests-at-100-rps)
        - [Watch it scale](#watch-it-scale)
    - [Clean up](#clean-up)

## Walkthrough

Prerequisites:

- provisioned an [AKS Cluster](https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough)
- your `kubeconfig` points to your cluster.  
- [Metric Server deployed](https://github.com/kubernetes-incubator/metrics-server#deployment) to your cluster ([aks does not come with it deployed](https://github.com/Azure/AKS/issues/318)). Validate by running `kubectl get --raw "/apis/metrics.k8s.io/v1beta1/nodes" | jq .`

Get this repository and cd to this folder (on your GOPATH):

```
go get -u github.com/Azure/azure-k8s-metrics-adapter
cd $GOPATH/src/github.com/Azure/azure-k8s-metrics-adapter/samples/servicebus-queue/
```

## Configure Application Insights

### Create  Application Insights

First thing to do is create an [application insights instance](https://docs.microsoft.com/en-us/azure/application-insights/app-insights-nodejs-quick-start#enable-application-insights).

### Get your instrumentation key

After the application instance is created [get your instrumentation key](https://docs.microsoft.com/en-us/azure/application-insights/app-insights-nodejs-quick-start#configure-app-insights-sdk).

### Get your appid and api key

Get your [appid and key](https://dev.applicationinsights.io/documentation/Authorization/API-key-and-App-ID). Then deploy the adapter.

#### Using Azure Application Insights API Key

```bash
helm install --name sample-release ../../charts/azure-k8s-metrics-adapter --namespace custom-metrics \ 
    --set appInsights.appId=<your app id> \
    --set appInsights.key=<your api key> \
    --set azureAuthentication.createSecret=true
```

> note: if you plan to use the adapter with External Metrics you may need additional configuration.  See the [Service Bus Queue Example](../servicebus-queue/).

#### Using Azure AD Pod Identity

If you prefer to use [Azure AD Pod Identity](/README.md#using-azure-ad-pod-identity), then you don't need to specify an Application Insights API key:

```bash
kubectl create secret generic app-insights-api -n custom-metrics --from-literal=app-insights-app-id=<appid>
```

Deploy this modified [adapter-aad-pod-identity.yaml](https://gist.github.com/jcorioland/947af2c02acd3bc2b4d8438f1e36a6bd) file that includes Azure Identity and Azure Identity Binding:

```bash
kubectl apply -f https://gist.githubusercontent.com/jcorioland/947af2c02acd3bc2b4d8438f1e36a6bd/raw/9ff013c18d3a76a9c41d9fce40ad445b166013fa/adapter-aad-pod-identity.yaml
```

> Note: the managed user identity you are using should be authorized to read the Azure Application Insights resource through RBAC.

## Deploy the app that will be scaled

Create a secret with the application insights key that you retrieved in the earlier step:

```bash
kubectl create secret generic appinsightskey --from-literal=instrumentation-key=<your-key-here>

kubectl apply -f deploy/rps-deployment.yaml
```

> optional: build and push to your own copy of the example with `docker build -t metric-rps-example -f webapp/Dockerfile webapp`

Double check you can hit the endpoint:

```bash
# there is probably a better way to get at that array
export RPS_ENDPOINT="$(k get svc rps-sample  -o json | jq .status.loadBalancer.ingress | jq -r '.[0]'.ip)"

curl http://$RPS_ENDPOINT
```

## Scale on Requests per Second (RPS)

### Deploy the Custom Metric Configuration

```bash
kubectl apply -f deploy/custom-metric.yaml
```

> note: the CustomMetric configuration is deployed per namespace.

You can list of the configured custom metrics via:

```
kubectl get acm #shortcut for custommetric
```

### Deploy the HPA

Deploy the HPA:

```bash
kubectl apply -f deploy/hpa.yaml
```

> note: the `metrics.pods.metricName` defined on the HPA must match the `metadata.name` on the CustomMetric declaration, in this case `rps`


### Put it under load and scale by RPS

[Hey](https://github.com/rakyll/hey) is a simple way to create load on an api from  the command line.

```bash
go get -u github.com/rakyll/hey http://$RPS_ENDPOINT

# 100000 requests at 100 RPS
hey -n 10000 -q 10 -c 10
```

### Watch it scale

In a separate window you can watch the HPA to see the RPS go up and the pods scale:

```bash
kubectl get hpa rps-sample -w
NAME         REFERENCE               TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
rps-sample   Deployment/rps-sample   0/10      2         10        2          4d                                                            
rps-sample   Deployment/rps-sample   36/10     2         10        2         4d                                            
rps-sample   Deployment/rps-sample   36/10     2         10        4         4d
rps-sample   Deployment/rps-sample   36/10     2         10        4         4d                
rps-sample   Deployment/rps-sample   36/10     2         10        4         4d                                                          
rps-sample   Deployment/rps-sample   36/10     2         10        4         4d                                                             
rps-sample   Deployment/rps-sample   49/10     2         10        4         4d                                                      
rps-sample   Deployment/rps-sample   49/10     2         10        4         4d                                                             
rps-sample   Deployment/rps-sample   49/10     2         10        4         4d                                                     
rps-sample   Deployment/rps-sample   49/10     2         10        4         4d                                                             
rps-sample   Deployment/rps-sample   49/10     2         10        4         4d                                                     
rps-sample   Deployment/rps-sample   33/10     2         10        4         4d                                                             
rps-sample   Deployment/rps-sample   33/10     2         10        4         4d                                                          
rps-sample   Deployment/rps-sample   25/10     2         10        4         4d                                                             
rps-sample   Deployment/rps-sample   29/10     2         10        4         4d                                                          
rps-sample   Deployment/rps-sample   24/10     2         10        4         4d                                                             
rps-sample   Deployment/rps-sample   0/10      2         10        4         4d                                        
```

## Clean up

Once you are done with this experiment you can delete you Application Insights instance via portal.

Also remove resources created in cluster:

```bash
kubectl delete -f deploy/hpa.yaml
kubectl delete -f deploy/custom-metric.yaml
kubectl delete -f deploy/rps-deployment.yaml
helm delete --purge sample-release
```