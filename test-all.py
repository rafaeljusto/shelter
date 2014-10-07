#!/usr/bin/env python

# Copyright 2014 Rafael Dantas Justo. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

import fnmatch
import getopt
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

def buildMainBinary():
  try:
    subprocess.check_call(["go", "build", "shelter.go"])

  except subprocess.CalledProcessError:
    print("Error building main binary")
    sys.exit(1)

  try:
    os.remove("shelter")
    os.remove("shelter.exe")

  except OSError:
    pass

def runUnitTests(benchmark):
  print("\n[[ UNIT TESTS ]]\n")

  success = True

  try:
    subprocess.check_call(["go", "install", "./..."])

    # Will turn on only when all problems are solved
    #subprocess.check_call(["go", "vet", "./..."])

    if benchmark:
      subprocess.check_call(["go", "test", "-cover", "-bench", ".", "-benchmem", "./..."])
    else:
      subprocess.check_call(["go", "test", "-cover", "./..."])

  except subprocess.CalledProcessError:
    success = False

  if not success:
    print("Errors during the unit test execution")
    sys.exit(1)

def runIntegrationTests():
  print("\n[[ INTEGRATION TESTS ]]\n")

  testFiles = []
  for root, dirnames, filenames in os.walk("testing"):
    for filename in fnmatch.filter(filenames, "*.go"):
      # Ignore utils directory
      if "utils" in root:
        continue

      testFiles.append(os.path.join(root, filename))

  testFiles = set(testFiles)
  testFiles = list(testFiles)
  testFiles.sort()

  success = True

  for testFile in testFiles:
    config = "-config=" + testFile[:len(testFile)-3] + ".conf"
    try:
      subprocess.check_call(["go", "run", "-race", testFile, config])

    except:
      success = False

  # One of the integration tests creates a temporary
  # log file, we just need to remove it
  try:
    os.remove("scan.log")

  except OSError:
    pass

  if not success:
    print("Errors during the integration test execution")
    sys.exit(1)

def runInterfaceTests():
  print("\n[[ INTERFACE TESTS ]]\n")

  success = True
  karmaConf = os.path.join("templates", "client", "tests", "karma.conf.js")
  commandLine = subprocess.list2cmdline(["karma", "start", karmaConf, "--single-run"])

  try:
    subprocess.check_call([commandLine], shell=True)

  except:
    success = False

  if not success:
    print("Errors during the interface test execution")
    sys.exit(1)

###################################################################

def usage():
  print("")
  print("Usage: " + sys.argv[0] + " [-h|--help] [-u|--unit] [-b|--bench] [-n|--integration] [-i|--interface]")
  print("  Where -h or --help is for showing this usage")
  print("        -u or --unit is to run only the unit tests")
  print("        -b or --bench is to run unit tests with benchmark")
  print("        -n or --integration is to run integration tests")
  print("        -i or --interface is to run only the interface tests")

def main(argv):
  try:
    opts, args = getopt.getopt(argv, "ubni", ["unit", "bench", "integration", "interface"])

  except getopt.GetoptError as err:
    print(str(err))
    usage()
    sys.exit(1)

  unitTestOnly = False
  integrationTestOnly = False
  interfaceTestOnly = False
  benchmark = False

  for key, value in opts:
    if key in ("-u", "--unit"):
      unitTestOnly = True

    elif key in ("-b", "--bench"):
      benchmark = True

    elif key in ("-n", "--integration"):
      integrationTestOnly = True

    elif key in ("-i", "--interface"):
      interfaceTestOnly = True

    elif key in ("-h", "--help"):
      usage()
      sys.exit(0)

  try:
    initialChecks()
    changePath()
    buildMainBinary()

    if unitTestOnly or not (integrationTestOnly or interfaceTestOnly):
      runUnitTests(benchmark)

    if integrationTestOnly or not (unitTestOnly or interfaceTestOnly):
      runIntegrationTests()

    if interfaceTestOnly or not (unitTestOnly or integrationTestOnly):
      runInterfaceTests()

  except KeyboardInterrupt:
    sys.exit(1)

if __name__ == "__main__":
  main(sys.argv[1:])
