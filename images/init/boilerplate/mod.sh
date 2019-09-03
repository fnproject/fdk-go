#!/bin/bash
set -euo pipefail
go mod init "$FN_FUNCTION_NAME" > /dev/null 2>&1
go get > /dev/null 2>&1
tar c go.mod func.go func.init.yaml
