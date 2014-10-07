/*
 * Copyright 2014 Rafael Dantas Justo. All rights reserved.
 * Use of this source code is governed by a GPL
 * license that can be found in the LICENSE file.
 */

describe("Domain form directive", function() {
  var scope, ctrl;

  beforeEach(module('shelter'));
  beforeEach(module('directives'));

  beforeEach(inject(function($rootScope, $compile, $injector) {
    localStorage.clear();

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET(/languages\/.+\.json/).respond(200, "{}");
    $httpBackend.flush();

    var elm = angular.element("<domainform domain='domain' form-ctrl='ctrl'></domain>");

    scope = $rootScope;
    scope.domain = {};
    scope.ctrl = {};

    $compile(elm)(scope);
    scope.$digest();

    ctrl = elm.isolateScope();
  }));

  it("should detect when a nameserver needs glue records", function() {
    expect(ctrl.needsGlue).not.toBeUndefined();

    expect(ctrl.needsGlue("test1.com.br.", "ns1.test1.com.br.")).toBe(true);
    expect(ctrl.needsGlue("test1.com.br.", "ns1.test1.com.br")).toBe(true);
    expect(ctrl.needsGlue("test1.com.br", "ns1.test1.com.br.")).toBe(true);
    expect(ctrl.needsGlue("test1.com.br", "ns1.test1.com.br")).toBe(true);
    expect(ctrl.needsGlue("test2.com.br.", "ns1.test1.com.br.")).toBe(false);
    expect(ctrl.needsGlue(undefined, "ns1.test1.com.br.")).toBe(false);
    expect(ctrl.needsGlue("test1.com.br.", undefined)).toBe(false);
  });

  it("should add to list", function() {
    expect(ctrl.addToList).not.toBeUndefined();

    var list = []; var obj = {};

    ctrl.addToList(obj, list);
    expect(list.length).toBe(1);
    expect(obj).not.toBe(list[0]);

    ctrl.addToList(undefined, list);
    expect(list.length).toBe(1);
  });

  it("should remove from list", function() {
    expect(ctrl.removeFromList).not.toBeUndefined();

    var list = []; var obj = {};
    list.push(obj);

    ctrl.removeFromList(0, list);
    expect(list.length).toBe(0);
    ctrl.removeFromList(0, list);
  });

  it("should query a domain correctly", inject(function($injector) {
    expect(ctrl.queryDomain).not.toBeUndefined();

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET("/domain/br./verification").respond(200, {
      fqdn: "br.",
      nameservers: [
        {
          host: "a.dns.br",
          ipv4: "200.160.0.10"
        }
      ],
      "dsset": [
        {
          keytag: 41674,
          algorithm: 5,
          digestType: 1,
          digest: "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C"
        }
      ]
    });

    ctrl.queryDomain("br.");
    $httpBackend.flush();

    expect(ctrl.domain).not.toBeUndefined();
    expect(ctrl.domain.nameservers.length).toBe(1);
    expect(ctrl.domain.nameservers[0].host).toBe("a.dns.br");
    expect(ctrl.domain.dsset.length).toBe(1);
    expect(ctrl.domain.dsset[0].keytag).toBe(41674);
    expect(ctrl.domain.owners).not.toBeUndefined();
    expect(ctrl.domain.dnskeys).not.toBeUndefined();
  }));

  it("should verify a domain", inject(function($injector) {
    expect(ctrl.verifyDomain).not.toBeUndefined();

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenPUT("/domain/br./verification").respond(200, {
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
      "dsset": [
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
    });

    ctrl.verifyDomain({
      fqdn: "br.",
      nameservers: [
        {
          host: "a.dns.br.",
          ipv4: "200.160.0.10"
        }
      ],
      dsset: [
        {
          keytag: "41674",
          algorithm: 5,
          digestType: 1,
          digest: "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C"
        }
      ],
    });

    $httpBackend.flush();

    expect(ctrl.verifyResult).not.toBeUndefined();
    expect(ctrl.verifyResult.nameservers.length).toBe(1);
    expect(ctrl.verifyResult.nameservers[0].lastStatus).toBe("OK");
    expect(ctrl.verifyResult.nameservers[0].lastOKAt).toBe("2014-03-25T11:00:00-03:00");
    expect(ctrl.verifyResult.nameservers[0].lastCheckAt).toBe("2014-03-25T11:00:00-03:00");
    expect(ctrl.verifyResult.dsset.length).toBe(1);
    expect(ctrl.verifyResult.dsset[0].lastStatus).toBe("OK");
    expect(ctrl.verifyResult.dsset[0].lastOKAt).toBe("2014-03-25T11:00:00-03:00");
    expect(ctrl.verifyResult.dsset[0].lastCheckAt).toBe("2014-03-25T11:00:00-03:00");
  }));

  it("should save a domain properly", inject(function($injector) {
    expect(ctrl.saveDomain).not.toBeUndefined();

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenPUT("/domain/br.").respond(201);

    ctrl.saveDomain({
      fqdn: "br.",
      nameservers: [
        {
          host: "a.dns.br.",
          ipv4: "200.160.0.10"
        }
      ],
      dsset: [
        {
          keytag: "41674",
          algorithm: 5,
          digestType: 1,
          digest: "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C"
        }
      ],
    });

    $httpBackend.flush();

    expect(ctrl.success).not.toBeUndefined();
    expect(ctrl.success).toBe(true);
  }));

  it("should store a CSV content correctly", inject(function($injector) {
    expect(ctrl.storeCSVFile).not.toBeUndefined();
    expect(ctrl.csv).not.toBeUndefined();

    ctrl.storeCSVFile("example1.com,ns1.example1.com$127.0.0.1$::1,ns2.example1.com$127.0.0.2$::2,,,,,123$5$1$EAA0978F38879DB70A53F9FF1ACF21D046A98B5C,,test@example1.com$en-us,,")
    expect(ctrl.csv.content).toBe("example1.com,ns1.example1.com$127.0.0.1$::1,ns2.example1.com$127.0.0.2$::2,,,,,123$5$1$EAA0978F38879DB70A53F9FF1ACF21D046A98B5C,,test@example1.com$en-us,,");
  }));

  it("should import a CSV correctly", inject(function($injector) {
    expect(ctrl.importCSV).not.toBeUndefined();
    expect(ctrl.csv).not.toBeUndefined();

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenPUT("/domain/example1.com").respond(201);
    $httpBackend.whenPUT("/domain/example2.com").respond(201);

    ctrl.csv.content = "example1.com,ns1.example1.com$127.0.0.1$::1,ns2.example1.com$127.0.0.2$::2,,,,,123$5$1$EAA0978F38879DB70A53F9FF1ACF21D046A98B5C,,test@example1.com$en-us,,\nexample2.com,ns1.example1.com$$,ns2.example1.com$$,,,,,,,,,";
    ctrl.importCSV();

    $httpBackend.flush();

    expect(ctrl.csv.domainsToUpload).toBe(2);
    expect(ctrl.csv.domainsUploaded).toBe(2);
    expect(ctrl.csv.success).toBe(2);
    expect(ctrl.csv.errors.length).toBe(0);
  }));

  it("should detect CSV format errors", inject(function($injector) {
    expect(ctrl.importCSV).not.toBeUndefined();
    expect(ctrl.csv).not.toBeUndefined();

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenPUT("/domain/example1.com").respond(201);

    ctrl.csv.content = "example1.com,ns1.example1.com$127.0.0.1$::1,ns2.example1.com$127.0.0.2$::2,,,,,123$5$1$EAA0978F38879DB70A53F9FF1ACF21D046A98B5C,,test@example1.com$en-us,,\nexample2.com,ns1.example1.com$$,ns2.example1.com$$,,,,,,,,";
    ctrl.importCSV();

    $httpBackend.flush();

    expect(ctrl.csv.domainsToUpload).toBe(2);
    expect(ctrl.csv.domainsUploaded).toBe(2);
    expect(ctrl.csv.success).toBe(1);
    expect(ctrl.csv.errors.length).toBe(1);
    expect(ctrl.csv.errors[0].lineNumber).toBe(2);
  }));

  it("should detect CSV network errors", inject(function($injector) {
    expect(ctrl.importCSV).not.toBeUndefined();
    expect(ctrl.csv).not.toBeUndefined();

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenPUT("/domain/example1.com").respond(500);
    $httpBackend.whenPUT("/domain/example2.com").respond(201);

    ctrl.csv.content = "example1.com,ns1.example1.com$127.0.0.1$::1,ns2.example1.com$127.0.0.2$::2,,,,,123$5$1$EAA0978F38879DB70A53F9FF1ACF21D046A98B5C,,test@example1.com$en-us,,\nexample2.com,ns1.example1.com$$,ns2.example1.com$$,,,,,,,,,";
    ctrl.importCSV();

    $httpBackend.flush();

    expect(ctrl.csv.domainsToUpload).toBe(2);
    expect(ctrl.csv.domainsUploaded).toBe(2);
    expect(ctrl.csv.success).toBe(1);
    expect(ctrl.csv.errors.length).toBe(1);
    expect(ctrl.csv.errors[0].lineNumber).toBe(1);
  }));
});