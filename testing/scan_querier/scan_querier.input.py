#!/usr/bin/env python

# Copyright 2014 Rafael Dantas Justo. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

import getopt
import sys
import subprocess
import urllib.request

class NS:
  def __init__(self):
    self.name = ""
    self.type = "NS"
    self.namserver = ""

  def __str__(self):
    return "{} {} {}".format(self.name, self.type, self.namserver)

class DS:
  def __init__(self):
    self.name = ""
    self.type = "DS"
    self.keytag = 0
    self.algorithm = 0
    self.digestType = 0
    self.digest = ""

  def __str__(self):
    return "{} {} {} {} {} {}".format(self.name, self.type, self.keytag,
      self.algorithm, self.digestType, self.digest)

class A:
  def __init__(self):
    self.name = ""
    self.type = "A"
    self.address = ""

  def __str__(self):
    return "{} {} {}".format(self.name, self.type, self.address)

class AAAA:
  def __init__(self):
    self.name = ""
    self.type = "AAAA"
    self.address = ""

  def __str__(self):
    return "{} {} {}".format(self.name, self.type, self.address)


def retrieveData(url):
  response = urllib.request.urlopen(url)
  data = response.read()
  response.close()
  return data.decode()

def buildZone(data):
  zone = []

  for line in data.split("\n"):
    lineParts = line.split()
    if len(lineParts) < 4:
      print(line)
      continue

    if lineParts[3] == "NS" and len(lineParts) == 5:
      ns = NS()
      ns.name = lineParts[0]
      ns.namserver = lineParts[4]
      zone.append(ns)

    elif lineParts[3] == "A" and len(lineParts) == 5:
      a = A()
      a.name = lineParts[0]
      a.address = lineParts[4]
      zone.append(a)

    elif lineParts[3] == "AAAA" and len(lineParts) == 5:
      aaaa = AAAA()
      aaaa.name = lineParts[0]
      aaaa.address = lineParts[4]
      zone.append(aaaa)

    elif lineParts[3] == "DS" and len(lineParts) == 8:
      ds = DS()
      ds.name = lineParts[0]
      ds.keytag = int(lineParts[4])
      ds.algorithm = int(lineParts[5])
      ds.digestType = int(lineParts[6])
      ds.digest = lineParts[7]
      zone.append(ds)

  return zone

def writeZone(zone, outputPath):
  output = open(outputPath, "w")

  for rr in zone:
    print(str(rr), file=output)

  output.close()

###################################################################

defaultURL = "http://www.internic.net/domain/root.zone"
defaultOutput = "scan_querier.input"

def usage():
  print("")
  print("Usage: " + sys.argv[0] + " [-h|--help] [-u|--url] [-o|--output]")
  print("  Where -h or --help is for showing this usage")
  print("        -u or --url is the URL of the source file")
  print("        -o or --output is the path where the Go code will written")

def main(argv):
  try:
    opts, args = getopt.getopt(argv, "u:o:", ["url", "output"])

  except getopt.GetoptError as err:
    print(str(err))
    usage()
    sys.exit(1)

  url = ""
  outputPath = ""

  for key, value in opts:
    if key in ("-u", "--url"):
      url = value

    elif key in ("-o", "--output"):
      outputPath = value

    elif key in ("-h", "--help"):
      usage()
      sys.exit(0)

  if len(url) == 0:
    url = defaultURL

  if len(outputPath) == 0:
    outputPath = defaultOutput

  try:
    data = retrieveData(url)
    rootZone = buildZone(data)
    writeZone(rootZone, outputPath)

  except KeyboardInterrupt:
    sys.exit(1)

if __name__ == "__main__":
  main(sys.argv[1:])
