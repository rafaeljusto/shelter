describe("Languages controller", function() {
  var scope, ctrl;

  beforeEach(module('shelter'));

  beforeEach(inject(function($rootScope, $controller, $injector) {
    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET("/languages/en_US.json").respond(200, "{}");
    $httpBackend.whenGET("/languages/pt_BR.json").respond(200, "{}");
    $httpBackend.flush()

    scope = $rootScope.$new();
    ctrl = $controller("languagesCtrl", {
      $scope: scope
    });
  }));

  it("should verify if the get language function returns the default language", function() {
    expect(scope.getLanguage()).toBe("en_US");
  });

  it("should change language correctly", inject(function($injector) {
    $httpBackend = $injector.get("$httpBackend");
    scope.changeLanguage("pt_BR");
    $httpBackend.flush()

    expect(scope.getLanguage()).toBe("pt_BR");
  }));
});