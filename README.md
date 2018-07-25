# azure-k8-metrics-adapter

An implementation of the Kubernetes [Custom Metrics API and External Metrics API](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis) for Azure Services. 

## Deploy

```
kubectl apply -f https://raw.githubusercontent.com/jsturtevant/azure-k8-metrics-adapter/master/deploy/adapter.yaml
```

After deployment you can query the api:

```bash
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1" | jq .
kubectl get --raw "/apis/external.metrics.k8s.io/v1beta1" | jq .
```

## Azure Setup

Enable [Managed Service Identity](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/tutorial-linux-vm-access-arm) on each of your AKS vms and give access to the resource the MSI access for each vm:

```bash
export RG=<aks resource group> 
export CLUSTER=<aks cluster name> 
export ACCESS_RG=<rg to give read access to>

NODE_RG="$(az aks show -n $CLUSTER -g $RG | jq -r .nodeResourceGroup)"
az vm list -g $NODE_RG
VMS="$(az vm list -g $NODE_RG | jq -r '.[] | select(.tags.creationSource | . and contains("aks")) | .name')"

while read -r vm; do
    echo "updating vm $vm..."
    msi="$(az vm identity assign -g $NODE_RG -n $vm | jq -r .systemAssignedIdentity)"

    echo "adding access with $msi..."
    az role assignment create --role Reader --assignee-object-id $msi --resource-group $ACCESS_RG
done <<< "$VMS"
```

## Development

### Get the source

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

