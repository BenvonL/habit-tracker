import {
  toLocalISODate,
  parseLocalISODate,
  shiftDate,
  countToLevel,
} from "./date";

describe("toLocalISODate", () => {
  it("uses the local calendar day", () => {
    expect(toLocalISODate(new Date(2026, 10, 4, 0, 30))).toBe("2026-11-04");
    expect(toLocalISODate(new Date(2026, 0, 1, 23, 59))).toBe("2026-01-01");
  });
});

describe("parseLocalISODate", () => {
  it("parses as local midnight", () => {
    const d = parseLocalISODate("2026-03-29");
    expect([d.getFullYear(), d.getMonth(), d.getDate(), d.getHours()]).toEqual(
      [2026, 2, 29, 0]
    );
  });
});

describe("shiftDate", () => {
  it("shifts across month and year boundaries", () => {
    expect(shiftDate("2026-01-01", -1)).toBe("2025-12-31");
    expect(shiftDate("2026-02-28", 1)).toBe("2026-03-01");
    expect(shiftDate("2026-11-04", -363)).toBe("2025-11-06");
  });
});

describe("countToLevel", () => {
  it("returns 0 for no activity", () => {
    expect(countToLevel(0, 10)).toBe(0);
    expect(countToLevel(5, 0)).toBe(0);
  });
  it("scales relative to the maximum and never exceeds 4", () => {
    expect(countToLevel(1, 10)).toBe(1);
    expect(countToLevel(5, 10)).toBe(2);
    expect(countToLevel(10, 10)).toBe(4);
    expect(countToLevel(12, 10)).toBe(4);
    expect(countToLevel(1, 1)).toBe(4);
  });
});
