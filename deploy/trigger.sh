#!/bin/bash

environment=${1:-$ENVIRONMENT}
version=${2:-$VERSION}
token=${3:-$K8S_TOKEN}
cluster=${4:-$K8S_CLUSTER}

if [ -z "$version" ]; then
  echo "VERSION should not be empty. Exiting with error."
  exit 1
fi

config_file=k8s-$environment-$version.yml

environment=$environment version=$version envsubst < deploy/$environment/deployment.yml  > $config_file

## apply deployment
kubectl -n rentals \
  apply --record -s $cluster \
  --token=$token \
  -f $config_file \
  --insecure-skip-tls-verify || exit 1

## check if the deploy was successful or not
kubectl -n rentals \
  rollout status deployment/go-boilerplate-api \
  -s $cluster \
  --token=$token \
  --insecure-skip-tls-verify || {
## if failed then undo it
kubectl -n rentals \
  rollout undo \
  -s $cluster \
  --token=$token \
  -f $config_file \
  --insecure-skip-tls-verify; exit 1; }
