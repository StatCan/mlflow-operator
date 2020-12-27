#!/bin/bash

envsubst < policy/secret.tmpl > policy/secret.txt

export MLFLOW_OIDC_CLIENT="${MLFLOW_OIDC_CLIENT:=XXXXX-XXXXX-XXXXX}"
export MLFLOW_OIDC_SECRET="${MLFLOW_OIDC_SECRET:=XXXXX-XXXXX-XXXXX}"
export MLFLOW_OIDC_DISCOVERY="${MLFLOW_OIDC_DISCOVERY:=XXXXX-XXXXX-XXXXX}"
export MLFLOW_JWT_JWKS="${MLFLOW_JWT_JWKS:=XXXXX-XXXXX-XXXXX}"
export MLFLOW_ADMIN_POLICY_GROUP="${MLFLOW_ADMIN_POLICY_GROUP:=XXXXX-XXXXX-XXXXX}"
export MLFLOW_POLICY_GROUP="${MLFLOW_POLICY_GROUP:=XXXXX-XXXXX-XXXXX}"
export MLFLOW_DOMAIN_NAME="${MLFLOW_DOMAIN_NAME:=default.example.ca}"

for patch in policy/patch-oidc*; do
  yq w -i $patch '[0].value' ${MLFLOW_OIDC_CLIENT}
  yq w -i $patch '[2].value' ${MLFLOW_OIDC_DISCOVERY}
done

for patch in policy/patch-jwt*; do
  yq w -i $patch '[0].value' ${MLFLOW_JWT_JWKS}
done

for patch in policy/patch-policy*; do
  yq w -i $patch '[0].value' "https://${MLFLOW_DOMAIN_NAME}"
  yq w -i $patch '[1].value' ${MLFLOW_ADMIN_POLICY_GROUP}
  yq w -i $patch '[2].value' ${MLFLOW_POLICY_GROUP}
done
