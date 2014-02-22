#!/usr/bin/env python

import getopt
import sys
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

def retrieveFile(url):
  response = urllib.request.urlopen(url)
  data = response.read()
  response.close()
  return data

def buildLanguages(data):
  group = GroupType.Unknown
  languages = []
  regions = []

  currentLanguage = Language()
  currentRegion = Region()

  for line in data.decode().split("\n"):
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

def writeLanguages(languages, regions):
  print("package model")
  print("""
// File generated using the language.py script, that is responsable for parsing the IANA
// Language Subtag Registry file obtained from 
// http://www.iana.org/assignments/language-subtag-registry/language-subtag-registry""")

  print("""
// List of possible language types""")
  print("const (")

  i = 1
  for language in languages:
    normalizeName = language.name.upper().replace(".", "")
    print ("  LanguageType" + normalizeName + " LanguageType = " + str(i) + " // " + language.description)
    i += 1

  print(")")

  print("""
// Used in the system everytime that we need to determinate a preferred language, the LanguageType
// is an enumerate that describe all possible languages""")
  print("type LanguageType int")

  print("""
// List of possible region types, that are a subcategory of the language""")
  print("const (")

  i = 1
  for region in regions:
    normalizeName = region.name.upper().replace(".", "")
    print ("  RegionType" + normalizeName + " RegionType = " + str(i) + " // " + region.description)
    i += 1

  print(")")

  print("""
// Used to determinate from a given language, what is the specific dialect (region) that we will
// assume""")
  print("type RegionType int")

###################################################################

defaultURL = "http://www.iana.org/assignments/language-subtag-registry/language-subtag-registry"

def usage():
  print("")
  print("Usage: " + sys.argv[0] + " [-h|--help] [-u|--url]")
  print("  Where -h or --help is for showing this usage")
  print("        -u or --url is the URL of the source file")

def main(argv):
  try:
    opts, args = getopt.getopt(argv, "u", ["url"])

  except getopt.GetoptError as err:
    print(str(err))
    usage()
    sys.exit(1)

  url = ""

  for key, value in opts:
    if key in ("-u", "--url"):
      url = value

    elif key in ("-h", "--help"):
      usage()
      sys.exit(0)

  if len(url) == 0:
    url = defaultURL

  try:
    data = retrieveFile(url)
    languages, regions = buildLanguages(data)
    writeLanguages(languages, regions)

  except KeyboardInterrupt:
    sys.exit(1)

if __name__ == "__main__":
  main(sys.argv[1:])
