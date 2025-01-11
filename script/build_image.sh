#!/usr/bin/env sh
#
# build image

REGISTRY="${REGISTRY:-zuzuka28}"

if [ -z "$SERVICE_NAME" ]; then
    echo "SERVICE_NAME is empty" && exit 1
fi

if [ -z "$LANG" ]; then
    echo "LANG is empty" && exit 1
fi

TAG="$(date +'%Y-%m-%d_%H-%M')_$(git rev-parse HEAD)"

docker build \
    -f "./prj/${SERVICE_NAME}/build/docker/Dockerfile" \
    -t "${REGISTRY}/${SERVICE_NAME}":"${TAG}"  \
    "prj/${SERVICE_NAME}"
