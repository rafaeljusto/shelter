// References:
// http://stackoverflow.com/questions/17646034/what-is-the-best-practice-for-making-an-ajax-call-in-angular-js

angular.module("shelter", [])

  ///////////////////////////////////////////////////////
  //                     Domain                        //
  ///////////////////////////////////////////////////////

  .factory("domainService", function($http) {
    return {
      retrieveDomains: function() {
        return $http.get("/domains")
          .then(
            function(response) {
              return response.data;
            },
            function(error) {
              return error.data;
            });
      },
      saveDomain: function(domain) {
        return $http.put("/domain/" + domain.FQDN)
          .then(
            function(response) {
              return response.data;
            },
            function(error) {
              return error.data;
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
      controller: function($scope, domainService) {
        $scope.saveDomain = function(domain) {
          domainService.saveDomain(domain).then(
            function(response) {
              // TODO
            },
            function(error) {
              // TODO
            });
        };
      }
    };
  })

  .controller("domainCtrl", function($scope, $timeout, domainService) {
    $scope.retrieveDomains = function() {
      domainService.retrieveDomains().then(
        function(response) {
          $scope.pagination = response;
          $timeout($scope.retrieveDomains, 5000);
        },
        function(error) {
          // TODO
          $timeout($scope.retrieveDomains, 5000);
        });
    };

    $scope.retrieveDomains();
  })



  ///////////////////////////////////////////////////////
  //                     Scan                          //
  ///////////////////////////////////////////////////////

  .factory("scanService", function($http) {
    return {
      retrieveScans: function() {
        return $http.get("/scans")
          .then(
            function(response) {
              return response.data;
            },
            function(error) {
              return error.data;
            });
      },
      retrieveCurrentScan: function() {
        return $http.get("/currentscan")
          .then(
            function(response) {
              return response.data;
            },
            function(error) {
              return error.data;
            });
      }
    };
  })

  .directive("scan", function() {
    return {
      restrict: 'E',
      scope: {
        scan: '='
      }
    };
  })

  .controller("scanCtrl", function($scope, $timeout, scanService) {
    $scope.retrieveScans = function() {
      scanService.retrieveScans().then(
        function(data) {
          // TODO
          $timeout($scope.retrieveScans, 5000);
        },
        function(error) {
          // TODO
          $timeout($scope.retrieveScans, 5000);
        });
    };

    $scope.retrieveCurrentScan = function() {
      scanService.retrieveCurrentScan().then(
        function(data) {
          // TODO
        },
        function(error) {
          // TODO
        });
    };

    $scope.retrieveScans();
  });
