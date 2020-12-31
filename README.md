# MLFlow Operator

MLFlow Operator generated via KubeBuilder to enable managing multiple MLFlow installs.

## Goals

The main goals of the operator are:

1. Ability to deploy MLFlow sites on top of Kubernetes
2. Provide best practices for application lifecycle
3. Facilitate proper devops (backups, monitoring and high-availability)
4. Provide basic Multi-Tenancy

## Workflow

Currently we have added an instances level folder that houses a `standard/tenant-1` kustomize deployment.

From this point on is just a straight copy of the folder and some slight customizations for additional instances:

* Minor tweaks to kustomization.yaml for that instance
* Add additional envvars to CI + `./params.sh`
* Add the redirect url to Azure AD

For programmatic access we will need to add custom header support to MLFLOW but for right now can use curl:

```sh
curl -v --cookie "oidc-cookie-XXXXX=XXXXX" https://mlflow-standard-tenant-1.covid.cloud.statcan.ca/api/2.0/preview/mlflow/experiments/list
```
