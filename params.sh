#!/bin/bash

# MLFlow
export MLFLOW_OPERATOR_IMAGE="${MLFLOW_OPERATOR_IMAGE:=controller}"
export MLFLOW_OPERATOR_IMAGE_VERSION="${MLFLOW_OPERATOR_IMAGE_VERSION:=latest}"
export CONFIG_FILE="${CONFIG_FILE:-operator/manager/kustomization.yaml}"
export MLFLOW_IMAGE="${MLFLOW_IMAGE:=mlflow}"
export MLFLOW_IMAGEPULLSECRET="${MLFLOW_IMAGEPULLSECRET:=registry-connection}"
export MLFLOW_OIDC_CLIENT="${MLFLOW_OIDC_CLIENT:=XXXXX-XXXXX-XXXXX}"
export MLFLOW_OIDC_SECRET="${MLFLOW_OIDC_SECRET:=XXXXX-XXXXX-XXXXX}"
export MLFLOW_OIDC_DISCOVERY="${MLFLOW_OIDC_DISCOVERY:=XXXXX-XXXXX-XXXXX}"
export MLFLOW_ADMIN_POLICY_GROUP="${MLFLOW_ADMIN_POLICY_GROUP:=XXXXX-XXXXX-XXXXX}"

export STANDARD_TENANT_1_DOMAIN_NAME="${STANDARD_TENANT_1_DOMAIN_NAME:=default.example.ca}"
export STANDARD_TENANT_1_PREFIX="${STANDARD_TENANT_1_PREFIX:=XXXXX-XXXXX-XXXXX}"
export STANDARD_TENANT_1_POLICY_GROUP="${STANDARD_TENANT_1_POLICY_GROUP:=XXXXX-XXXXX-XXXXX}"

yq w -i ${CONFIG_FILE} 'images[0].newName' ${MLFLOW_OPERATOR_IMAGE}
yq w -i ${CONFIG_FILE} 'images[0].newTag' ${MLFLOW_OPERATOR_IMAGE_VERSION}

envsubst < instances/standard/tenant-1/secret.tmpl > instances/standard/tenant-1/secret.txt

for patch in instances/*/tenant*/patch-ing*; do
  yq w -i $patch '[0].value' ${STANDARD_TENANT_1_DOMAIN_NAME}
  yq w -i $patch '[1].value' ${STANDARD_TENANT_1_PREFIX}
done

for patch in instances/*/tenant*/patch-oidc*; do
  yq w -i $patch '[0].value' ${MLFLOW_OIDC_CLIENT}
  yq w -i $patch '[1].value' ${STANDARD_TENANT_1_PREFIX}
  yq w -i $patch '[2].value' ${MLFLOW_OIDC_DISCOVERY}
done

for patch in instances/*/tenant*/patch-policy*; do
  yq w -i $patch '[0].value' "https://${STANDARD_TENANT_1_DOMAIN_NAME}"
  yq w -i $patch '[1].value' ${MLFLOW_ADMIN_POLICY_GROUP}
  yq w -i $patch '[2].value' ${STANDARD_TENANT_1_POLICY_GROUP}
  yq w -i $patch '[3].value' ${STANDARD_TENANT_1_PREFIX}
  yq w -i $patch '[4].value' ${STANDARD_TENANT_1_PREFIX}
done

for patch in instances/*/tenant*/patch-trackingserver*; do
  yq w -i $patch '[0].value' ${MLFLOW_IMAGE}
  yq w -i $patch '[1].value' ${MLFLOW_IMAGEPULLSECRET}
done
