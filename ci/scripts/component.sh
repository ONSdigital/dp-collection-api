#!/bin/bash -eux

pushd dp-collection-api

  echo "$(</etc/os-release)"

  make test-component
popd
