#!/usr/bin/env python

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
    jsTestPath = os.path.join(goPathPart, "src", "github.com", 
      "rafaeljusto", "shelter", "templates", "client", "tests")
    if os.path.exists(jsTestPath):
      return jsTestPath

  return ""

def changePath():
  projectPath = findPath()
  if len(projectPath) == 0:
    print("Project not found")
    sys.exit(1)

  os.chdir(projectPath)

def runJSTests():
  try:
    subprocess.check_call(["karma", "start", "karma.conf.js", "--single-run"])

  except:
    print("Errors during the JS test execution")
    sys.exit(1)

def main(argv):
  try:
    initialChecks()
    changePath()
    runJSTests()

  except KeyboardInterrupt:
    sys.exit(1)

if __name__ == "__main__":
  main(sys.argv[1:])
