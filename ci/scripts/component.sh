#!/bin/bash -eux

pushd dp-collection-api
  export MEMONGO_DOWNLOAD_URL=https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-debian10-4.2.15.tgz
  make test-component
popd
