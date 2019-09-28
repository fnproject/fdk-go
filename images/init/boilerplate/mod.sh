#!/bin/sh
set -euo pipefail

# Pipe stdout to file, as Fn CLI init-image reads stdout to untar (outputting messages to stderr is fine).
go mod init "$FN_FUNCTION_NAME" > /tmp/go_mod_output
go get > /tmp/go_get_output

tar c go.mod func.go func.init.yaml
