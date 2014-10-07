/*
 * Copyright 2014 Rafael Dantas Justo. All rights reserved.
 * Use of this source code is governed by a GPL
 * license that can be found in the LICENSE file.
 */

describe("Domains controller", function() {
  var scope, ctrl;

  beforeEach(module('shelter'));

  beforeEach(inject(function($rootScope, $controller, $injector) {
    localStorage.clear();

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
    $httpBackend.whenGET("/domains/").respond(200, result);
    $httpBackend.whenGET("/domains/?page=1&pagesize=20").respond(200, result);

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
    $httpBackend.whenGET("/domains/").respond(200, result);
    $httpBackend.whenGET("/domains/?page=1&pagesize=20").respond(200, result);

    scope.retrieveDomainsByURI("/domains/?page=1&pagesize=20");
    $httpBackend.flush();

    expect(scope.pagination.numberOfItems).toBe(1);
    expect(scope.pagination.numberOfPages).toBe(1);
    expect(scope.pagination.pageSize).toBe(20);
    expect(scope.pagination.domains.length).toBe(1);
    expect(scope.pagination.domains).toEqual(result.domains);
  }));

  it("should check if all domains were selected", inject(function($injector) {
    scope.pagination = {
      domains: [
        { fqdn: "example1.com.br." },
        { fqdn: "example2.com.br." },
        { fqdn: "example3.com.br." }
      ]
    };
    scope.selectedDomains = [
      { fqdn: "example1.com.br." },
      { fqdn: "example2.com.br." },
      { fqdn: "example3.com.br." }
    ];

    expect(scope.allDomainsSelected()).toBe(true);

    scope.selectedDomains = [
      { fqdn: "example1.com.br." },
      { fqdn: "example2.com.br." },
    ];

    expect(scope.allDomainsSelected()).toBe(false);
  }));

  it("should select all domains of the current page", inject(function($injector) {
    scope.pagination = {
      domains: [
        { fqdn: "example1.com.br." },
        { fqdn: "example2.com.br." },
        { fqdn: "example3.com.br." }
      ]
    };
    scope.selectedDomains = [
      { fqdn: "example1.com.br." },
    ];

    scope.selectAllDomains();
    expect(scope.selectedDomains.length).toBe(3);
    expect(scope.selectedDomains[0].fqdn).toBe("example1.com.br.");
    expect(scope.selectedDomains[1].fqdn).toBe("example2.com.br.");
    expect(scope.selectedDomains[2].fqdn).toBe("example3.com.br.");
  }));

  it("should deselect all domains of the current page", inject(function($injector) {
    scope.pagination = {
      domains: [
        { fqdn: "example1.com.br." },
        { fqdn: "example2.com.br." },
        { fqdn: "example3.com.br." }
      ]
    };
    scope.selectedDomains = [
      { fqdn: "example1.com.br." },
      { fqdn: "example2.com.br." },
      { fqdn: "example3.com.br." },
      { fqdn: "example4.com.br." }
    ];

    scope.selectAllDomains();
    expect(scope.selectedDomains.length).toBe(1);
    expect(scope.selectedDomains[0].fqdn).toBe("example4.com.br.");
  }));

  it("should remove selected domains", inject(function($injector) {
    var result = {
      numberOfItems: 2,
      numberOfPages: 1,
      pageSize: 20,
      domains: [
        { fqdn: "example1.com.br." },
        { fqdn: "example2.com.br." }
      ]
    };

    scope.selectedDomains = [
      {
        fqdn: "example1.com.br.",
        links: [
          { types: [ "self" ], href: "/domain/example1.com.br." }
        ]
      },
      {
        fqdn: "example2.com.br.",
        links: [
          { types: [ "self" ], href: "/domain/example1.com.br." }
        ]
      }
    ];

    httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET("/domains/").respond(200, result);
    $httpBackend.whenDELETE("/domain/example1.com.br.").respond(204);
    $httpBackend.whenDELETE("/domain/example2.com.br.").respond(204);

    scope.removeDomains();
    $httpBackend.flush();

    expect(scope.selectedDomains.length).toBe(0);
  }));
});