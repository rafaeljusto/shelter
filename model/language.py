#!/usr/bin/env python

# Copyright 2014 Rafael Dantas Justo. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

# Script responsable for building the language.go file from the IANA Language Subtag Registry file
# format. We should run this only when we found a new language that should be used in the system and
# is not yet listed in our Go file

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

  print("""// model - Description of the objects
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.""", file=output)
  print("package model", file=output)
  print("""
// File generated using the language.py script, that is responsable for parsing the IANA
// Language Subtag Registry file obtained from
// http://www.iana.org/assignments/language-subtag-registry/language-subtag-registry""", file=output)

  print("""import (
  "strings"
)""", file=output)

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
// Structure used to identify if a language exists or not
var (
  languageTypes map[LanguageType]bool = map[LanguageType]bool{""", file=output)

  for language in languages:
    normalizeName = language.name.upper().replace(".", "")
    print ("    LanguageType" + normalizeName + ": true,", file=output)

  print("""}
)""", file=output)

  print("""
// Used to verify if a language is valid or not. Don't need lock because we have many
// readers and no writers. Rob Pike sad that it's ok
// (https://groups.google.com/forum/#!msg/golang-nuts/HpLWnGTp-n8/hyUYmnWJqiQJ)
func LanguageTypeExists(languageType string) bool {
  // Normalize input
  languageType = strings.ToLower(languageType)
  languageType = strings.TrimSpace(languageType)

  _, ok := languageTypes[LanguageType(languageType)]
  return ok
}""", file=output)

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

  print("""
// Structure used to identify if a region exists or not
var (
  regionTypes map[RegionType]bool = map[RegionType]bool{""", file=output)

  for region in regions:
    normalizeName = region.name.upper().replace(".", "")
    print ("    RegionType" + normalizeName + ": true,", file=output)

  print("""}
)""", file=output)

  print("""
// Used to verify if a region is valid or not. Don't need lock because we have many
// readers and no writers. Rob Pike sad that it's ok
// (https://groups.google.com/forum/#!msg/golang-nuts/HpLWnGTp-n8/hyUYmnWJqiQJ)
func RegionTypeExists(regionType string) bool {
  // Normalize input
  regionType = strings.ToUpper(regionType)
  regionType = strings.TrimSpace(regionType)

  _, ok := regionTypes[RegionType(regionType)]
  return ok
}

// Useful function to check if a language with region or not is valid
func IsValidLanguage(language string) bool {
  if LanguageTypeExists(language) {
    return true
  }

  languageParts := strings.Split(language, "-")
  if len(languageParts) != 2 {
    return false
  }

  return LanguageTypeExists(languageParts[0]) && RegionTypeExists(languageParts[1])
}""", file=output)

  output.close()

  subprocess.call(["gofmt", "-w", outputPath])

###################################################################

defaultURL = "http://www.iana.org/assignments/language-subtag-registry/language-subtag-registry"
defaultOutput = "language.go"

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
    languages, regions = buildLanguages(data)
    writeLanguages(languages, regions, outputPath)

  except KeyboardInterrupt:
    sys.exit(1)

if __name__ == "__main__":
  main(sys.argv[1:])
