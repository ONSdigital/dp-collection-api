#!/bin/bash -eux

pushd dp-collection-api
  make test-component
popd
