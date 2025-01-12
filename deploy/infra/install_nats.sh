#!/usr/bin/env sh
#
# install nats into k8s


helm repo add nats https://nats-io.github.io/k8s/helm/charts/

cat <<EOF | helm upgrade --install -n nats --create-namespace nats nats/nats -f -
    config:
      cluster:
        enabled: true
        replicas: 3
      jetstream:
        enabled: true
        fileStore:
          pvc:
            size: 1Gi

    podTemplate:
      topologySpreadConstraints:
        kubernetes.io/hostname:
          maxSkew: 1
          whenUnsatisfiable: DoNotSchedule
EOF
