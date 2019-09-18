#!/bin/bash

set -ex

# If it exists, show contents of cache directory
if [ -d go/pkg/dep ]; then
    echo "Concourse cache directory: go/pkg/dep"
    ls -al go/pkg/dep
else
    echo "go/pkg/dep does not exist"
fi

# make depend fails if run as gpadmin with a dep ensure git-remote-https signal 11 failure
export GOPATH="$PWD/go"
export PATH="$PWD/go/bin:$PATH"

make -C go/src/github.com/greenplum-db/gpupgrade dep

tar zcf go/pkg/go-dep-cache.tgz -C go/pkg dep
