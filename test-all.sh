#!/bin/sh

: ${GOPATH:?"Need to set GOPATH"}
option=$1

# Main binary must build
cd $GOPATH/src/github.com/rafaeljusto/shelter && go build shelter.go
return_code=$?
cd $GOPATH/src/github.com/rafaeljusto/shelter && rm -f shelter

# If there was any error building the main binary, we shouldn't continue
if [ $return_code -ne 0 ]; then
  exit $return_code
fi

# Unit tests
echo "\n[[ UNIT TESTS ]]\n"

cd $GOPATH/src/github.com/rafaeljusto/shelter && find . -type f -name '*_test.go' -exec dirname {} \; | sort -u | awk '{printf("github.com/rafaeljusto/shelter%s ", substr($1,2))}' | xargs -I@ sh -c "go install @ && go test -cover @"
return_code=$?

# If there was any error in the unit tests, we shouldn't run
# the integration tests!
if [ $return_code -ne 0 ]; then
  exit $return_code
fi

# Check for unit test only
if [ "$option" = "-unit" ]; then
  exit 0
fi

# Integration tests
echo "\n[[ INTEGRATION TESTS ]]\n"

cd $GOPATH/src/github.com/rafaeljusto/shelter && find . -type f -wholename './testing/*.go' | grep -v 'utils' | sort -u | awk -F '/' '{ print $0 " -config=\"" substr($0, 1, length($0)-3) ".conf\"" }' | xargs -I@ sh -c "go run -race @"
return_code=$?
rm -f scan.log

exit $return_code
