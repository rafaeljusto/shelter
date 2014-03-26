describe("Domain controller", function() {
  var scope, ctrl;

  beforeEach(module('shelter'));

  beforeEach(inject(function($rootScope, $controller) {
    scope = $rootScope.$new();
    ctrl = $controller("domainCtrl", {
      $scope: scope
    });
  }));

  it("should have an empty domain", function() {
    expect(scope.emptyDomain).toEqual(emptyDomain);
  });
});