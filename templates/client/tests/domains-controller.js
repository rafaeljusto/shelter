/*
 * Copyright 2014 Rafael Dantas Justo. All rights reserved.
 * Use of this source code is governed by a GPL
 * license that can be found in the LICENSE file.
 */

describe("Domains controller", function() {
  var scope, ctrl;

  beforeEach(module('shelter'));

  beforeEach(inject(function($rootScope, $controller, $injector) {
    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET(/languages\/.+\.json/).respond(200, "{}");
    $httpBackend.flush();

    scope = $rootScope.$new();
    ctrl = $controller("domainsCtrl", {
      $scope: scope
    });
  }));

  it("should retrieve the domains", inject(function($injector) {
    var result = {
      numberOfItems: 1,
      numberOfPages: 1,
      pageSize: 20,
      domains: [
        {
          fqdn: "br.",
          nameservers: [
            {
              host: "a.dns.br",
              ipv4: "200.160.0.10",
              lastStatus: "OK",
              lastOKAt: "2014-03-25T11:00:00-03:00",
              lastCheckAt: "2014-03-25T11:00:00-03:00"
            }
          ],
          dsset: [
            {
              keytag: 41674,
              algorithm: 5,
              digestType: 1,
              digest: "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
              lastStatus: "OK",
              lastOKAt: "2014-03-25T11:00:00-03:00",
              lastCheckAt: "2014-03-25T11:00:00-03:00"
            }
          ]
        }
      ]
    };

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET("/domains/?expand").respond(200, result);
    $httpBackend.whenGET("/domains/?expand&page=1&pagesize=20").respond(200, result);

    scope.retrieveDomains(1, 20);
    $httpBackend.flush();

    expect(scope.pagination.numberOfItems).toBe(1);
    expect(scope.pagination.numberOfPages).toBe(1);
    expect(scope.pagination.pageSize).toBe(20);
    expect(scope.pagination.domains.length).toBe(1);
    expect(scope.pagination.domains).toEqual(result.domains);
  }));

  it("should retrieve the domains by URI", inject(function($injector) {
    var result = {
      numberOfItems: 1,
      numberOfPages: 1,
      pageSize: 20,
      domains: [
        {
          fqdn: "br.",
          nameservers: [
            {
              host: "a.dns.br",
              ipv4: "200.160.0.10",
              lastStatus: "OK",
              lastOKAt: "2014-03-25T11:00:00-03:00",
              lastCheckAt: "2014-03-25T11:00:00-03:00"
            }
          ],
          dsset: [
            {
              keytag: 41674,
              algorithm: 5,
              digestType: 1,
              digest: "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
              lastStatus: "OK",
              lastOKAt: "2014-03-25T11:00:00-03:00",
              lastCheckAt: "2014-03-25T11:00:00-03:00"
            }
          ]
        }
      ]
    };

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET("/domains/?expand").respond(200, result);
    $httpBackend.whenGET("/domains/?expand&page=1&pagesize=20").respond(200, result);

    scope.retrieveDomainsByURI("/domains/?expand&page=1&pagesize=20");
    $httpBackend.flush();

    expect(scope.pagination.numberOfItems).toBe(1);
    expect(scope.pagination.numberOfPages).toBe(1);
    expect(scope.pagination.pageSize).toBe(20);
    expect(scope.pagination.domains.length).toBe(1);
    expect(scope.pagination.domains).toEqual(result.domains);
  }));
});