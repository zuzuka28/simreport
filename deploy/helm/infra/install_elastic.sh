#!/usr/bin/env sh
#
# install elasticsearch and kibana into k8s

helm repo add elastic https://helm.elastic.co
helm repo update

helm install elastic-operator-crds elastic/eck-operator-crds

helm install elastic-operator elastic/eck-operator -n elastic --create-namespace \
    --set=installCRDs=false \
    --set=managedNamespaces='{elastic}' \
    --set=createClusterScopedResources=false \
    --set=webhook.enabled=false \
    --set=config.validateStorageClass=false

cat <<EOF | kubectl apply -f -
apiVersion: elasticsearch.k8s.elastic.co/v1
kind: Elasticsearch
metadata:
  name: alpha
  namespace: elastic
spec:
  version: 8.17.0
  nodeSets:
  - name: default
    count: 1
    config:
      node.store.allow_mmap: false
EOF

ELASTIC_PASSWORD="$(kubectl get secret -n elastic alpha-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')"

export ELASTIC_PASSWORD
