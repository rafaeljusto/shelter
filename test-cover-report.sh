#!/bin/sh

: ${GOPATH:?"Need to set GOPATH"}

cd $GOPATH/src/github.com/rafaeljusto/shelter/ && find . -type f -name '*_test.go' -exec dirname {} \; | sort -u | awk '{ print "github.com/rafaeljusto/shelter" substr($1,2) }' | xargs -I@ sh -c "go install @ && go test -coverprofile='cover-profile.out' -cover @ && go tool cover -html=cover-profile.out"
rm -f cover-profile.out
