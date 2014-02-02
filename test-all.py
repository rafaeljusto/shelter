#!/usr/bin/env python

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
  goPathParts = goPath.split(";")
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

  except OSError:
    pass

  # TODO: In windows the main binary is probably 
  #       generated as shelter.exe, we should also
  #        remove it

def runUnitTests(benchmark):
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

      if benchmark:
        subprocess.check_call(["go", "test", "-cover", "-bench", ".", "-benchmem", goPackage])
      else:
        subprocess.check_call(["go", "test", "-cover", goPackage])

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

###################################################################

def usage():
  print("")
  print("Usage: " + sys.argv[0] + " [-h|--help] [-u|--unit] [-b|--bench]")
  print("  Where -h or --help is for showing this usage")
  print("        -u or --unit is to run only the unit tests")
  print("        -b or --bench is to run unit tests with benchmark")

def main(argv):
  try:
    opts, args = getopt.getopt(argv, "ub", ["unit", "bench"])

  except getopt.GetoptError as err:
    print(str(err))
    usage()
    sys.exit(1)

  unitTestOnly = False
  benchmark = False

  for key, value in opts:
    if key in ("-u", "--unit"):
      unitTestOnly = True

    elif key in ("-b", "--bench"):
      benchmark = True

    elif key in ("-h", "--help"):
      usage()
      sys.exit(0)

  try:
    initialChecks()
    changePath()
    buildMainBinary()
    runUnitTests(benchmark)

    if not unitTestOnly:
      runIntegrationTests()

  except KeyboardInterrupt:
    sys.exit(1)

if __name__ == "__main__":
  main(sys.argv[1:])
