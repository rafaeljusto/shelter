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
});