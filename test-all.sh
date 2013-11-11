#!/bin/sh

: ${GOPATH:?"Need to set GOPATH"}
current_dir=$PWD

# Unit tests

cd $GOPATH/src/shelter && find . -type f -name '*_test.go' -exec dirname {} \; | sort -u | awk '{printf("shelter%s ", substr($1,2))}' | xargs -I@ sh -c "go install @ && go test @"
return_code=$?

# If there was any error in the unit tests, we shouldn't run
# the integration tests!
if [ $return_code -ne 0 ]; then
  cd $current_dir
  exit $return_code
fi

# Integration tests

cd $GOPATH/src/shelter && find . -type f -wholename './integration_tests/*.go' | sort -u | awk -F '/' '{ print $0 " -config=\"" substr($0, 0, length($0)-2) ".conf\"" }' | xargs -I@ sh -c "go run @"
return_code=$?

cd $current_dir
exit $return_code
