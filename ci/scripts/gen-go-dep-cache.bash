#!/bin/bash

set -ex

# If it exists, show contents of cache directory
if [ -d go/pkg/dep/sources ]; then
    echo "Concourse cache directory: go/pkg/dep/sources"
    ls -al go/pkg/dep/sources
else
    echo "go/pkg/dep/sources does not exist"
fi

# make depend fails if run as gpadmin with a dep ensure git-remote-https signal 11 failure
export GOPATH="$PWD/go"
export PATH="$PWD/go/bin:$PATH"

time make -C go/src/github.com/greenplum-db/gpupgrade dep

if [ -d go/pkg/dep/sources ]; then
    echo "Concourse cache directory: go/pkg/dep/sources"
    ls -al go/pkg/dep/sources
else
    echo "go/pkg/dep/sources does not exist"
fi

tar zcf go/pkg/go-dep-cache.tgz -C go/pkg dep
