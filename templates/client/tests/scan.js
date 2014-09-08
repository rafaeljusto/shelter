/*
 * Copyright 2014 Rafael Dantas Justo. All rights reserved.
 * Use of this source code is governed by a GPL
 * license that can be found in the LICENSE file.
 */

describe("Scan directive", function() {
  var scope, ctrl;

  beforeEach(module('shelter'));
  beforeEach(module('directives'));

  beforeEach(inject(function($rootScope, $compile, $injector) {
    localStorage.clear();

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET(/languages\/.+\.json/).respond(200, "{}");
    $httpBackend.flush();

    var elm = angular.element("<scan scan='scan'></scan>");

    scope = $rootScope;
    scope.scan = {};

    $compile(elm)(scope);
    scope.$digest();

    ctrl = elm.isolateScope();
  }));

  it("should count number of statistics", function() {
    expect(ctrl.countStatistics).not.toBeUndefined();

    expect(ctrl.countStatistics({
      "key01": 2,
      "key02": 6,
      "key03": 1,
      "key04": 1
    })).toBe(10);
  });

  it("should verify if the get language function returns the default language", inject(function($translate) {
    expect(ctrl.getLanguage).not.toBeUndefined();
    expect(ctrl.getLanguage()).toBe($translate.preferredLanguage());
    expect(ctrl.getLanguage()).not.toBe("");
    expect(ctrl.getLanguage()).not.toBe(undefined);
  }));

  it("should retrieve a scan", inject(function($injector) {
    expect(ctrl.retrieveScan).not.toBeUndefined();

    $httpBackend = $injector.get("$httpBackend");
    $httpBackend.whenGET("/scan/2014-06-03T11:25:21.114-03:00").respond(200, {"status": "EXECUTED"});

    var scan = {
      links: [
        {
          types: [ "self" ],
          href: "/scan/2014-06-03T11:25:21.114-03:00"
        }
      ],
      etag: 1
    };

    ctrl.retrieveScan(scan);
    $httpBackend.flush();

    expect(ctrl.freshScan.status).toBe("EXECUTED");
  }));
});