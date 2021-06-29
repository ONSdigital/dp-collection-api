#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-collection-api
  make lint
popd
