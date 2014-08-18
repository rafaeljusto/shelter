#!/usr/bin/env python

# Copyright 2014 Rafael Dantas Justo. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

import fnmatch
import os
import subprocess
import sys

def initialChecks():
  if "GOPATH" not in os.environ:
    print("Need to set GOPATH")
    sys.exit(1)

def findPath():
  goPath = os.environ["GOPATH"]
  goPathParts = goPath.split(":")
  for goPathPart in goPathParts:
    projectPath = os.path.join(goPathPart, "src", "github.com",
      "rafaeljusto", "shelter")
    if os.path.exists(projectPath):
      return projectPath

  return ""

def changePath():
  projectPath = findPath()
  if len(projectPath) == 0:
    print("Project not found")
    sys.exit(1)

  os.chdir(projectPath)

def runCoverReport():
  print("\n[[ UNIT TESTS ]]\n")

  goPackages = []
  for root, dirnames, filenames in os.walk("."):
    for filename in fnmatch.filter(filenames, "*_test.go"):
      # TODO: We should test this in Windows
      goPackage = "github.com/rafaeljusto/shelter" + root[1:]
      goPackages.append(goPackage)

  goPackages = set(goPackages)
  goPackages = list(goPackages)
  goPackages.sort()

  success = True

  for goPackage in goPackages:
    try:
      subprocess.check_call(["go", "install", goPackage])
      subprocess.check_call(["go", "test", "-coverprofile=cover-profile.out", "-cover", goPackage])
      subprocess.check_call(["go", "tool", "cover", "-html=cover-profile.out"])
    except subprocess.CalledProcessError:
      success = False

  # Remove the temporary file created for the
  # covering reports
  try:
    os.remove("cover-profile.out")
  except OSError:
    pass

  if not success:
    print("Errors during the unit test execution")
    sys.exit(1)

###################################################################

if __name__ == "__main__":
  try:
    initialChecks()
    changePath()
    runCoverReport()
  except KeyboardInterrupt:
    sys.exit(1)
