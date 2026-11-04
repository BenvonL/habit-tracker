// Date helpers that work in the user's local time zone.
// toISOString() would use UTC and put entries made shortly after midnight
// (e.g. 00:30 in Germany) on the previous day.

const pad = (n) => String(n).padStart(2, "0");

// Formats a Date as YYYY-MM-DD in local time.
export function toLocalISODate(date = new Date()) {
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(
    date.getDate()
  )}`;
}

// Parses a YYYY-MM-DD string as a local date (not UTC midnight).
export function parseLocalISODate(value) {
  const [year, month, day] = value.split("-").map(Number);
  return new Date(year, month - 1, day);
}

// Returns the YYYY-MM-DD string numDays after (or before) the given one.
export function shiftDate(value, numDays) {
  const date = parseLocalISODate(value);
  date.setDate(date.getDate() + numDays);
  return toLocalISODate(date);
}

// Maps a count to a heatmap colour level from 0 (nothing) to 4 (max),
// relative to the highest count of the habit, since units differ
// (litres, hours, ...).
export function countToLevel(count, maxCount) {
  if (!count || count <= 0 || !maxCount) {
    return 0;
  }
  return Math.min(4, Math.max(1, Math.ceil((count / maxCount) * 4)));
}
