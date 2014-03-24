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

  it("verify if a domain has errors", inject(function($compile, $rootScope) {
    expect(elm.scope().$$childTail.hasErrors).not.toBeUndefined();

    scope.domain = {
      nameservers: [
        { status: "OK" }
      ],
      dsset: [
        { status: "OK" }
      ]
    };

    scope.$digest();
    expect(elm.scope().$$childTail.hasErrors(scope.domain)).toBe(false);

    scope.domain = {
      nameservers: [
        { status: "OK" }
      ],
      dsset: [
        { status: "OK" },
        { status: "TIMEOUT" }
      ]
    };

    scope.$digest();
    expect(elm.scope().$$childTail.hasErrors(scope.domain)).toBe(true);

    scope.domain = {
      nameservers: [
        { status: "OK" },
        { status: "TIMEOUT" }
      ],
      dsset: [
        { status: "OK" }
      ]
    };

    scope.$digest();
    expect(elm.scope().$$childTail.hasErrors(scope.domain)).toBe(true);

    scope.domain = {
      nameservers: [
        { status: "NOTCHECKED" }
      ],
      dsset: [
        { status: "NOTCHECKED" }
      ]
    };

    scope.$digest();
    expect(elm.scope().$$childTail.hasErrors(scope.domain)).toBe(false);
  }));
});