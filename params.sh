#!/bin/bash

# Operator image version
export MLFLOW_OPERATOR_IMAGE="${MLFLOW_OPERATOR_IMAGE:=controller}"
export MLFLOW_OPERATOR_IMAGE_VERSION="${MLFLOW_OPERATOR_IMAGE_VERSION:=latest}"
export CONFIG_FILE="${CONFIG_FILE:-config/manager/kustomization.yaml}"
yq w -i ${CONFIG_FILE} 'images[0].newName' ${MLFLOW_OPERATOR_IMAGE}
yq w -i ${CONFIG_FILE} 'images[0].newTag' ${MLFLOW_OPERATOR_IMAGE_VERSION}
