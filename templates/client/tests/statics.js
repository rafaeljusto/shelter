/*
 * Copyright 2014 Rafael Dantas Justo. All rights reserved.
 * Use of this source code is governed by a GPL
 * license that can be found in the LICENSE file.
 */

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