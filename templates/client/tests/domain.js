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

    scope.domain = {
      nameservers: [
        { lastStatus: "OK" }
      ],
      dsset: [
        { lastStatus: "OK" }
      ]
    };

    scope.$digest();
    expect(ctrl.hasErrors(scope.domain)).toBe(false);

    scope.domain = {
      nameservers: [
        { lastStatus: "OK" }
      ],
      dsset: [
        { lastStatus: "OK" },
        { lastStatus: "TIMEOUT" }
      ]
    };

    scope.$digest();
    expect(ctrl.hasErrors(scope.domain)).toBe(true);

    scope.domain = {
      nameservers: [
        { lastStatus: "OK" },
        { lastStatus: "TIMEOUT" }
      ],
      dsset: [
        { lastStatus: "OK" }
      ]
    };

    scope.$digest();
    expect(ctrl.hasErrors(scope.domain)).toBe(true);

    scope.domain = {
      nameservers: [
        { lastStatus: "NOTCHECKED" }
      ],
      dsset: [
        { lastStatus: "NOTCHECKED" }
      ]
    };

    scope.$digest();
    expect(ctrl.hasErrors(scope.domain)).toBe(false);
  });

  it("verify if a domain was checked", function() {
    expect(ctrl.wasChecked).not.toBeUndefined();

    scope.domain = {
      nameservers: [
        { lastStatus: "NOTCHECKED" }
      ],
      dsset: [
        { lastStatus: "NOTCHECKED" }
      ]
    };

    scope.$digest();
    expect(ctrl.wasChecked(scope.domain)).toBe(false);

    scope.domain = {
      nameservers: [
        { lastStatus: "NOTCHECKED" }
      ],
      dsset: [
        { lastStatus: "OK" }
      ]
    };

    scope.$digest();
    expect(ctrl.wasChecked(scope.domain)).toBe(true);

    scope.domain = {
      nameservers: [
        { lastStatus: "OK" }
      ],
      dsset: [
        { lastStatus: "NOTCHECKED" }
      ]
    };

    scope.$digest();
    expect(ctrl.wasChecked(scope.domain)).toBe(true);
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
  });
});