#!/bin/bash -eux

pushd dp-collection-api
  make build
  cp build/dp-collection-api Dockerfile.concourse ../build
popd
