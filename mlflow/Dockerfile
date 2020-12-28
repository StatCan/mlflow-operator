FROM python:3.8.2-slim

COPY entrypoint.sh /

RUN apt-get update && \
    apt-get install -y git && \
    rm -rf /var/lib/apt/lists/*

RUN pip --no-cache-dir install --quiet --no-use-pep517 \
      'mlflow==1.13.0' \
      'azure-storage-blob==12.6.0' \
      'azure-storage-file-share==12.3.0' \
      'azure-storage-file-datalake==12.2.0' \
      'azure-storage-queue==12.1.4' \
      'msrestazure~=0.6.4' \
      'psycopg2-binary==2.8.6'

# REQUIRED FOR AZURE:
# ENV AZURE_STORAGE_CONNECTION_STRING "XXXXX"
# ENV MLFLOW_SERVER_FILE_STORE "sqlite:///mnt/mlruns/mlflow.db"
# ENV MLFLOW_SERVER_DEFAULT_ARTIFACT_ROOT "wasbs://XXXXX@XXXXX.blob.core.windows.net/mlartifacts"

ENV MLFLOW_SERVER_HOST 0.0.0.0
ENV MLFLOW_SERVER_PORT 5000
ENV MLFLOW_SERVER_WORKERS 4

CMD ["/entrypoint.sh"]
