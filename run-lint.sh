#!/usr/bin/env bash

set -euo pipefail

[[ -n ${DEBUG:-} ]] && set -x

gopath="$(go env GOPATH)"

if ! [[ -x "$gopath/bin/golangci-lint" ]]; then
	echo >&2 'Installing golangci-lint'
	curl --silent --fail --location \
		https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$gopath/bin" v1.43.0
fi

# configured by .golangci.yml
"$gopath/bin/golangci-lint" run
