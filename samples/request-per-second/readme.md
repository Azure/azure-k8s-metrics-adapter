# Requests per Second Custom Metric Scaling
This is an example on using custom metric from Application insights to scale a deployment.

> note: this is currently a work in progress

- [Requests per Second Custom Metric Scaling](#requests-per-second-custom-metric-scaling)
    - [Create Application Insights](#create-application-insights)
    - [Get your instrumentation key](#get-your-instrumentation-key)
    - [Get your appid and api key](#get-your-appid-and-api-key)
    - [Build the nodejs application](#build-the-nodejs-application)
    - [Deploy your app](#deploy-your-app)
    - [Deploy the HPA](#deploy-the-hpa)
    - [Put it under load](#put-it-under-load)
    - [The Raw query](#the-raw-query)

## Create  Application Insights 

https://docs.microsoft.com/en-us/azure/application-insights/app-insights-nodejs-quick-start#enable-application-insights

## Get your instrumentation key

https://docs.microsoft.com/en-us/azure/application-insights/app-insights-nodejs-quick-start#configure-app-insights-sdk

## Get your appid and api key
Get your key:

https://dev.applicationinsights.io/documentation/Authorization/API-key-and-App-ID

Add the following to your adapter deployment manifest:

```
- name: APP_INSIGHTS_APP_ID
valueFrom:
    secretKeyRef:
    name: app-insights-api
    key: app-insights-app-id
- name: APP_INSIGHTS_KEY
valueFrom:
    secretKeyRef:
    name: app-insights-api
    key: app-insights-key
```

Create a secret for the adapter to use:

```
kubectl create secret generic app-insights-api -n custom-metrics --from-literal=app-insights-app-id=<appid> --from-literal=app-insights-key=<key> 
```

## Build the nodejs application

```
docker build -t metric-rps-example -f webapp/Dockerfile webapp
```

> optional: push to your own repository

## Deploy your app

Create a secret with the application insights key:

```
kubectl create secret generic appinsightskey --from-literal=instrumentation-key=<your-key-here>

kubectl apply -f deploy/rps-deployment.yaml
```

Double check you can hit the endpoint:

```bash
# there is probably a better way to get at that array
export RPS_ENDPOINT="$(k get svc rps-sample  -o json | jq .status.loadBalancer.ingress | jq -r '.[0]'.ip)"

curl http://$RPS_ENDPOINT
```

## Deploy the HPA

Deploy the HPA:

```bash
kubectl apply -f deploy/hpa.yaml
```


## Put it under load

```
go get -u github.com/rakyll/hey

# 10000 requests 5 RPS
hey -n 10000 -q 5 -c 5 
```

##  The Raw query
Can use for debugging: https://dev.applicationinsights.io/apiexplorer/metrics

GET /v1/apps/<yourkey>/metrics/performanceCounters/requestsPerSecond?timespan=PT5M&interval=PT1M HTTP/1.1
