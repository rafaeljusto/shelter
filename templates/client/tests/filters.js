/*
 * Copyright 2014 Rafael Dantas Justo. All rights reserved.
 * Use of this source code is governed by a GPL
 * license that can be found in the LICENSE file.
 */

describe("Filters", function() {
  var datetimeFilter, rangeFilter;

  beforeEach(module('shelter'));

  beforeEach(inject(function($filter) {
    localStorage.clear();
    datetimeFilter = $filter("datetime");
    rangeFilter = $filter("range");
  }));

  it("should format correctly a datetime", function() {
    expect(datetimeFilter("2014-03-26T09:55:29Z", "en_US")).toBe("Wednesday, March 26 2014 9:55 AM");
    expect(datetimeFilter("2014-03-26T09:55:29.12345Z", "en_US")).toBe("Wednesday, March 26 2014 9:55 AM");
    expect(datetimeFilter("2014-03-26T09:55:29Z", "pt_BR")).toBe("quarta-feira, 26 de março de 2014 às 09:55");
    expect(datetimeFilter("2014-03-26T09:55:29Z", "es_ES")).toBe("miércoles, 26 de marzo del 2014 9:55");
    expect(datetimeFilter("Invalid date", "pt_BR")).toBe("");
    expect(datetimeFilter("2014-03-26T09:55:29Z", undefined)).toBe("");
    expect(datetimeFilter("1968-01-01T09:55:29Z", "pt_BR")).toBe("");

  });
});