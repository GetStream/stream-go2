#!/usr/bin/env bash

set -euo pipefail

[[ -n ${DEBUG:-} ]] && set -x

gopath="$(go env GOPATH)"

if ! [[ -x "$gopath/bin/golangci-lint" ]]; then
	echo >&2 'Installing golangci-lint'
	curl --silent --fail --location \
		https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$gopath/bin" v1.24.0
fi

# configured by .golangci.yml
"$gopath/bin/golangci-lint" run

install_impi() {
	impi_dir="$(mktemp -d)"
	trap 'rm -rf -- ${impi_dir}' EXIT

	GO111MODULE=off GOPATH="${impi_dir}" \
		GOBIN="${gopath}/bin" \
		go get github.com/pavius/impi/cmd/impi
}

# install impi that ensures import grouping is done consistently
if ! [[ -x "${gopath}/bin/impi" ]]; then
	echo >&2 'Installing impi'
	install_impi
fi

"$gopath/bin/impi" \
	--local github.com/GetStream/stream-go2 \
	--scheme stdThirdPartyLocal \
	./...
