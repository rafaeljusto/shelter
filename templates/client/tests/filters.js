describe("Filters", function() {
  var datetimeFilter, rangeFilter;

  beforeEach(module('shelter'));

  beforeEach(inject(function($filter) {
    datetimeFilter = $filter("datetime");
    rangeFilter = $filter("range");
  }));

  it("should format correctly a datetime", function() {
    moment().zone("-03:00");
    expect(datetimeFilter("2014-03-26T09:55:29-03:00", "en_US")).toBe("March 26th 2014, 9:55:29 am");
    expect(datetimeFilter("2014-03-26T09:55:29-03:00", "pt_BR")).toBe("Março 26º 2014, 9:55:29 am");
  });
});