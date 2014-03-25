// References:
// http://stackoverflow.com/questions/17646034/what-is-the-best-practice-for-making-an-ajax-call-in-angular-js

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
  "pt-BR"
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

// Apply source list to destination list, adding new elements, removing old ones and
// keeping the ones that are equal
function mergeList(source, destination, areEqual, mergeObject) {
  if (!source || !destination) {
    return;
  }

  // Check new items
  for (var i = 0; i < source.length; i++) {
    var found = false;
    for (var j = 0; j < destination.length; j++) {
      if (areEqual(source[i], destination[j])) {
        if (mergeObject) {
          mergeObject(source[i], destination[j]);
        }

        found = true;
        break;
      }
    }

    if (!found) {
      destination.push(source[i]);
    }
  }

  // Check removed items
  for (var j = 0; j < destination.length; j++) {
    var found = false;
    for (var i = 0; i < source.length; i++) {
      if (areEqual(source[i], destination[j])) {
        found = true;
        break;
      }
    }

    if (!found) {
      destination.splice(j, 1);
    }
  }
}

angular.module("shelter", ["ngAnimate", "ngCookies", "pascalprecht.translate"])

  .config(function($translateProvider, $httpProvider) {
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
      var datetime = moment(input, ["YYYY-MM-DDTHH:mm:ss.SSSZ", "YYYY-MM-DDTHH:mm:ssZ"])
      if (datetime.isValid() && language) {
        // Detect empty datetime
        if (datetime.unix() <= 0) {
          return ""
        }

        language = language.replace("_", "-");
        moment.lang(language);
        return datetime.format("MMMM Do YYYY, h:mm:ss a");
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
      retrieveDomains: function(page, pageSize) {
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

        uri = "/domains/" + uri;

        return $http.get(uri)
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
      removeDomain: function(fqdn) {
        return $http.delete("/domain/" + fqdn)
          .then(
            function(response) {
              return response;
            },
            function(response) {
              return response;
            });
      },
      saveDomain: function(domain) {
        return $http.put("/domain/" + domain.fqdn, domain, {
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
      }
    };
  })

  .directive("domain", function() {
    return {
      restrict: 'E',
      scope: {
        domain: '='
      },
      templateUrl: "/directives/domain.html",
      controller: function($scope, $translate, domainService) {
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
                $scope.error = null;
                $translate("Verify result").then(function(translation) {
                  $scope.success = translation;
                  $scope.verifyResult = response.data;
                });

              } else if (response.status == 400) {
                $scope.success = null;
                $scope.verifyResult = null;
                $scope.error = response.data.message;

              } else {
                $scope.success = null;
                $scope.verifyResult = null;
                $translate("Server error").then(function(translation) {
                  $scope.error = translation;
                });
              }

              $scope.verifyWorking = false;
            });
        };

        $scope.removeDomain = function(fqdn) {
          $scope.removeWorking = true;

          domainService.removeDomain(fqdn).then(
            function(response) {
              if (response.status == 204) {
                $scope.error = null;
                $translate("Domain removed").then(function(translation) {
                  $scope.success = translation;
                  $scope.verifyResult = null;
                });

              } else if (response.status == 400) {
                $scope.success = null;
                $scope.verifyResult = null;
                $scope.error = response.data.message;

              } else {
                $scope.success = null;
                $scope.verifyResult = null;
                $translate("Server error").then(function(translation) {
                  $scope.error = translation;
                });
              }

              $scope.removeWorking = false;
            });
        }
      }
    };
  })

  .directive("domainform", function() {
    return {
      restrict: 'E',
      scope: {
        domain: '='
      },
      templateUrl: "/directives/domainform.html",
      controller: function($scope, $translate, domainService) {
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

        $scope.needsGlue = function(fqdn, host) {
          if (fqdn == undefined || host == undefined ||
            fqdn.length == 0 || host.length < fqdn.length) {

            return false;
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
                $scope.error = null;
                $scope.domain = response.data;

              } else if (response.status == 400) {
                $scope.error = null;
                $scope.error = response.data.message;

              } else {
                $scope.success = null;
                $translate("Server error").then(function(translation) {
                  $scope.error = translation;
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
                $scope.error = null;
                $translate("Verify result").then(function(translation) {
                  $scope.success = translation;
                  $scope.verifyResult = response.data;
                });

              } else if (response.status == 400) {
                $scope.error = null;
                $scope.error = response.data.message;

              } else {
                $scope.success = null;
                $translate("Server error").then(function(translation) {
                  $scope.error = translation;
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

          domainService.saveDomain(domain).then(
            function(response) {
              if (response.status == 201 || response.status == 204) {
                $translate("Domain created").then(function(translation) {
                  $scope.error = null;
                  $scope.success = translation;
                });

              } else if (response.status == 400) {
                $scope.error = null;
                $scope.error = response.data.message;

              } else {
                $scope.success = null;
                $translate("Server error").then(function(translation) {
                  $scope.error = translation;
                });
              }

              $scope.saveWorking = false;
            });
        };
      }
    };
  })

  .controller("domainCtrl", function($scope) {
    $scope.emptyDomain = angular.copy(emptyDomain);
  })

  .controller("domainsCtrl", function($scope, $translate, $timeout, domainService) {
    $scope.pageSizes = [ 20, 40, 60, 80, 100 ];

    $scope.retrieveDomains = function(page, pageSize) {
      domainService.retrieveDomains(page, pageSize).then(
        function(response) {
          if (response.status == 200) {
            $scope.error = null;

            if (!$scope.pagination) {
              $scope.pagination = response.data;

            } else {
              $scope.pagination.numberOfItems = response.data.numberOfItems;
              $scope.pagination.numberOfPages = response.data.numberOfPages;
              $scope.pagination.pageSize = response.data.pageSize;

              mergeList(response.data.domains,
                $scope.pagination.domains,
                function(networkDomain, domain) {
                  return networkDomain.fqdn == domain.fqdn;
                },
                function(networkDomain, domain) {
                  domain.nameservers = networkDomain.nameservers;
                  domain.dsset = networkDomain.dsset;
                  domain.owners = networkDomain.owners;
                });
            }

          } else if (response.status == 400) {
            $scope.error = null;
            $scope.error = response.data.message;

          } else {
            $scope.success = null;
            $translate("Server error").then(function(translation) {
              $scope.error = translation;
            });
          }

          $timeout(function() {
            if ($scope.pagination) {
              $scope.retrieveDomains($scope.pagination.page, $scope.pagination.pageSize);
            } else {
              $scope.retrieveDomains();
            }
          }, 5000);
        });
    };

    $scope.retrieveDomains();
  })



  ///////////////////////////////////////////////////////
  //                     Scan                          //
  ///////////////////////////////////////////////////////

  .factory("scanService", function($http) {
    return {
      retrieveScans: function(page, pageSize) {
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

        return $http.get(uri)
          .then(
            function(response) {
              return response;
            },
            function(response) {
              return response;
            });
      },
      retrieveCurrentScan: function() {
        return $http.get("/scan/current")
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
      controller: function($scope, $translate) {
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
      }
    };
  })

  .controller("scanCtrl", function($scope, $translate, $timeout, scanService) {
    $scope.pageSizes = [ 20, 40, 60, 80, 100 ];

    $scope.getLanguage = function() {
      return $translate.use();
    };

    $scope.retrieveScans = function(page, pageSize) {
      scanService.retrieveScans(page, pageSize).then(
        function(response) {
          if (response.status == 200) {
            $scope.error = null;

            if (!$scope.pagination) {
              $scope.pagination = response.data;

            } else {
              $scope.pagination.numberOfItems = response.data.numberOfItems;
              $scope.pagination.numberOfPages = response.data.numberOfPages;
              $scope.pagination.pageSize = response.data.pageSize;

              mergeList(response.data.scans,
                $scope.pagination.scans,
                function(networkScan, scan) {
                  return networkScan.startedAt == scan.startedAt;
                },
                function(networkScan, scan) {
                  scan.status = networkScan.status;
                  scan.finishedAt = networkScan.finishedAt;
                  scan.domainsScanned = networkScan.domainsScanned;
                  scan.domainsWithDNSSECScanned = networkScan.domainsWithDNSSECScanned;
                  scan.nameserverStatistics = networkScan.nameserverStatistics;
                  scan.dsStatistics = networkScan.dsStatistics;
                });
            }

          } else if (response.status == 400) {
            $scope.error = null;
            $scope.error = response.data.message;

          } else {
            $scope.success = null;
            $translate("Server error").then(function(translation) {
              $scope.error = translation;
            });
          }

          $timeout(function() {
            if ($scope.pagination) {
              $scope.retrieveScans($scope.pagination.page, $scope.pagination.pageSize);
            } else {
              $scope.retrieveScans();
            }
          }, 5000);
        });
    };

    $scope.retrieveCurrentScan = function() {
      scanService.retrieveCurrentScan().then(
        function(response) {
          if (response.status == 200) {
            $scope.error = null;
            $scope.currentScan = response.data;

          } else if (response.status == 400) {
            $scope.error = null;
            $scope.error = response.data.message;

          } else {
            $scope.success = null;
            $translate("Server error").then(function(translation) {
              $scope.error = translation;
            });
          }

          $timeout($scope.retrieveCurrentScan, 1000);
        });
    };

    $scope.retrieveCurrentScan();
    $scope.retrieveScans();
  });