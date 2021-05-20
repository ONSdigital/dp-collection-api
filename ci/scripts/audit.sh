#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-collection-api
  make audit
popd
