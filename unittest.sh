#!/usr/bin/env bash
set -euo pipefail

# DEBUG
go version
go env

# TEST
COVER_PROFILE=coverage.out
COVER_THRESHOLD=100

go mod download
go generate ./...
CGO_ENABLED=0 go test -coverprofile=${COVER_PROFILE} $(go list ./... | grep -v testings) -v

# CHECK THRESHOLD
COVERAGE=$(go tool cover -func=${COVER_PROFILE} | grep total | awk '{print substr($3, 1, length($3) - 1)}')
echo "Total Coverage is: ${COVERAGE} by threshold: ${COVER_THRESHOLD}."
echo "${COVERAGE} ${COVER_THRESHOLD}" | awk '{if (!($1 >= $2)) { print "Coverage: " $1 "%" ", Expected threshold: " $2 "%"; exit 1 } }'
go tool cover -html=${COVER_PROFILE}