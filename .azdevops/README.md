# Setting up Azure Dev Ops build pipelines

## Creating a pipeline
* In AzDO Pipelines, go to Build Pipelines > New > New Build Pipeline
* If you have have the preview feature "New YAML pipeline creation experience" enabled, make sure to click the link at the bottom of the first page to use the visual designer instead
* Set up your source with your repository and master branch
* Choose the 'Configuration as code' > YAML template to start out
* Set the YAML file path to `.azdevops/image-pipeline.yml`

You can get back to this screen (the visual designer) at any time from your saved pipeline by going to Edit > ••• > Pipeline settings.

## Setting up a variable group
* In AzDO Pipelines go to Library > + Variable Group
* To make minimal changes to the pipeline YAML file, name it 'Metrics Adapter' (otherwise, name it whatever you want and [change the YAML file](./image-pipeline.yml#L8))
* You'll need the following variables (secrets are denoted with \*\*\*\*\* as their example value):

| Name | Description | Example |
| --- | --- | --- |
| `modulePath` | Standard working directory (makes YAML files cleaner) | $(GOPATH)/src/github.com/Azure/azure-k8s-metrics-adapter \* |
| `GOBIN` | Golang bin directory for projects | $(GOPATH)/bin \* |
| `GOPATH` | Golang project directory | $(system.defaultWorkingDirectory)/go \* |
| `GOROOT` | Determines the version of Go used by AzDO | /usr/local/go1.11 \* |
| `HELM_VERSION` | Version of Helm to use | 2.12.0 |
| `MINIKUBE_VERSION` | Version of Minikube to use | 0.32.0 |
| `DOCKER_USER` | Docker username | user |
| `DOCKER_PASS` | Docker password (see note below) | \*\*\*\*\* |
| `REGISTRY` | Container registry address, leave empty if using DockerHub | myregistry.azurecr.io |
| `IMAGE` | Image name without the registry appended to the front | user/metrics-adapter-test |
| `SUBSCRIPTION_ID` | Azure subscription ID that the service bus namespace belongs to | <GUID\> |
| `SERVICEBUS_CONNECTION_STRING` | Service Bus namespace connection string | \*\*\*\*\* |
| `SERVICEBUS_NAMESPACE` | Service Bus namespace | my-namespace  |
| `SERVICEBUS_RESOURCE_GROUP` | Resource group containing the Service Bus namespace | my-resource-group |
| `SP_CLIENT_ID` | Service principal app ID | \*\*\*\*\* |
| `SP_CLIENT_SECRET` | Service principal password | \*\*\*\*\* |
| `SP_TENANT_ID` | Service principal tenant ID | \*\*\*\*\* |

\* Suggested variable value in AzDO

* You can then either check the option "Allow access to all pipelines" or, if you have multiple pipelines in your project and want to limit access to only the current pipeline, go to its visual designer > Variables tab > Variable groups > Link variable group and add your new variable group.

### A note about Docker passwords
Due to the way the password is entered on login in `image-pipeline.yml`, special characters might cause issues with logging in.

## Set up build triggers
This can be done in multiple ways depending on what you want and how you prefer to set it up (visual designer vs YAML). In order to keep most of the pipeline encoded in YAML, it's probably easier to set up triggers using the YAML syntax for [CI triggers](https://docs.microsoft.com/en-us/azure/devops/pipelines/yaml-schema?view=azure-devops&tabs=schema&viewFallbackFrom=azdevops#trigger) and [PR validation](https://docs.microsoft.com/en-us/azure/devops/pipelines/yaml-schema?view=azure-devops&tabs=schema&viewFallbackFrom=azdevops#pr-trigger). **To set up builds that trigger on PRs from forks, you need to use the visual designer** - go to Triggers > Pull request validation > Override the YAML pull request trigger from here > Build pull requests from forks of this repository. You'll need to turn on "Make secrets available to builds of forks." For an additionally layer of security, turn on "Only trigger builds for collaborators' pull request comments" - this will require a repository collaborator to comment `/azp run` on the PR before a build a triggered.

## Creating the necessary Service Bus Queues
This requires one queue for each version of kubernetes that will be tested - currently, 1.10 through 1.13. It's set up to default to queues named `externalq-10`, `externalq-11`, `externalq-12`, and `externalq-13`. These can be changed by editing the variables in [image-pipeline.yml](./image-pipeline.yml#L54). Do not use one queue for all combined tests - they need to individually test message counts.
