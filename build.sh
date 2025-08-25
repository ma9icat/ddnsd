#! /usr/bin/env sh

set -e

if [ -f .env ]; then
    export $(cat .env | grep -v '#' | awk '/=/ {print $1}')
fi

REGISTRY=${REGISTRY:-"ddnsd"}

function build() {
  local -a args
  args+=(--file Dockerfile)
  args+=(--tag "${REGISTRY}$(basename $PWD):latest")
  args+=(.)
  docker build "${args[@]}"
}

function push() {
  docker push "${REGISTRY}$(basename $PWD):latest"
}

function main() {
  build
  push
}

main "$@"