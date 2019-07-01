#!/usr/bin/env bash


main() {
  run_gofmt
  run_golint
  run_govet
}

function log_info() {
  now=$(date)
  printf "$now INFO: $@\n"
}


function log_error() {
  now=$(date)
  printf "$now ERROR: $@\n"
}

function log_hint() {
  now=$(date)
  printf "$now >>>: $@\n"
}

function run_gofmt() {
  GOFMT_FILES=$(gofmt -l .)
  if [[ -n "$GOFMT_FILES" ]]; then
    log_error "gofmt failed for the following files: \n$GOFMT_FILES"
    log_hint "please run 'gofmt -w .' on your changes before committing."
    exit 1
  fi
}

function run_golint() {
  GOLINT_ERRORS=$(golint ./... | grep -v "Id should be")
  if [[ -n "$GOLINT_ERRORS" ]]; then
    log_error "golint failed for the following reasons: \n$GOLINT_ERRORS"
    log_hint "please run 'golint ./...' on your changes before committing."
    exit 1
  fi
}

function run_govet() {
  GOVET_ERRORS=$(go vet ./*.go 2>&1)
  if [[ -n "$GOVET_ERRORS" ]]; then
    log_error "go vet failed for the following reasons: \n$GOVET_ERRORS"
    log_hint "please run 'go vet ./*.go' on your changes before committing."
    exit 1
  fi
}

main
