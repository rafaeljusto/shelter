/*
 * Copyright 2014 Rafael Dantas Justo. All rights reserved.
 * Use of this source code is governed by a GPL
 * license that can be found in the LICENSE file.
 */

describe("mergeList", function() {
  it("Should add a new item", function() {
    a = [
      {
        key1: "value1",
        key2: "value2"
      },
      {
        key1: "value3",
        key2: "value4"
      },
    ];

    b = [
      {
        key1: "value1",
        key2: "value2"
      },
      {
        key1: "value3",
        key2: "value4"
      },
      {
        key1: "value5",
        key2: "value6"
      },
    ];

    mergeList(b, a,
      function(newItem, oldItem) {
        return newItem.key1 == oldItem.key1;
      },
      function(newItem, oldItem) {
        oldItem.key2 = newItem.key2;
      });

    expect(a.length).toBe(3);
    expect(a[2].key2).toBe("value6");
  });

  it("Should remove old items and update existing item", function() {
    a = [
      {
        key1: "value1",
        key2: "value2"
      },
      {
        key1: "value3",
        key2: "value4"
      },
    ];

    b = [
      {
        key1: "value1",
        key2: "value5"
      }
    ];

    mergeList(b, a,
      function(newItem, oldItem) {
        return newItem.key1 == oldItem.key1;
      },
      function(newItem, oldItem) {
        oldItem.key2 = newItem.key2;
      });

    expect(a.length).toBe(1);
    expect(a[0].key2).toBe("value5");
  });

  it("Should ignore the operation when an input is null", function() {
    a = [
      {
        key1: "value1",
        key2: "value2"
      },
      {
        key1: "value3",
        key2: "value4"
      },
    ];

    b = [
      {
        key1: "value1",
        key2: "value5"
      }
    ];

    mergeList(undefined, a,
      function(newItem, oldItem) {
        return newItem.key1 == oldItem.key1;
      },
      function(newItem, oldItem) {
        oldItem.key2 = newItem.key2;
      });

    expect(a.length).toBe(2);

    a = undefined
    mergeList(b, a,
      function(newItem, oldItem) {
        return newItem.key1 == oldItem.key1;
      },
      function(newItem, oldItem) {
        oldItem.key2 = newItem.key2;
      });

    expect(a).toBe(undefined);
  });

  it("Should add the first item", function() {
    a = [];

    b = [
      {
        key1: "value1",
        key2: "value2"
      }
    ];

    mergeList(b, a,
      function(newItem, oldItem) {
        return newItem.key1 == oldItem.key1;
      },
      function(newItem, oldItem) {
        oldItem.key2 = newItem.key2;
      });

    expect(a.length).toBe(1);
    expect(a[0].key2).toBe("value2");

    b = [];

    a = [
      {
        key1: "value1",
        key2: "value2"
      }
    ];

    mergeList(b, a,
      function(newItem, oldItem) {
        return newItem.key1 == oldItem.key1;
      },
      function(newItem, oldItem) {
        oldItem.key2 = newItem.key2;
      });

    expect(a.length).toBe(0);
  });

  it("Should deal with null structures", function() {
    a = [];

    b = [
      {
        key1: "value1",
        key2: "value2"
      }
    ];

    a = mergeList(b, undefined,
      function(newItem, oldItem) {
        return newItem.key1 == oldItem.key1;
      },
      function(newItem, oldItem) {
        oldItem.key2 = newItem.key2;
      });

    expect(a.length).toBe(1);
    expect(a[0].key2).toBe("value2");

    a = [
      {
        key1: "value1",
        key2: "value2"
      }
    ];

    a = mergeList(undefined, a,
      function(newItem, oldItem) {
        return newItem.key1 == oldItem.key1;
      },
      function(newItem, oldItem) {
        oldItem.key2 = newItem.key2;
      });

    expect(a.length).toBe(0);
  });
});

describe("findLink", function() {
  it("should find a existing link", function() {
    var links = [
      {
        href: "/path/to/object1",
        types: ["first", "prev", "last"]
      },
      {
        href: "/path/to/object2",
        types: ["next"]
      }
    ];

    expect(findLink(links, "first")).toBe("/path/to/object1");
    expect(findLink(links, "prev")).toBe("/path/to/object1");
    expect(findLink(links, "last")).toBe("/path/to/object1");
    expect(findLink(links, "next")).toBe("/path/to/object2");
  });

  it("should return empty when the link doesn't exist", function() {
    var links = [
      {
        href: "/path/to/object1",
        types: ["first", "prev", "last"]
      },
      {
        href: "/path/to/object2",
        types: ["next"]
      }
    ];

    expect(findLink(links, "current")).toBe("");
  });

  it("should return empty when the structures are undefined", function() {
    var links = [
      {
        href: "/path/to/object1",
        types: ["first", "prev", "last"]
      },
      {
        href: "/path/to/object2",
        types: ["next"]
      }
    ];

    expect(findLink(undefined, "first")).toBe("");
    expect(findLink(links, undefined)).toBe("");
  });
});

describe("verificationResponseToHTML", function() {
  it("should convert structure correctly", function() {
    var data = {
      fqdn: "test.com.br.",
      nameservers: [
        { host: "ns1.test.com.br.", lastStatus: "OK" },
        { host: "ns2.test.com.br.", lastStatus: "TIMEOUT" }
      ],
      dsset: [
        { keytag: 1234, lastStatus: "OK" },
        { keytag: 4321, lastStatus: "SIGERR" }
      ]
    }

    expect(verificationResponseToHTML(data)).toBe("<h3>test.com.br.</h3><hr/>" +
      "<h3>NS</h3><table style='margin:auto'>" +
      "<tr><th style='text-align:left'>ns1.test.com.br.</th><td>OK</td></tr>" +
      "<tr><th style='text-align:left'>ns2.test.com.br.</th><td>TIMEOUT</td></tr>" +
      "</table><h3>DS</h3><table style='margin:auto'>" +
      "<tr><th style='text-align:left'>1234</th><td>OK</td></tr>" +
      "<tr><th style='text-align:left'>4321</th><td>SIGERR</td></tr>" +
      "</table>");
  });

  it ("should not convert an invalid structure", function() {
    expect(verificationResponseToHTML({})).toBe("");
  });

  it ("should not convert an undefined structure", function() {
    expect(verificationResponseToHTML(undefined)).toBe("");
  });
});