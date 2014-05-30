/*
 * Copyright 2014 Rafael Dantas Justo. All rights reserved.
 * Use of this source code is governed by a GPL
 * license that can be found in the LICENSE file.
 */

describe("Languages controller", function() {
  var scope, ctrl;

  beforeEach(module('shelter'));

  beforeEach(inject(function($rootScope, $controller, $injector) {
    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET(/languages\/.+\.json/).respond(200, "{}");
    $httpBackend.flush();

    scope = $rootScope.$new();
    ctrl = $controller("languagesCtrl", {
      $scope: scope
    });
  }));

  it("should verify if the get language function returns the default language", inject(function($translate) {
    expect(scope.getLanguage()).toBe($translate.preferredLanguage());
    expect(scope.getLanguage()).not.toBe("");
    expect(scope.getLanguage()).not.toBe(undefined);
  }));

  it("should change language correctly", inject(function($injector) {
    $httpBackend = $injector.get("$httpBackend");

    // Must be an unknown language, because we don't known the default system language
    scope.changeLanguage("xx_ZZ");
    $httpBackend.flush();

    expect(scope.getLanguage()).toBe("xx_ZZ");
  }));
});