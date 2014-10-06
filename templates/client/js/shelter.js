/*
 * Copyright 2014 Rafael Dantas Justo. All rights reserved.
 * Use of this source code is governed by a GPL
 * license that can be found in the LICENSE file.
 */

// Screen refresh rate in seconds
var refreshRate = 5;

var dnskeyFlags = [
  {id: 256, name:"ZSK"},
  {id: 257, name:"KSK"},
];

var algorithms = [
  {id: 1, name:"RSA/MD5"},
  {id: 2, name:"DH"},
  {id: 3, name:"DSA/SHA1"},
  {id: 4, name:"ECC"},
  {id: 5, name:"RSA/SHA1"},
  {id: 6, name:"DSA/SHA1-NSEC3"},
  {id: 7, name:"RSA/SHA1-NSEC3"},
  {id: 8, name:"RSA/SHA256"},
  {id: 10, name:"RSA/SHA512"},
  {id: 12, name:"GOST R"},
  {id: 13, name:"ECDSA/SHA256"},
  {id: 14, name:"ECDSA/SHA384"}
];

var dsDigestTypes = [
  {id: 1, name: "SHA1"},
  {id: 2, name: "SHA256"},
  {id: 3, name: "GOST94"},
  {id: 4, name: "SHA384"},
  {id: 5, name: "SHA512"},
];

var ownerLanguages = [
  "en-US",
  "pt-BR",
  "es-ES"
];

var emptyNameserver = {
  host: "",
  ipv4: "",
  ipv6: ""
};

var emptyDNSKEY = {
  flags: 257,
  algorithm: 8,
  publickKey: ""
};

var emptyDS = {
  keytag: "",
  algorithm: 8,
  digestType: 2,
  digest: ""
};

var emptyOwner = {
  email: "",
  language: "en-US"
};

var emptyDomain = {
  fqdn: "",
  nameservers: [
    angular.copy(emptyNameserver),
    angular.copy(emptyNameserver)
  ],
  dnskeys: [],
  dsset: [
    angular.copy(emptyDS)
  ],
  owners: [
    angular.copy(emptyOwner)
  ]
};

var emptyScan = {
  status: "NOTCHECKED",
  startedAt: "0000-00-00T00:00:00.000Z",
  finishedAt: "0000-00-00T00:00:00.000Z",
  domainsScanned: 0,
  domainsWithDNSSECScanned: 0,
  nameserverStatistics: {},
  dsStatistics: {}
};

// Look for a specific link href given a type. According to W3C link can have many types
// for a given URL (http://www.w3.org/TR/html401/types.html#type-links), so we use this
// function to make our life easier
function findLink(links, expectedType) {
  var href = "";
  if (!links) {
    return href;
  }

  links.forEach(function(link) {
    if (!link.types) {
      return;
    }

    link.types.forEach(function(type) {
      if (type == expectedType) {
        href = link.href;
      }
    });
  });

  return href;
}

function verificationResponseToHTML(data) {
  var result = "";

  if (!data) {
    return result;
  }

  if (data.fqdn) {
    result += "<h3>" + data.fqdn + "</h3><hr/>";
  }

  if (data.nameservers) {
    result += "<h3>NS</h3>" +
      "<table style='margin:auto'>";

    data.nameservers.forEach(function(ns) {
      result += "<tr>" +
        "<th style='text-align:left'>" + ns.host + "</th>" + 
        "<td>" + ns.lastStatus + "</td></tr>";
    });

    result += "</table>";
  }

  if (data.dsset) {
    result += "<h3>DS</h3>" +
      "<table style='margin:auto'>";

    data.dsset.forEach(function(ds) {
      result += "<tr>" +
        "<th style='text-align:left'>" + ds.keytag + "</th>" + 
        "<td>" + ds.lastStatus + "</td></tr>";
    });

    result += "</table>";
  }

  return result;
}

angular.module("shelter", ["ngAnimate", "ngCookies", "pascalprecht.translate"])

  .config(function($translateProvider, $httpProvider, $anchorScrollProvider) {
    $translateProvider.useStaticFilesLoader({
      prefix: "/languages/",
      suffix: ".json"
    });

    $translateProvider.fallbackLanguage("en_US");
    $translateProvider.determinePreferredLanguage();
    $translateProvider.useLocalStorage();

    $httpProvider.defaults.headers.common = {
      "Accept-Language": function() {
        if ($translateProvider.use() == undefined) {
          return $translateProvider.use();
        } else {
          return $translateProvider.use().replace("_", "-");
        }
      }
    };

    $anchorScrollProvider.disableAutoScrolling();
  })

  .filter("range", function() {
    return function(input, total) {
      total = parseInt(total);
      for (var i = 0; i < total; i++) {
        input.push(i);
      }
      return input;
    };
  })

  .filter("datetime", function() {
    return function(input, language) {
      // We will keep the timezone from what we received, but we aren't showing the
      // timezone in the result string. This could be a problem!

      // http://momentjs.com/docs/#/parsing/string-formats/
      // Note: Parsing multiple formats is considerably slower than parsing a single format. If you
      // can avoid it, it is much faster to parse a single format.
      var datetime = moment(input, "YYYY-MM-DDTHH:mm:ss.SSSZ").parseZone();
      if (datetime.isValid() && language) {
        // Detect empty datetime
        if (datetime.unix() <= 0) {
          return ""
        }

        language = language.replace("_", "-");
        moment.lang(language);
        return datetime.format("LLLL");
      }

      return "";
    };
  })

  ///////////////////////////////////////////////////////
  //                     Languages                     //
  ///////////////////////////////////////////////////////

  .controller("languagesCtrl", function($scope, $translate) {
    $scope.changeLanguage = function(language) {
      $translate.use(language);
    };

    $scope.getLanguage = function() {
      return $translate.use();
    };
  })

  ///////////////////////////////////////////////////////
  //                     Domain                        //
  ///////////////////////////////////////////////////////

  .factory("domainService", function($http) {
    return {
      queryDomain: function(fqdn) {
        return $http.get("/domain/" + fqdn + "/verification")
          .then(
            function(response) {
              return response;
            },
            function(response) {
              return response;
            });
      },
      retrieveDomain: function(uri) {
        return $http.get(uri)
          .then(
            function(response) {
              return response;
            },
            function(response) {
              return response;
            });
      },
      retrieveDomains: function(uri, etag) {
        return $http.get(uri, {
            headers: {
              "If-None-Match": etag
            }
        })
          .then(
            function(response) {
              return response;
            },
            function(response) {
              return response;
            });
      },
      verifyDomain: function(domain) {
        return $http.put("/domain/" + domain.fqdn + "/verification", domain, {
          headers: {
            "Content-Type": "application/json"
          }
        })
          .then(
            function(response) {
              return response;
            },
            function(response) {
              return response;
            });
      },
      removeDomain: function(uri, etag) {
        return $http.delete(uri, {
            headers: {
              "If-Match": etag
            }
        })
          .then(
            function(response) {
              return response;
            },
            function(response) {
              return response;
            });
      },
      saveDomain: function(domain, etag) {
        return $http.put("/domain/" + domain.fqdn, domain, {
          headers: {
            "Content-Type": "application/json",
            "If-Match": etag
          }
        })
          .then(
            function(response) {
              return response;
            },
            function(response) {
              return response;
            });
      }
    };
  })

  .directive("domain", function() {
    return {
      restrict: 'E',
      scope: {
        domain: '=',
        selectedDomains: '='
      },
      templateUrl: "/directives/domain.html",
      controller: function($scope, $translate, domainService) {
        $scope.freshDomain = emptyDomain;

        $scope.hasErrors = function(domain) {
          errors = false;

          if (domain.nameservers) {
            domain.nameservers.forEach(function(nameserver) {
              if (nameserver.lastStatus != "NOTCHECKED" && nameserver.lastStatus != "OK") {
                errors = true;
              }
            });
          }

          if (domain.dsset) {
            domain.dsset.forEach(function(ds) {
              if (ds.lastStatus != "NOTCHECKED" && ds.lastStatus != "OK") {
                errors = true;
              }
            });
          }

          return errors;
        };

        $scope.wasChecked = function(domain) {
          checked = false;

          if (domain.nameservers) {
            domain.nameservers.forEach(function(nameserver) {
              if (nameserver.lastStatus != "NOTCHECKED") {
                checked = true;
              }
            });
          }

          if (domain.dsset) {
            domain.dsset.forEach(function(ds) {
              if (ds.lastStatus != "NOTCHECKED") {
                checked = true;
              }
            });
          }

          return checked;
        };

        $scope.dateDefined = function(input) {
          var datetime = moment(input, ["YYYY-MM-DDTHH:mm:ss.SSSZ", "YYYY-MM-DDTHH:mm:ssZ"])
          if (datetime.isValid()) {
            return (datetime.unix() > 0);
          }
          return false;
        };

        $scope.getLanguage = function() {
          return $translate.use();
        };

        $scope.getAlgorithm = function(id) {
          for (var i = 0; i < algorithms.length; i++) {
            if (algorithms[i].id == id) {
              return algorithms[i].name;
            }
          }

          return id;
        }

        $scope.getDSDigestType = function(id) {
          for (var i = 0; i < dsDigestTypes.length; i++) {
            if (dsDigestTypes[i].id == id) {
              return dsDigestTypes[i].name;
            }
          }

          return id;
        }

        $scope.showDSDigest = function(digest) {
          if (digest == undefined) {
            return "";
          }

          if (digest.length >= 24) {
            return digest.substr(0, 12) + "..." + digest.substr(digest.length - 12, 12);
          }

          return digest;
        };

        $scope.toggleDetails = function(domain) {
          if (!$scope.details) {
            $scope.retrieveDomain(domain);
          }
          $scope.details = !$scope.details;
        };

        $scope.retrieveDomain = function(domain) {
          var uri = findLink(domain.links, "self");

          domainService.retrieveDomain(uri).then(
            function(response) {
              if (response.status == 200) {
                $scope.freshDomain = response.data;
                $scope.freshDomain.etag = response.headers.Etag;

                if ($scope.freshDomain.nameservers == undefined) {
                  $scope.freshDomain.nameservers = [];
                }

                if ($scope.freshDomain.dnskeys == undefined) {
                  $scope.freshDomain.dnskeys = [];
                }

                if ($scope.freshDomain.dsset == undefined) {
                  $scope.freshDomain.dsset = [];
                }

                if ($scope.freshDomain.owners == undefined) {
                  $scope.freshDomain.owners = [];
                }

              } else if (response.status == 400) {
                alertify.error(response.data.message);

              } else {
                $translate("Server error").then(function(translation) {
                  alertify.error(response.data.message);
                });
              }
            }
          );
        };

        $scope.verifyDomain = function(domain) {
          $scope.verifyWorking = true;

          // If the domain is imported without DS records, the DS set attribute will be undefined
          if (domain.dsset) {
            // Convert keytag to number. For now I don't want to use valueAsNumber function
            // because the code will remain the same size and older browsers will break
            for (var i = 0; i < domain.dsset.length; i += 1) {
              var keytag = parseInt(domain.dsset[i].keytag, 10);
              if (keytag == NaN) {
                keytag = 0;
              }
              domain.dsset[i].keytag = keytag;
            }
          }

          domainService.verifyDomain(domain).then(
            function(response) {
              if (response.status == 200) {
                $translate("Verify result").then(function(translation) {
                  alertify.success(translation);
                });

                $scope.verifyResult = response.data; // For unit tests only

                // Cannot replace the whole object because we lose the reference
                domain.nameservers = response.data.nameservers;
                domain.dsset = response.data.dsset;

              } else if (response.status == 400) {
                alertify.error(response.data.message);

              } else {
                $translate("Server error").then(function(translation) {
                  alertify.error(translation);
                });
              }

              $scope.verifyWorking = false;
            });
        };

        $scope.removeDomain = function(domain) {
          $scope.removeWorking = true;

          var uri = findLink(domain.links, "self");
          domainService.removeDomain(uri, domain.etag).then(
            function(response) {
              if (response.status == 204) {
                $translate("Domain removed").then(function(translation) {
                  alertify.success(translation);
                });
                $scope.success = true; // For unit tests only

              } else if (response.status == 400) {
                alertify.error(response.data.message);

              } else {
                $translate("Server error").then(function(translation) {
                  alertify.error(translation);
                });
              }

              $scope.removeWorking = false;
            });
        };

        $scope.selectDomain = function(fqdn) {
          if (fqdn in $scope.selectedDomains) {
            delete $scope.selectedDomains[fqdn];
          } else {
            $scope.selectedDomains[fqdn] = true;
          }
        };

        $scope.isDomainSelected = function(fqdn) {
          return (fqdn in $scope.selectedDomains);
        };
      }
    };
  })

  .directive("domainform", function() {
    return {
      restrict: 'E',
      scope: {
        domain: '=',
        formCtrl: '='
      },
      templateUrl: "/directives/domainform.html",
      controller: function($scope, $rootScope, $translate, $compile, domainService) {
        if ($scope.domain == emptyDomain) {
          $scope.domain = angular.copy($scope.domain);
        }

        $scope.emptyNameserver = angular.copy(emptyNameserver);
        $scope.emptyDNSKEY = angular.copy(emptyDNSKEY);
        $scope.emptyDS = angular.copy(emptyDS);
        $scope.emptyOwner = angular.copy(emptyOwner);

        // Don't need to copy constant values
        $scope.dnskeyFlags = dnskeyFlags;
        $scope.algorithms = algorithms;
        $scope.dsDigestTypes = dsDigestTypes;
        $scope.ownerLanguages = ownerLanguages;

        $scope.csv = {
          working: false,
          content: "",
          domainsToUpload: 0,
          domainsUploaded: 0,
          success: 0,
          errors: []
        };

        // Easy way to allow extenal scopes to call a function in the directive
        $scope.formCtrl = $scope.formCtrl || {};
        $scope.formCtrl.clear = function() {
          $scope.domain = angular.copy(emptyDomain);
          $scope.csv = {
            working: false,
            content: "",
            domainsToUpload: 0,
            domainsUploaded: 0,
            success: 0,
            errors: []
          };
        };

        $scope.needsGlue = function(fqdn, host) {
          if (fqdn == undefined || host == undefined ||
            fqdn.length == 0 || host.length < fqdn.length) {

            return false;
          }

          // Add final dot to FQDN and host to make a valid comparission

          if (fqdn.slice(-1) != ".") {
            fqdn += ".";
          }

          if (host.length > 0 && host.slice(-1) != ".") {
            host += ".";
          }

          return host.indexOf(fqdn, host.length - fqdn.length) !== -1;
        };

        $scope.addToList = function(object, list) {
          if (object == undefined) {
            return;
          }

          list.push(angular.copy(object));
        };

        $scope.removeFromList = function(index, list) {
          if (index >= 0 && index < list.length) {
            list.splice(index, 1);
          }
        };

        $scope.queryDomain = function(fqdn) {
          $scope.importWorking = true;

          domainService.queryDomain(fqdn).then(
            function(response) {
              if (response.status == 200) {
                $scope.domain = response.data;

                if ($scope.domain.nameservers == undefined) {
                  $scope.domain.nameservers = [];
                }
                if ($scope.domain.dsset == undefined) {
                  $scope.domain.dsset = [];
                }
                if ($scope.domain.dnskeys == undefined) {
                  $scope.domain.dnskeys = [];
                }
                if ($scope.domain.owners == undefined) {
                  $scope.domain.owners = [];
                }

              } else if (response.status == 400) {
                alertify.error(response.data.message);

              } else {
                $translate("Server error").then(function(translation) {
                  alertify.error(translation);
                });
              }

              $scope.importWorking = false;
            });
        };

        $scope.verifyDomain = function(domain) {
          $scope.verifyWorking = true;

          // If the domain is imported without DS records, the DS set attribute will be undefined
          if (domain.dsset) {
            // Convert keytag to number. For now I don't want to use valueAsNumber function
            // because the code will remain the same size and older browsers will break
            for (var i = 0; i < domain.dsset.length; i += 1) {
              var keytag = parseInt(domain.dsset[i].keytag, 10);
              if (keytag == NaN) {
                keytag = 0;
              }
              domain.dsset[i].keytag = keytag;
            }
          }

          domainService.verifyDomain(domain).then(
            function(response) {
              if (response.status == 200) {
                $translate("Verify result").then(function(translation) {
                  alertify.success(translation);
                });

                alertify.alert(verificationResponseToHTML(response.data));
                $scope.verifyResult = response.data; // For unit tests only

              } else if (response.status == 400) {
                alertify.error(response.data.message);

              } else {
                $translate("Server error").then(function(translation) {
                  alertify.error(translation);
                });
              }

              $scope.verifyWorking = false;
            });
        };

        $scope.saveDomain = function(domain) {
          $scope.saveWorking = true;

          // If the domain is imported without DS records, the DS set attribute will be undefined
          if (domain.dsset) {
            // Convert keytag to number. For now I don't want to use valueAsNumber function
            // because the code will remain the same size and older browsers will break
            for (var i = 0; i < domain.dsset.length; i += 1) {
              var keytag = parseInt(domain.dsset[i].keytag, 10);
              if (keytag == NaN) {
                keytag = 0;
              }
              domain.dsset[i].keytag = keytag;
            }
          }

          domainService.saveDomain(domain, domain.etag).then(
            function(response) {
              if (response.status == 201 || response.status == 204) {
                $translate("Domain created").then(function(translation) {
                  alertify.success(translation);
                });
                $rootScope.menu = "domains";
                $scope.success = true; // For unit tests only

              } else if (response.status == 400) {
                alertify.error(response.data.message);

              } else {
                $translate("Server error").then(function(translation) {
                  alertify.error(translation);
                });
              }

              $scope.saveWorking = false;
            });
        };

        $scope.storeCSVFile = function(content) {
          $scope.csv.content = content;
        };

        $scope.importCSV = function() {
          $scope.csv.working = true;
          $scope.csv.domainsUploaded = 0;
          $scope.csv.success = 0;
          $scope.csv.errors = [];

          var lines = $scope.csv.content.split("\n");
          $scope.csv.domainsToUpload = lines.length;

          lines.forEach(function(line, lineNumber) {
            // Start line number as one instead of zero
            lineNumber += 1;

            if (line.length == 0) {
              $scope.csv.domainsUploaded += 1;
              if ($scope.csv.domainsUploaded == $scope.csv.domainsToUpload) {
                $scope.csv.working = false;
              }
              return;
            }

            // Parse first level using comma:
            // fqdn,ns1,ns2,ns3,ns4,ns5,ns6,ds1,ds2,owner1,owner2,owner3
            var fields = line.split(",");

            if (fields.length != 12) {
              $translate("CSV fields count error").then(function(translation) {
                $scope.csv.errors.push({
                  message: translation,
                  lineNumber: lineNumber
                });
              });

              $scope.csv.domainsUploaded += 1;
              if ($scope.csv.domainsUploaded == $scope.csv.domainsToUpload) {
                $scope.csv.working = false;
              }
              return;
            }

            var domain = {};
            domain.fqdn = fields[0];
            domain.nameservers = [];
            domain.dsset = [];
            domain.owners = [];

            // Loop between the nameservers fields
            for (var i = 1; i <= 6; i++) {
              if (fields[i].length == 0) {
                continue;
              }

              // Parse the nameserver:
              // host$ipv4$ipv6
              var nsParts = fields[i].split("$");

              if (nsParts.length > 3) {
                $translate("CSV NS error").then(function(translation) {
                  $scope.csv.errors.push({
                    message: translation,
                    lineNumber: lineNumber
                  });
                });
                continue;
              }

              var nameserver = {
                host: nsParts[0]
              };

              if (nsParts.length > 1) {
                nameserver.ipv4 = nsParts[1];
              }

              if (nsParts.length > 2) {
                nameserver.ipv4 = nsParts[2];
              }

              domain.nameservers.push(nameserver);
            }

            // Loop between the DS set fields
            for (var i = 7; i <= 8; i++) {
              if (fields[i].length == 0) {
                continue;
              }

              // Parse the DS:
              // keytag$algorithm$digestType$digest
              var dsParts = fields[i].split("$");

              if (dsParts.length != 4) {
                $translate("CSV DS error").then(function(translation) {
                  $scope.csv.errors.push({
                    message: translation,
                    lineNumber: lineNumber
                  });
                });
                continue;
              }

              var ds = {
                keytag: parseInt(dsParts[0], 10),
                algorithm: parseInt(dsParts[1], 10),
                digestType: parseInt(dsParts[2], 10),
                digest: dsParts[3]
              };

              if (ds.keytag == NaN) {
                ds.keytag = 0;
              }

              if (ds.algorithm == NaN) {
                ds.algorithm = 0;
              }

              if (ds.digestType == NaN) {
                ds.digestType = 0;
              }

              domain.dsset.push(ds);
            }

            // Loop between the owner fields
            for (var i = 9; i <= 11; i++) {
              if (fields[i].length == 0) {
                continue;
              }

              // Parse the owner:
              // email$language
              var ownerParts = fields[i].split("$");

              if (ownerParts.length != 2) {
                $translate("CSV owener error").then(function(translation) {
                  $scope.csv.errors.push({
                    message: translation,
                    lineNumber: lineNumber
                  });
                });
                continue;
              }

              domain.owners.push({
                email: ownerParts[0],
                language: ownerParts[1]
              });
            }

            domainService.saveDomain(domain, domain.etag).then(
              function(response) {
                if (response.status == 201 || response.status == 204) {
                  $scope.csv.success += 1;
                  $scope.success = true; // For unit tests only

                } else if (response.status == 400) {
                  $scope.csv.errors.push({
                    message: response.data.message,
                    lineNumber: lineNumber
                  });

                } else {
                  $translate("Server error").then(function(translation) {
                    $scope.csv.errors.push({
                      message: translation,
                      lineNumber: lineNumber
                    });
                  });
                }

                $scope.csv.domainsUploaded += 1;
                if ($scope.csv.domainsUploaded == $scope.csv.domainsToUpload) {
                  $scope.csv.working = false;
                }
              });
          });
        };
      }
    };
  })

  .controller("domainCtrl", function($scope) {
    $scope.emptyDomain = angular.copy(emptyDomain);
  })

  .controller("domainsCtrl", function($scope, $translate, $timeout, $anchorScroll, $location, domainService) {
    $scope.pageSizes = [ 20, 40, 60, 80, 100 ];
    $scope.lastRetrieveDomains = moment();
    $scope.selectedDomains = {};

    $scope.countSelectedDomains = function() {
      return Object.keys($scope.selectedDomains).length;
    };

    $scope.retrieveDomains = function(page, pageSize, filter, successFunction) {
      var uri = "";

      if (page != undefined) {
        if (uri.length == 0) {
          uri += "?";
        } else {
          uri += "&";
        }
        uri += "page=" + page;
      }

      if (pageSize != undefined) {
        if (uri.length == 0) {
          uri += "?";
        } else {
          uri += "&";
        }
        uri += "pagesize=" + pageSize;
      }

      if (filter != undefined) {
        if (uri.length == 0) {
          uri += "?";
        } else {
          uri += "&";
        }
        uri += "filter=" + filter;
      }

      uri = "/domains/" + uri;

      $scope.retrieveDomainsByURI(uri, successFunction);
    };

    $scope.retrieveDomainsByURI = function(uri, successFunction) {
      $scope.retrieveDomainsURI = uri;
      $scope.lastRetrieveDomains = moment();
      
      domainService.retrieveDomains($scope.retrieveDomainsURI, $scope.etag).then(
        function(response) {
          $scope.processDomainsResult(response, successFunction);
        });
    };

    $scope.processDomainsResult = function(response, successFunction) {
      if (response.status == 200) {
        $scope.etag = response.headers.Etag;

        if (!$scope.pagination) {
          $scope.pagination = response.data;

        } else {
          $scope.pagination.page = response.data.page;
          $scope.pagination.numberOfItems = response.data.numberOfItems;
          $scope.pagination.numberOfPages = response.data.numberOfPages;
          $scope.pagination.pageSize = response.data.pageSize;
          $scope.pagination.links = response.data.links;
          $scope.pagination.domains = response.data.domains;
        }

        if (successFunction != undefined) {
          successFunction();
        }

      } else if (response.status == 304) {
        // Not modified

      } else if (response.status == 400) {
        alertify.error(response.data.message);

      } else {
        $translate("Server error").then(function(translation) {
          alertify.error(translation);
        });
      }
    };

    $scope.retrieveDomainsWorker = function() {
      if (!$scope.retrieveDomainsURI) {
        $scope.retrieveDomains();
      }

      var now = moment();
      if (now.subtract($scope.lastRetrieveDomains).seconds() >= refreshRate) {
        domainService.retrieveDomains($scope.retrieveDomainsURI, $scope.etag).then(
          function(response) {
            $scope.processDomainsResult(response);
          });
          $scope.lastRetrieveDomains = moment();
      }

      $timeout($scope.retrieveDomainsWorker, refreshRate * 1000);
    };

    $scope.findLink = function(pagination, type) {
      if (pagination == undefined || type == undefined) {
        return "";
      }

      return findLink(pagination.links, type);
    };

    $scope.scrollToPagination = function() {
      // We can't do this to fast, because the HTML page is not render yet
      $timeout(function() {
        $location.hash("domainPagination");
        $anchorScroll();
        $location.hash("");
      }, 100);
    };

    $scope.retrieveDomainsWorker();
  })

  // Directive retrieved from:
  // http://veamospues.wordpress.com/2014/01/27/reading-files-with-angularjs/
  .directive('onReadFile', function ($parse) {
    return {
      restrict: 'A',
      scope: false,
      link: function(scope, element, attrs) {
        var fn = $parse(attrs.onReadFile);
              
        element.on('change', function(onChangeEvent) {
          var reader = new FileReader();
                  
          reader.onload = function(onLoadEvent) {
            scope.$apply(function() {
              fn(scope, {$fileContent:onLoadEvent.target.result});
            });
          };

          reader.readAsText((onChangeEvent.srcElement || onChangeEvent.target).files[0]);
        });
      }
    };
  })


  ///////////////////////////////////////////////////////
  //                     Scan                          //
  ///////////////////////////////////////////////////////

  .factory("scanService", function($http) {
    return {
      retrieve: function(uri, etag) {
        return $http.get(uri, {
          headers: {
            "If-None-Match": etag
          }
        })
          .then(
            function(response) {
              return response;
            },
            function(response) {
              return response;
            });
      }
    };
  })

  .directive("scan", function() {
    return {
      restrict: 'E',
      scope: {
        scan: '='
      },
      templateUrl: "/directives/scan.html",
      controller: function($scope, $translate, scanService) {
        $scope.freshDomain = emptyScan;

        $scope.countStatistics = function(statistics) {
          var counter = 0;
          for (var status in statistics) {
            counter += statistics[status];
          }
          return counter;
        };

        $scope.getLanguage = function() {
          return $translate.use();
        };

        $scope.toggleDetails = function(scan) {
          if (!$scope.details) {
            $scope.retrieveScan(scan);
          }
          $scope.details = !$scope.details;
        };

        $scope.retrieveScan = function(scan) {
          var uri = findLink(scan.links, "self");

          scanService.retrieve(uri).then(
            function(response) {
              if (response.status == 200) {
                $scope.freshScan = response.data;
                $scope.freshScan.etag = response.headers.Etag;

              } else if (response.status == 400) {
                alertify.error(response.data.message);

              } else {
                $translate("Server error").then(function(translation) {
                  alertify.error(translation);
                });
              }
            }
          );
        };
      }
    };
  })

  .controller("scanCtrl", function($scope, $translate, $timeout, $anchorScroll, $location, scanService) {
    $scope.pageSizes = [ 20, 40, 60, 80, 100 ];
    $scope.lastRetrieveScans = moment();

    $scope.getLanguage = function() {
      return $translate.use();
    };

    $scope.retrieveScans = function(page, pageSize, successFunction) {
      var uri = "";

      if (page != undefined) {
        if (uri.length == 0) {
          uri += "?";
        } else {
          uri += "&";
        }
        uri += "page=" + page;
      }

      if (pageSize != undefined) {
        if (uri.length == 0) {
          uri += "?";
        } else {
          uri += "&";
        }
        uri += "pagesize=" + pageSize;
      }

      uri = "/scans/" + uri;

      $scope.retrieveScansByURI(uri, successFunction);
    };

    $scope.retrieveScansByURI = function(uri, successFunction) {
      $scope.retrieveScansURI = uri;
      $scope.lastRetrieveScans = moment();
      
      scanService.retrieve($scope.retrieveScansURI, $scope.etag).then(
        function(response) {
          $scope.processScansResult(response, successFunction);
        });
    };

    $scope.processScansResult = function(response, successFunction) {
      if (response.status == 200) {
        $scope.etag = response.headers.Etag;

        if (!$scope.pagination) {
          $scope.pagination = response.data;

        } else {
          $scope.pagination.page = response.data.page;
          $scope.pagination.numberOfItems = response.data.numberOfItems;
          $scope.pagination.numberOfPages = response.data.numberOfPages;
          $scope.pagination.pageSize = response.data.pageSize;
          $scope.pagination.links = response.data.links;
          $scope.pagination.scans = response.data.scans;
        }

        $scope.currentScanURI = findLink($scope.pagination.links, "current");

        if (successFunction != undefined) {
          successFunction();
        }

      } else if (response.status == 304) {
        // Not modified

      } else if (response.status == 400) {
        alertify.error(response.data.message);

      } else {
        $translate("Server error").then(function(translation) {
          alertify.error(translation);
        });
      }
    };

    $scope.retrieveCurrentScan = function() {
      // We dind't got the current scan URI yet, wait until we do
      if (!$scope.currentScanURI) {
        $timeout($scope.retrieveCurrentScan, refreshRate*1000);
        return;
      }

      scanService.retrieve($scope.currentScanURI, $scope.currentScanEtag).then(
        function(response) {
          if (response.status == 200) {
            $scope.currentScanEtag = response.headers.Etag;
            $scope.currentScan = response.data.scans[0];

          } else if (response.status == 304) {
            // Not modified

          } else if (response.status == 400) {
            alertify.error(response.data.message);

          } else {
            $translate("Server error").then(function(translation) {
              alertify.error(translation);
            });
          }

          $timeout($scope.retrieveCurrentScan, refreshRate*1000);
        });
    };

    $scope.retrieveScansWorker = function() {
      if (!$scope.retrieveScansURI) {
        $scope.retrieveScans();
      }

      var now = moment();
      if (now.subtract($scope.lastRetrieveScans).seconds() >= refreshRate) {
        scanService.retrieve($scope.retrieveScansURI, $scope.etag).then(
          function(response) {
            $scope.processScansResult(response);
          });
        $scope.lastRetrieveScans = moment();
      }

      $timeout($scope.retrieveScansWorker, refreshRate*1000);
    };

    $scope.findLink = function(pagination, type) {
      if (pagination == undefined || type == undefined) {
        return "";
      }

      return findLink(pagination.links, type);
    };

    $scope.scrollToPagination = function() {
      // We can't do this to fast, because the HTML page is not render yet
      $timeout(function() {
        $location.hash("scanPagination");
        $anchorScroll();
        $location.hash("");
      }, 100);
    };

    $scope.retrieveCurrentScan();
    $scope.retrieveScansWorker();
  });