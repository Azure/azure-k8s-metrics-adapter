# Setting up Azure Dev Ops build pipelines

## Creating a pipeline
* In ADO Pipelines, go to Build Pipelines > New > New Build Pipeline
* In the first step, click the link at the bottom to use the visual designer
* Set up your source with the repository/branch the pipeline YAML files are in
* Choose the 'Configuration as code' > YAML template to start out
* Set the YAML file path
You can get back to this screen (the visual designer) at any time from your saved pipeline by going to Edit > ••• > Pipeline settings.

Do this for each of `image-pipeline.yml` and `deploy-pipeline.yml`. For at least `image-pipeline`, you'll want to name the ADO pipeline accordingly or change how it's referred to in `deploy-pipeline.yml`.

## Setting up a variable group
* In ADO Pipelines go to Library > + Variable Group
* To make minimal changes to the pipeline YAML files, name it 'Metrics Adapter' (otherwise, name it whatever you want and change the YAML files)
* You'll need the following variables (secrets are denoted with \*\*\*\*\* as their example value):

| Name | Description | Example |
| --- | --- | --- |
| `modulePath` | Standard working directory (makes YAML files cleaner) | $(GOPATH)/src/github.com/Azure/azure-k8s-metrics-adapter \* |
| `GOBIN` | Golang bin directory for projects | $(GOPATH)/bin \* |
| `GOPATH` | Golang project directory | $(system.defaultWorkingDirectory)/go \* |
| `GOROOT` | Determines the version of Go used by ADO | /usr/local/go1.11 \* |
| `HELM_VERSION` | Version of Helm to use | 2.12.0 |
| `MINIKUBE_VERSION` | Version of Minikube to use | 0.32.0 |
| `DOCKER_USER` | Docker username | user |
| `DOCKER_PASS` | Docker password (see note below) | \*\*\*\*\* |
| `REGISTRY` | Container registry address (use example if using DockerHub) | https://index.docker.io/v1/ |
| `FULL_IMAGE` | Full name of the image, excluding the tag | user/metrics-adapter-test |
| `SUBSCRIPTION_ID` | Azure subscription ID that the service bus namespace belongs to | <GUID\> |
| `SERVICEBUS_CONNECTION_STRING` | Service Bus namespace connection string | \*\*\*\*\* |
| `SERVICEBUS_NAMESPACE` | Service Bus namespace | my-namespace  |
| `SERVICEBUS_RESOURCE_GROUP` | Resource group containing the Service Bus namespace | my-resource-group |
| `SP_CLIENT_ID` | Service principal app ID | \*\*\*\*\* |
| `SP_CLIENT_SECRET` | Service principal password | \*\*\*\*\* |
| `SP_TENANT_ID` | Service principal tenant ID | \*\*\*\*\* |

\* Suggested variable value in ADO

* In each pipelines' visual designer, go to the Variables tab > Variable groups > Link variable group and add the new variable group

## Set up build triggers
This can be done in multiple ways depending on what you want and how you prefer to set it up (visual designer vs YAML). In order to keep most of the pipeline encoded in YAML, it's probably easier to set up triggers using the YAML syntax for [CI triggers](https://docs.microsoft.com/en-us/azure/devops/pipelines/yaml-schema?view=azure-devops&tabs=schema&viewFallbackFrom=azdevops#trigger) and [PR validation](https://docs.microsoft.com/en-us/azure/devops/pipelines/yaml-schema?view=azure-devops&tabs=schema&viewFallbackFrom=azdevops#pr-trigger). (Also, there are triggers specific to my branches set up in `image-pipeline.yml`, which you'll want to change.) Build completion triggers are not yet supported in YAML, so in the visual designer of `deploy-pipeline`, go to Triggers > Build completion > + Add and set a triggering build.

### A note about Docker passwords
Due to the way the password is entered on login in `image-pipeline.yml`, special characters might cause issues with logging in.
