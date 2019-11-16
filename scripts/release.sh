#!/bin/bash -e

clrEcho() {
  CODE=$1
  echo -e "\033[${CODE}${@:2}\033[0m"
}

if [ -z "${TAG}" ]; then
  read -p "tag: " TAG
fi

[ -z "${TAG}" ] && clrEcho 031m tag is empty && exit 1

git tag ${TAG} && git push origin ${TAG}