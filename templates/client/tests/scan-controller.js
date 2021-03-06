/*
 * Copyright 2014 Rafael Dantas Justo. All rights reserved.
 * Use of this source code is governed by a GPL
 * license that can be found in the LICENSE file.
 */

describe("Scan controller", function() {
  var scope, ctrl;

  beforeEach(module('shelter'));

  beforeEach(inject(function($rootScope, $controller, $injector) {
    localStorage.clear();

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET(/languages\/.+\.json/).respond(200, "{}");
    $httpBackend.flush();

    scope = $rootScope.$new();
    ctrl = $controller("scanCtrl", {
      $scope: scope
    });
  }));

  it("should verify if the get language function returns the default language", inject(function($translate) {
    expect(scope.getLanguage()).toBe($translate.preferredLanguage());
    expect(scope.getLanguage()).not.toBe("");
    expect(scope.getLanguage()).not.toBe(undefined);
  }));

  it("should return the retrieved scans properly", inject(function($injector) {
    var result = {
      numberOfItems: 1,
      numberOfPages: 1,
      pageSize: 20,
      scans: [
        {
          status: "EXECUTED",
          startedAt: "2014-03-26T08:30:00-03:00",
          finishedAt: "2014-03-26T09:15:00-03:00",
          domainsScanned: 150,
          domainsWithDNSSECScanned: 20,
          nameserverStatistics: {
            OK: 110,
            TIMEOUT: 40
          },
          dsStatistics: {
            OK: 15,
            EXPSIG: 5
          }
        }
      ]
    };

    var current = {
      numberOfItems: 1,
      numberOfPages: 1,
      pageSize: 20,
      scans: [
        {
          status: "WAITINGEXECUTION",
          scheduledAt: "2014-03-27T06:00:00-03:00"
        }
      ]
    };

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET("/scans/").respond(200, result);
    $httpBackend.whenGET("/scans/?page=1&pagesize=20").respond(200, result);
    $httpBackend.whenGET("/scans/?current").respond(200, current);

    scope.retrieveScans(1, 20);
    $httpBackend.flush();

    expect(scope.pagination.numberOfItems).toBe(1);
    expect(scope.pagination.numberOfPages).toBe(1);
    expect(scope.pagination.pageSize).toBe(20);
    expect(scope.pagination.scans.length).toBe(1);
    expect(scope.pagination.scans).toEqual(result.scans);
  }));

  it("should return the retrieved scans by URI properly", inject(function($injector) {
    var result = {
      numberOfItems: 1,
      numberOfPages: 1,
      pageSize: 20,
      scans: [
        {
          status: "EXECUTED",
          startedAt: "2014-03-26T08:30:00-03:00",
          finishedAt: "2014-03-26T09:15:00-03:00",
          domainsScanned: 150,
          domainsWithDNSSECScanned: 20,
          nameserverStatistics: {
            OK: 110,
            TIMEOUT: 40
          },
          dsStatistics: {
            OK: 15,
            EXPSIG: 5
          }
        }
      ]
    };

    var current = {
      numberOfItems: 1,
      numberOfPages: 1,
      pageSize: 20,
      scans: [
        {
          status: "WAITINGEXECUTION",
          scheduledAt: "2014-03-27T06:00:00-03:00"
        }
      ]
    };

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET("/scans/").respond(200, result);
    $httpBackend.whenGET("/scans/?page=1&pagesize=20").respond(200, result);
    $httpBackend.whenGET("/scans/?current").respond(200, current);

    scope.retrieveScansByURI("/scans/?page=1&pagesize=20");
    $httpBackend.flush();

    expect(scope.pagination.numberOfItems).toBe(1);
    expect(scope.pagination.numberOfPages).toBe(1);
    expect(scope.pagination.pageSize).toBe(20);
    expect(scope.pagination.scans.length).toBe(1);
    expect(scope.pagination.scans).toEqual(result.scans);
  }));

  it("should return the current scan properly", inject(function($injector, $timeout) {
    var result = {
      numberOfItems: 1,
      numberOfPages: 1,
      pageSize: 20,
      scans: [
        {
          status: "EXECUTED",
          startedAt: "2014-03-26T08:30:00-03:00",
          finishedAt: "2014-03-26T09:15:00-03:00",
          domainsScanned: 150,
          domainsWithDNSSECScanned: 20,
          nameserverStatistics: {
            OK: 110,
            TIMEOUT: 40
          },
          dsStatistics: {
            OK: 15,
            EXPSIG: 5
          }
        }
      ]
    };

    var current = {
      numberOfItems: 1,
      numberOfPages: 1,
      pageSize: 20,
      scans: [
        {
          status: "WAITINGEXECUTION",
          scheduledAt: "2014-03-27T06:00:00-03:00"
        }
      ]
    };

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET("/scans/").respond(200, result);
    $httpBackend.whenGET("/scans/?current").respond(200, current);

    scope.currentScanURI = "/scans/?current";
    $timeout.flush();
    $httpBackend.flush();

    expect(scope.currentScan).toEqual(current.scans[0]);  
    
  }));
});