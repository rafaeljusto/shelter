describe("Domain form directive", function() {
  var scope, ctrl;

  beforeEach(module('shelter'));
  beforeEach(module('directives'));

  beforeEach(inject(function($rootScope, $compile, $injector) {
    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET("/languages/en_US.json").respond(200, "{}");
    $httpBackend.flush()

    var elm = angular.element("<domainform domain='domain'></domain>");

    scope = $rootScope;
    scope.domain = {};

    $compile(elm)(scope);
    scope.$digest();

    ctrl = elm.isolateScope();
  }));

  it("should detect when a nameserver needs glue records", function() {
    expect(ctrl.needsGlue).not.toBeUndefined();

    expect(ctrl.needsGlue("test1.com.br.", "ns1.test1.com.br.")).toBe(true);
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
    $httpBackend.flush()

    expect(ctrl.domain).not.toBeUndefined();
    expect(ctrl.domain.nameservers.length).toBe(1);
    expect(ctrl.domain.nameservers[0].host).toBe("a.dns.br");
    expect(ctrl.domain.dsset.length).toBe(1);
    expect(ctrl.domain.dsset[0].keytag).toBe(41674);
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

    $httpBackend.flush()

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

    $httpBackend.flush()

    expect(ctrl.success).not.toBeUndefined();
    expect(ctrl.success).toBe("Domain created");
  }));
});