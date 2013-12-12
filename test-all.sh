#!/bin/sh

: ${GOPATH:?"Need to set GOPATH"}
current_dir=$PWD

# Unit tests
echo "\n[[ UNIT TESTS ]]\n"

cd $GOPATH/src/shelter && find . -type f -name '*_test.go' -exec dirname {} \; | sort -u | awk '{printf("shelter%s ", substr($1,2))}' | xargs -I@ sh -c "go install @ && go test -cover @"
return_code=$?

# If there was any error in the unit tests, we shouldn't run
# the integration tests!
if [ $return_code -ne 0 ]; then
  cd $current_dir
  exit $return_code
fi

# Integration tests
echo "\n[[ INTEGRATION TESTS ]]\n"

cd $GOPATH/src/shelter && find . -type f -wholename './testing/*.go' | grep -v 'utils' | sort -u | awk -F '/' '{ print $0 " -config=\"" substr($0, 1, length($0)-3) ".conf\"" }' | xargs -I@ sh -c "go run -race @"
return_code=$?
rm -f scan.log

cd $current_dir
exit $return_code
