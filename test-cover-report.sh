#!/bin/sh

: ${GOPATH:?"Need to set GOPATH"}
current_dir=$PWD

cd $GOPATH/src/shelter && find . -type f -name '*_test.go' -exec dirname {} \; | sort -u | awk '{ print "shelter" substr($1,2) }' | xargs -I@ sh -c "go install @ && go test -coverprofile='cover-profile.out' -cover @ && go tool cover -html=cover-profile.out"
rm -f cover-profile.out

cd $current_dir
