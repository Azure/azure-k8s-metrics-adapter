# azure-k8-metrics-adapter

An implementation of the Kubernetes [Custom Metrics API and External Metrics API](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis) for Azure Services.  

# Build and deploy

```bash
go get github.com/jsturtevant/azure-k8-metrics-adapter
cd $GOPATH/github.com/jsturtevant/azure-k8-metrics-adapter

export REGISTRY=<your-registry>
make container-build
make container-push
kubectl apply -f deploy/manifests
```

After deployment you can query the api:

```bash
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1" | jq .
```