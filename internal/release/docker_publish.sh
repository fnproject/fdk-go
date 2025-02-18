#!/usr/bin/env bash

set -eu

REGCTL_BIN=regctl
# Test regctl is on path
$REGCTL_BIN --help

TEMPDIR=$(mktemp -d)
cd "${TEMPDIR}"

function cleanup {
    rm -rf "${TEMPDIR}"
}
trap cleanup EXIT

{
  $REGCTL_BIN image copy iad.ocir.io/oraclefunctionsdevelopm/fnproject/go:1.23 docker.io/fnproject/go:1.23;
  $REGCTL_BIN image copy iad.ocir.io/oraclefunctionsdevelopm/fnproject/go:1.23-dev docker.io/fnproject/go:1.23-dev;
  $REGCTL_BIN image copy iad.ocir.io/oraclefunctionsdevelopm/fnproject/go:1.24 docker.io/fnproject/go:1.24;
  $REGCTL_BIN image copy iad.ocir.io/oraclefunctionsdevelopm/fnproject/go:1.24-dev docker.io/fnproject/go:1.24-dev;
}
