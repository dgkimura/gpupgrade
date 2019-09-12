#!/bin/bash

set -ex

# make depend fails if run as gpadmin with a dep ensure git-remote-https signal 11 failure
export GOPATH="$PWD/go"
export PATH="$PWD/go/bin:$PATH"

ls -al go/pkg/dep/*

make -C go/src/github.com/greenplum-db/gpupgrade
