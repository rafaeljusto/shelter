#!/bin/sh

# Copyright 2014 Rafael Dantas Justo. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

# Travis CI tool doesn't support python scripts in a Go project,
# so this script was created thinking on the limited enviroment
# of the CI tool, if you're testing the project, please use the
# Python script (test-all.py)

: ${GOPATH:?"Need to set GOPATH"}

workspace=`echo $GOPATH | cut -d: -f1`
cd $workspace/src/github.com/rafaeljusto/shelter

# Main binary must build
go build shelter.go
return_code=$?
rm -f shelter

# If there was any error building the main binary, we shouldn't continue
if [ $return_code -ne 0 ]; then
  exit $return_code
fi

# Unit tests
echo "\n[[ UNIT TESTS ]]\n"

go install ./... && go test -cover ./...
return_code=$?

# If there was any error in the unit tests, we shouldn't run
# the integration tests!
if [ $return_code -ne 0 ]; then
  exit $return_code
fi

# For now on we are only testing the unit and interface layers in Travis. Because are tests takes
# too long and consume resources that Travis don't have. The full test will be performed by the
# integration tool.

# Integration tests
#echo "\n[[ INTEGRATION TESTS ]]\n"

#find . -type f -wholename './testing/*.go' | grep -v 'utils' | sort -u | awk -F '/' '{ print $0 " -config=\"" substr($0, 1, length($0)-3) ".conf\"" }' | xargs -I@ sh -c "go run -race @"
#return_code=$?
#rm -f scan.log

# If there was any error in the integration tests, we shouldn't
# run the interface tests!
#if [ $return_code -ne 0 ]; then
#  exit $return_code
#fi

# Interface tests
echo "\n[[ INTERFACE TESTS ]]\n"
./node_modules/karma/bin/karma start templates/client/tests/karma.conf.js --single-run
return_code=$?

exit $return_code
