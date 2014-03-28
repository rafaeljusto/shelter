/*
 * Copyright 2014 Rafael Dantas Justo. All rights reserved.
 * Use of this source code is governed by a GPL
 * license that can be found in the LICENSE file.
 */

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