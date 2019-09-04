#!/bin/sh
set -euo pipefail
go mod init "$FN_FUNCTION_NAME" > /tmp/go_mod_output 2>&1
go get > /tmp/go_get_output 2>&1
tar c go.mod func.go func.init.yaml
