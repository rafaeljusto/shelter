#!/usr/bin/env python

import getopt
import sys
import subprocess
import urllib.request

class Language:
  def __init__(self):
    self.name = ""
    self.description = ""

class Region:
  def __init__(self):
    self.name = ""
    self.description = ""

class GroupType:
  Unknown = 0
  Language = 1
  Region = 2

def retrieveData(url):
  response = urllib.request.urlopen(url)
  data = response.read()
  response.close()
  return data.decode()

def buildLanguages(data):
  group = GroupType.Unknown
  languages = []
  regions = []

  currentLanguage = Language()
  currentRegion = Region()

  for line in data.split("\n"):
    if line.startswith("Type:"):
      if line[6:] == "language":
        group = GroupType.Language
        currentLanguage = Language()

      elif line[6:] == "region":
        group = GroupType.Region
        currentRegion = Region()

      else:
        group = GroupType.Unknown

    elif line.startswith("Subtag:"):
      if group == GroupType.Language:
        currentLanguage.name = line[8:]

      elif group == GroupType.Region:
        currentRegion.name = line[8:]

    elif line.startswith("Description:"):
      if group == GroupType.Language:
        if len(currentLanguage.description) > 0:
          currentLanguage.description += " and " + line[13:]

        else:
          currentLanguage.description = line[13:]
          languages.append(currentLanguage)

      if group == GroupType.Region:
        if len(currentRegion.description) > 0:
          currentRegion.description += " and " + line[13:]

        else:
          currentRegion.description = line[13:]
          regions.append(currentRegion)

  return languages, regions

def writeLanguages(languages, regions, outputPath):
  output = open(outputPath, "w")

  print("package model", file=output)
  print("""
// File generated using the language.py script, that is responsable for parsing the IANA
// Language Subtag Registry file obtained from 
// http://www.iana.org/assignments/language-subtag-registry/language-subtag-registry""", file=output)

  print("""
// List of possible language types""", file=output)
  print("const (", file=output)

  for language in languages:
    normalizeName = language.name.upper().replace(".", "")
    print ("  LanguageType" + normalizeName + " LanguageType = \"" + language.name + "\" "+
      "// " + language.description, file=output)

  print(")", file=output)

  print("""
// Used in the system everytime that we need to determinate a preferred language, the LanguageType
// is an enumerate that describe all possible languages""", file=output)
  print("type LanguageType string", file=output)

  print("""
// List of possible region types, that are a subcategory of the language""", file=output)
  print("const (", file=output)

  i = 1
  for region in regions:
    normalizeName = region.name.upper().replace(".", "")
    print ("  RegionType" + normalizeName + " RegionType = \"" + region.name + "\" "+
      "// " + region.description, file=output)
    i += 1

  print(")", file=output)

  print("""
// Used to determinate from a given language, what is the specific dialect (region) that we will
// assume""", file=output)
  print("type RegionType string", file=output)
  output.close()

  subprocess.call(["gofmt", "-w", outputPath])

###################################################################

defaultURL = "http://www.iana.org/assignments/language-subtag-registry/language-subtag-registry"

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
    print("Output path was not given")
    usage()
    sys.exit(1)

  try:
    data = retrieveData(url)
    languages, regions = buildLanguages(data)
    writeLanguages(languages, regions, outputPath)

  except KeyboardInterrupt:
    sys.exit(1)

if __name__ == "__main__":
  main(sys.argv[1:])
