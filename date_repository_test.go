package main

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var _ DateRepository = MySQLRepository{}

var dateRowColumns = []string{"date_id", "date_date", "date_count", "habit_id"}

const (
	selectAllDates = "SELECT date_id, date_date, date_count, habit_id FROM date WHERE habit_id = ? ORDER BY date_date"
	selectTodays   = "SELECT date_id, date_date, date_count, habit_id FROM date WHERE date_date = ? AND habit_id = ?"
	upsertCount    = "INSERT INTO date (date_date, date_count, habit_id) VALUES (?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE date_count = date_count + VALUES(date_count)"
)

func TestGetAllDates(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectAllDates).WithArgs(1).WillReturnRows(
		sqlmock.NewRows(dateRowColumns).
			AddRow(10, "2026-11-01", 2, 1).
			AddRow(11, "2026-11-02", 5, 1),
	)

	got, err := repo.GetAllDates(1)
	if err != nil {
		t.Fatalf("GetAllDates: %v", err)
	}
	want := []Date{
		{ID: 10, Date: "2026-11-01", Count: 2, HabitID: 1},
		{ID: 11, Date: "2026-11-02", Count: 5, HabitID: 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetAllDatesEmptyReturnsEmptySlice(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectAllDates).WithArgs(1).WillReturnRows(sqlmock.NewRows(dateRowColumns))

	got, err := repo.GetAllDates(1)
	if err != nil {
		t.Fatalf("GetAllDates: %v", err)
	}
	if got == nil || len(got) != 0 {
		t.Errorf("got %#v, want empty non-nil slice", got)
	}
}

func TestGetAllDatesQueryError(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectAllDates).WithArgs(1).WillReturnError(errDB)

	if _, err := repo.GetAllDates(1); !errors.Is(err, errDB) {
		t.Errorf("got error %v, want %v", err, errDB)
	}
}

func TestGetAllDatesScanError(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectAllDates).WithArgs(1).WillReturnRows(
		sqlmock.NewRows(dateRowColumns).AddRow(10, "2026-11-01", "lots", 1),
	)

	if _, err := repo.GetAllDates(1); err == nil {
		t.Error("expected a scan error, got nil")
	}
}

func TestGetAllDatesRowError(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectAllDates).WithArgs(1).WillReturnRows(
		sqlmock.NewRows(dateRowColumns).
			AddRow(10, "2026-11-01", 2, 1).
			RowError(0, errDB),
	)

	if _, err := repo.GetAllDates(1); !errors.Is(err, errDB) {
		t.Errorf("got error %v, want %v", err, errDB)
	}
}

func TestGetTodaysCountForHabit(t *testing.T) {
	repo, mock := newMockRepository(t)
	today := time.Now().Format(dateLayout)
	mock.ExpectQuery(selectTodays).WithArgs(today, 4).WillReturnRows(
		sqlmock.NewRows(dateRowColumns).AddRow(20, today, 3, 4),
	)

	got, err := repo.GetTodaysCountForHabit(4)
	if err != nil {
		t.Fatalf("GetTodaysCountForHabit: %v", err)
	}
	want := Date{ID: 20, Date: today, Count: 3, HabitID: 4}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetTodaysCountForHabitUsesLocalDate(t *testing.T) {
	// Pin the local zone far from UTC so that "today" differs from the UTC
	// date for part of every day. The query must use the local date.
	orig := time.Local
	time.Local = time.FixedZone("UTC+14", 14*60*60)
	t.Cleanup(func() { time.Local = orig })

	repo, mock := newMockRepository(t)
	localToday := time.Now().Format(dateLayout)
	mock.ExpectQuery(selectTodays).WithArgs(localToday, 4).WillReturnRows(sqlmock.NewRows(dateRowColumns))

	got, err := repo.GetTodaysCountForHabit(4)
	if err != nil {
		t.Fatalf("GetTodaysCountForHabit: %v", err)
	}
	if got.Date != localToday {
		t.Errorf("got date %q, want local date %q", got.Date, localToday)
	}
}

func TestGetTodaysCountForHabitNothingLogged(t *testing.T) {
	repo, mock := newMockRepository(t)
	today := time.Now().Format(dateLayout)
	mock.ExpectQuery(selectTodays).WithArgs(today, 4).WillReturnRows(sqlmock.NewRows(dateRowColumns))

	got, err := repo.GetTodaysCountForHabit(4)
	if err != nil {
		t.Fatalf("GetTodaysCountForHabit: %v", err)
	}
	// No row is not an error: the caller gets today's date with count 0.
	want := Date{Date: today, Count: 0, HabitID: 4}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetTodaysCountForHabitError(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectTodays).WillReturnError(errDB)

	if _, err := repo.GetTodaysCountForHabit(4); !errors.Is(err, errDB) {
		t.Errorf("got error %v, want %v", err, errDB)
	}
}

func TestAddCountUpsertsInOneStatement(t *testing.T) {
	repo, mock := newMockRepository(t)
	// A single upsert, no read-then-write, so two quick saves cannot lose a count.
	mock.ExpectExec(upsertCount).
		WithArgs("2026-11-27", 3, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.AddCount(Date{Date: "2026-11-27", Count: 3, HabitID: 2}); err != nil {
		t.Fatalf("AddCount: %v", err)
	}
}

func TestAddCountError(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectExec(upsertCount).WillReturnError(errDB)

	if err := repo.AddCount(Date{Date: "2026-11-27", Count: 3, HabitID: 2}); !errors.Is(err, errDB) {
		t.Errorf("got error %v, want %v", err, errDB)
	}
}

func TestAddCountRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		date Date
	}{
		{"missing habit id", Date{Date: "2026-11-27", Count: 1}},
		{"zero count", Date{Date: "2026-11-27", HabitID: 1}},
		{"empty date", Date{Count: 1, HabitID: 1}},
		{"wrong date format", Date{Date: "27.11.2026", Count: 1, HabitID: 1}},
		{"impossible date", Date{Date: "2026-02-30", Count: 1, HabitID: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, _ := newMockRepository(t)
			// No statement may reach the database; the cleanup check enforces that.
			if err := repo.AddCount(tt.date); err == nil {
				t.Errorf("AddCount(%+v): expected an error, got nil", tt.date)
			}
		})
	}
}

func TestAddCountFromJSON(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectExec(upsertCount).
		WithArgs("2026-11-27", 4, 9).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := []byte(`{"date":"2026-11-27","count":4,"habit_id":9}`)
	if err := repo.AddCountFromJSON(req); err != nil {
		t.Fatalf("AddCountFromJSON: %v", err)
	}
}

func TestAddCountFromJSONInvalid(t *testing.T) {
	repo, _ := newMockRepository(t)
	if err := repo.AddCountFromJSON([]byte(`{"count":"three"}`)); err == nil {
		t.Error("expected an error for invalid JSON, got nil")
	}
}
