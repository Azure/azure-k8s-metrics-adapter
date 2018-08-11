# Samples for authentication with Azure
There are several was to authenticate the adapter with Azure.  

## Using Azure AD Application ID and Secret
To create secrets for use with [AD Service Principal](adapter-with-ad-sp.yaml):

```
az ad sp create-for-rbac -n azure-metric-adapter 
az role assignment create --role "Monitoring Reader" --assignee-object-id <objectid> --resource-group sb-external-example

# use output from create-for-rbac to create secret
kubectl create secret generic adapter-service-principal -n custom-metrics --from-literal=azure-tenant-id=<tenantid> --from-literal=azure-client-id=<azure-client-id>  --from-literal=azure-client-secret=<azure-client-secret>
```

## Application Insights Permissions
https://docs.microsoft.com/en-us/azure/application-insights/app-insights-resources-roles-access-control

https://stackoverflow.com/questions/42978366/authenticate-on-application-insights-rest-api-with-aad?rq=1