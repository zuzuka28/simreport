#!/usr/bin/env sh
#
# install s3 into k8s

helm repo add minio https://charts.min.io/

helm install \
    --set replicas=1 \
    --set mode=standalone \
    --set resources.requests.memory=512Mi \
    --create-namespace \
    -n minio minio minio/minio
