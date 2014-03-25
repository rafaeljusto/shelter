describe("Domain directive", function() {
  var elm, scope, ctrl;

  beforeEach(module('shelter'));
  beforeEach(module('directives'));

  beforeEach(inject(function($rootScope, $compile, $injector) {
    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET("/languages/en_US.json").respond(200, "{}");
    $httpBackend.flush()

    elm = angular.element("<domain domain='domain'></domain>");

    scope = $rootScope;
    scope.domain = {};

    $compile(elm)(scope);
    scope.$digest();

    ctrl = elm.scope().$$childTail;
  }));

  it("verify if a domain has errors", function() {
    expect(ctrl.hasErrors).not.toBeUndefined();

    expect(ctrl.hasErrors({
      nameservers: [
        { lastStatus: "OK" }
      ],
      dsset: [
        { lastStatus: "OK" }
      ]
    })).toBe(false);

    expect(ctrl.hasErrors({
      nameservers: [
        { lastStatus: "OK" }
      ],
      dsset: [
        { lastStatus: "OK" },
        { lastStatus: "TIMEOUT" }
      ]
    })).toBe(true);

    expect(ctrl.hasErrors({
      nameservers: [
        { lastStatus: "OK" },
        { lastStatus: "TIMEOUT" }
      ],
      dsset: [
        { lastStatus: "OK" }
      ]
    })).toBe(true);

    expect(ctrl.hasErrors({
      nameservers: [
        { lastStatus: "NOTCHECKED" }
      ],
      dsset: [
        { lastStatus: "NOTCHECKED" }
      ]
    })).toBe(false);
  });

  it("verify if a domain was checked", function() {
    expect(ctrl.wasChecked).not.toBeUndefined();

    expect(ctrl.wasChecked({
      nameservers: [
        { lastStatus: "NOTCHECKED" }
      ],
      dsset: [
        { lastStatus: "NOTCHECKED" }
      ]
    })).toBe(false);

    expect(ctrl.wasChecked({
      nameservers: [
        { lastStatus: "NOTCHECKED" }
      ],
      dsset: [
        { lastStatus: "OK" }
      ]
    })).toBe(true);

    expect(ctrl.wasChecked({
      nameservers: [
        { lastStatus: "OK" }
      ],
      dsset: [
        { lastStatus: "NOTCHECKED" }
      ]
    })).toBe(true);
  });

  it("verify if the date is defined", function() {
    expect(ctrl.dateDefined).not.toBeUndefined();
    expect(ctrl.dateDefined("2014-03-24T14:13:15-03:00")).toBe(true);
    expect(ctrl.dateDefined("2014-03-24T14:13:15Z")).toBe(true);
    expect(ctrl.dateDefined("1969-01-01T00:00:00Z")).toBe(false);
    expect(ctrl.dateDefined("This is not a date")).toBe(false);
  });

  it("verify if the get language function returns the default language", function() {
    expect(ctrl.getLanguage).not.toBeUndefined();
    expect(ctrl.getLanguage()).toBe("en_US");
  });

  it("should return the correct algorithm name", function() {
    expect(ctrl.getAlgorithm).not.toBeUndefined();
    expect(ctrl.getAlgorithm(1)).toBe("RSA/MD5");
    expect(ctrl.getAlgorithm(2)).toBe("DH");
    expect(ctrl.getAlgorithm(3)).toBe("DSA/SHA1");
    expect(ctrl.getAlgorithm(4)).toBe("ECC");
    expect(ctrl.getAlgorithm(5)).toBe("RSA/SHA1");
    expect(ctrl.getAlgorithm(6)).toBe("DSA/SHA1-NSEC3");
    expect(ctrl.getAlgorithm(7)).toBe("RSA/SHA1-NSEC3");
    expect(ctrl.getAlgorithm(8)).toBe("RSA/SHA256");
    expect(ctrl.getAlgorithm(10)).toBe("RSA/SHA512");
    expect(ctrl.getAlgorithm(12)).toBe("GOST R");
    expect(ctrl.getAlgorithm(13)).toBe("ECDSA/SHA256");
    expect(ctrl.getAlgorithm(14)).toBe("ECDSA/SHA384");
    expect(ctrl.getAlgorithm(99)).toBe(99);
  });

  it("should return the correct DS digest type name", function() {
    expect(ctrl.getDSDigestType).not.toBeUndefined();
    expect(ctrl.getDSDigestType(1)).toBe("SHA1");
    expect(ctrl.getDSDigestType(2)).toBe("SHA256");
    expect(ctrl.getDSDigestType(3)).toBe("GOST94");
    expect(ctrl.getDSDigestType(4)).toBe("SHA384");
    expect(ctrl.getDSDigestType(5)).toBe("SHA512");
    expect(ctrl.getDSDigestType(99)).toBe(99);
  });

  it("should show only a part of the DS digest", function() {
    expect(ctrl.showDSDigest).not.toBeUndefined();
    expect(ctrl.showDSDigest(undefined)).toBe("");
    expect(ctrl.showDSDigest("EAA0978F38879DB70A53F9FF1ACF21D046A98B5C")).toBe("EAA0978F3887...21D046A98B5C");
    expect(ctrl.showDSDigest("EAA0978F38879DB70A53F9F")).toBe("EAA0978F38879DB70A53F9F");
  });
});