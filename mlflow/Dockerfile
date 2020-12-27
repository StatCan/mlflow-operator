FROM python:3.8.2-slim

COPY entrypoint.sh /

RUN apt-get update && \
    apt-get install -y git && \
    rm -rf /var/lib/apt/lists/*

RUN pip --no-cache-dir install --quiet --no-use-pep517 \
      'mlflow==1.9.1' \
      'azure-storage-blob==12.3.2' \
      'azure-storage-file-share==12.1.2' \
      'azure-storage-file-datalake==12.0.2' \
      'azure-storage-queue==12.1.2' \
      'msrestazure~=0.6.3' \
      'psycopg2-binary==2.8.5'

# REQUIRED FOR AZURE:
# ENV AZURE_STORAGE_ACCESS_KEY DefaultEndpointsProtocol=https;EndpointSuffix=core.windows.net;AccountName=XXXXX;AccountKey=XXXXX
# ENV MLFLOW_SERVER_FILE_STORE /mnt/mlruns
# ENV MLFLOW_SERVER_DEFAULT_ARTIFACT_ROOT wasbs://mlflow@XXXXX.blob.core.windows.net/mlartifacts

ENV MLFLOW_SERVER_HOST 0.0.0.0
ENV MLFLOW_SERVER_PORT 5000
ENV MLFLOW_SERVER_WORKERS 4

CMD ["/entrypoint.sh"]
