#!/bin/bash

export MLFLOW_SERVER_FILE_STORE="${MLFLOW_SERVER_FILE_STORE:-./mlruns}"
export MLFLOW_SERVER_DEFAULT_ARTIFACT_ROOT="${MLFLOW_SERVER_DEFAULT_ARTIFACT_ROOT:-./mlruns}"

mlflow db upgrade $MLFLOW_SERVER_FILE_STORE

mlflow server \
    --backend-store-uri "$MLFLOW_SERVER_FILE_STORE" \
    --default-artifact-root "$MLFLOW_SERVER_DEFAULT_ARTIFACT_ROOT" \
    --host "$MLFLOW_SERVER_HOST" \
    --port "$MLFLOW_SERVER_PORT" \
    --workers "$MLFLOW_SERVER_WORKERS"
