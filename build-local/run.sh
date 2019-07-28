#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SCRIPT_DIR=${DIR}

UNAME_OUT="$(uname -s)"
case "${UNAME_OUT}" in
    Linux*)     machine=linux_amd64;;
    Darwin*)    machine=darwin_amd64;;
    *)          machine="UNKNOWN:${UNAME_OUT}"
esac

${SCRIPT_DIR}/../dist/${machine}/memcache
