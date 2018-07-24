# azure-k8-metrics-adapter

An implementation of the Kubernetes [Custom Metrics API and External Metrics API](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis) for Azure Services. 

# Deploy

```
kubectl apply -f https://raw.githubusercontent.com/jsturtevant/azure-k8-metrics-adapter/master/deploy/adapter.yaml
```

After deployment you can query the api:

```bash
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1" | jq .
kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1" | jq .
```

# Development

## Get the source

```bash
go get github.com/jsturtevant/azure-k8-metrics-adapter
cd $GOPATH/github.com/jsturtevant/azure-k8-metrics-adapter
```

### Use Skaffold
Before you run the command below be sure to:

- Download [skaffold](https://github.com/GoogleContainerTools/skaffold#installation) 
- Log in to your container registry: `docker login`
- Have your K8 context set to the cluster you want to deploy to: `kubectl config use-context`

```bash
make dev
```

