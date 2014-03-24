describe("Domain directive", function() {
  var elm, scope;

  beforeEach(module('shelter'));
  beforeEach(module('directives'));

  beforeEach(inject(function($rootScope, $compile, $injector) {
    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET("/languages/en_US.json").respond("");

    elm = angular.element("<domain domain='domain'></domain>");

    scope = $rootScope;
    scope.domain = {};

    $compile(elm)(scope);
    scope.$digest();
  }));

  it("verify if a domain has errors", function() {
    expect(elm.scope().$$childTail.hasErrors).not.toBeUndefined();

    scope.domain = {
      nameservers: [
        { lastStatus: "OK" }
      ],
      dsset: [
        { lastStatus: "OK" }
      ]
    };

    scope.$digest();
    expect(elm.scope().$$childTail.hasErrors(scope.domain)).toBe(false);

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
    expect(elm.scope().$$childTail.hasErrors(scope.domain)).toBe(true);

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
    expect(elm.scope().$$childTail.hasErrors(scope.domain)).toBe(true);

    scope.domain = {
      nameservers: [
        { lastStatus: "NOTCHECKED" }
      ],
      dsset: [
        { lastStatus: "NOTCHECKED" }
      ]
    };

    scope.$digest();
    expect(elm.scope().$$childTail.hasErrors(scope.domain)).toBe(false);
  });

  it("verify if a domain was checked", function() {
    expect(elm.scope().$$childTail.wasChecked).not.toBeUndefined();

    scope.domain = {
      nameservers: [
        { lastStatus: "NOTCHECKED" }
      ],
      dsset: [
        { lastStatus: "NOTCHECKED" }
      ]
    };

    scope.$digest();
    expect(elm.scope().$$childTail.wasChecked(scope.domain)).toBe(false);

    scope.domain = {
      nameservers: [
        { lastStatus: "NOTCHECKED" }
      ],
      dsset: [
        { lastStatus: "OK" }
      ]
    };

    scope.$digest();
    expect(elm.scope().$$childTail.wasChecked(scope.domain)).toBe(true);

    scope.domain = {
      nameservers: [
        { lastStatus: "OK" }
      ],
      dsset: [
        { lastStatus: "NOTCHECKED" }
      ]
    };

    scope.$digest();
    expect(elm.scope().$$childTail.wasChecked(scope.domain)).toBe(true);
  });

  it("verify if the date is defined", function() {
    expect(elm.scope().$$childTail.dateDefined).not.toBeUndefined();
    expect(elm.scope().$$childTail.dateDefined("2014-03-24T14:13:15-03:00")).toBe(true);
    expect(elm.scope().$$childTail.dateDefined("2014-03-24T14:13:15Z")).toBe(true);
    expect(elm.scope().$$childTail.dateDefined("1969-01-01T00:00:00Z")).toBe(false);
    expect(elm.scope().$$childTail.dateDefined("This is not a date")).toBe(false);
  });
});