// References:
// http://stackoverflow.com/questions/17646034/what-is-the-best-practice-for-making-an-ajax-call-in-angular-js

angular.module("shelter", ["ngCookies", "pascalprecht.translate"])

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
      var datetime = moment("YYYY-MM-DDTHH:mm:ssZ", input)
      if (datetime.isValid()) {
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
  })

  ///////////////////////////////////////////////////////
  //                     Domain                        //
  ///////////////////////////////////////////////////////

  .factory("domainService", function($http) {
    return {
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
        $scope.getLanguage = function() {
          return $translate.use();
        };

        $scope.saveDomain = function(domain) {
          domainService.saveDomain(domain).then(
            function(response) {
              // TODO
            });
        };
      }
    };
  })

  .controller("domainCtrl", function($scope, $timeout, $translate, domainService) {
    $scope.dnskeyFlags = [
      {id: 256, name:"ZSK"},
      {id: 257, name:"KSK"},
    ];

    $scope.algorithms = [
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

    $scope.dsDigestTypes = [
      {id: 1, name: "SHA1"},
      {id: 2, name: "SHA256"},
      {id: 3, name: "GOST94"},
      {id: 4, name: "SHA384"},
      {id: 5, name: "SHA512"},
    ];

    $scope.ownerLanguages = [
      "en-US",
      "pt-BR"
    ];

    $scope.emptyNameserver = {
      host: "",
      ipv4: "",
      ipv6: ""
    };

    $scope.emptyDNSKEY = {
      flags: 257,
      algorithm: 8,
      publickKey: ""
    };

    $scope.emptyDS = {
      keytag: "",
      algorithm: 8,
      digestType: 2,
      digest: ""
    };

    $scope.emptyOwner = {
      email: "",
      language: "en-US"
    };

    $scope.domain = {
      fqdn: "",
      nameservers: [
        angular.copy($scope.emptyNameserver),
        angular.copy($scope.emptyNameserver)
      ],
      dnskeys: [],
      dsset: [
        angular.copy($scope.emptyDS)
      ],
      owners: [
        angular.copy($scope.emptyOwner)
      ]
    };

    $scope.addToList = function(object, list) {
      list.push(angular.copy(object));
    };

    $scope.removeFromList = function(index, list) {
      if (index >= 0 && index < list.length) {
        list.splice(index, 1);
      }
    };

    $scope.saveDomain = function() {
      // Convert keytag to number. For now I don't want to use valueAsNumber function
      // because the code will remain the same size and older browsers will break
      for (var i = 0; i < $scope.domain.dsset.length; i += 1) {
        var keytag = parseInt($scope.domain.dsset[i].keytag, 10);
        if (keytag == NaN) {
          keytag = 0;
        }
        $scope.domain.dsset[i].keytag = keytag;
      }

      domainService.saveDomain($scope.domain).then(
        function(response) {
          if (response.status == 400) {
            $scope.error = response.data.message;
          } else if (response.status != 201 && response.status != 204) {
            $scope.error = $translate("Server error");
          }
        });
    };
  })

  .controller("domainsCtrl", function($scope, $timeout, domainService) {
    $scope.pageSizes = [ 20, 40, 60, 80, 100 ];

    $scope.retrieveDomains = function(page, pageSize) {
      domainService.retrieveDomains(page, pageSize).then(
        function(response) {
          if (response.status == 400) {
            $scope.error = response.data.message;
          } else {
           $scope.pagination = response.data;
          }
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
        return $http.get("/currentscan")
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
        $scope.getLanguage = function() {
          return $translate.use();
        };
      }
    };
  })

  .controller("scanCtrl", function($scope, $timeout, scanService) {
    $scope.pageSizes = [ 20, 40, 60, 80, 100 ];

    $scope.retrieveScans = function(page, pageSize) {
      scanService.retrieveScans(page, pageSize).then(
        function(response) {
          if (response.status == 200) {
            $scope.pagination = response.data;
          } else {
            // TODO
          }
        });
    };

    $scope.retrieveCurrentScan = function() {
      scanService.retrieveCurrentScan().then(
        function(response) {
          // TODO
        });
    };

    $scope.retrieveScans();
  });
