#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -v -coverprofile=profile.out -coverpkg ./... $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

bash <(curl -s https://codecov.io/bash)

rm coverage.txt
