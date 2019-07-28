#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SCRIPT_DIR=${DIR}

goreleaser --config=${SCRIPT_DIR}/../.goreleaser.yml --skip-publish --snapshot --rm-dist
