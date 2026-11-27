package main

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var _ HabitRepository = MySQLRepository{}

var habitRowColumns = []string{"habit_id", "habit_name", "habit_unit", "habit_pomodoro", "habit_why"}

const (
	selectAllHabits = "SELECT habit_id, habit_name, habit_unit, habit_pomodoro, habit_why FROM habit ORDER BY habit_id"
	selectHabit     = "SELECT habit_id, habit_name, habit_unit, habit_pomodoro, habit_why FROM habit WHERE habit_id = ?"
	insertHabit     = "INSERT INTO habit (habit_name, habit_unit, habit_pomodoro, habit_why) VALUES (?, ?, ?, ?)"
	updateHabit     = "UPDATE habit SET habit_name = ?, habit_unit = ?, habit_pomodoro = ?, habit_why = ? WHERE habit_id = ?"
	deleteHabit     = "DELETE FROM habit WHERE habit_id = ?"
)

var errDB = errors.New("db unavailable")

func TestGetAllHabits(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectAllHabits).WillReturnRows(
		sqlmock.NewRows(habitRowColumns).
			AddRow(1, "yoga", "hours", false, "relax").
			AddRow(2, "reading", "pages", true, "learn"),
	)

	got, err := repo.GetAllHabits()
	if err != nil {
		t.Fatalf("GetAllHabits: %v", err)
	}
	want := []Habit{
		{ID: 1, Name: "yoga", Unit: "hours", Pomodoro: false, Why: "relax"},
		{ID: 2, Name: "reading", Unit: "pages", Pomodoro: true, Why: "learn"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetAllHabitsEmptyReturnsEmptySlice(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectAllHabits).WillReturnRows(sqlmock.NewRows(habitRowColumns))

	got, err := repo.GetAllHabits()
	if err != nil {
		t.Fatalf("GetAllHabits: %v", err)
	}
	// The frontend expects a JSON array, so an empty result must not be nil.
	if got == nil || len(got) != 0 {
		t.Errorf("got %#v, want empty non-nil slice", got)
	}
}

func TestGetAllHabitsQueryError(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectAllHabits).WillReturnError(errDB)

	if _, err := repo.GetAllHabits(); !errors.Is(err, errDB) {
		t.Errorf("got error %v, want %v", err, errDB)
	}
}

func TestGetAllHabitsScanError(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectAllHabits).WillReturnRows(
		sqlmock.NewRows(habitRowColumns).AddRow("not-a-number", "yoga", "hours", false, "relax"),
	)

	if _, err := repo.GetAllHabits(); err == nil {
		t.Error("expected a scan error, got nil")
	}
}

func TestGetAllHabitsRowError(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectAllHabits).WillReturnRows(
		sqlmock.NewRows(habitRowColumns).
			AddRow(1, "yoga", "hours", false, "relax").
			RowError(0, errDB),
	)

	if _, err := repo.GetAllHabits(); !errors.Is(err, errDB) {
		t.Errorf("got error %v, want %v", err, errDB)
	}
}

func TestGetHabit(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectHabit).WithArgs(7).WillReturnRows(
		sqlmock.NewRows(habitRowColumns).AddRow(7, "yoga", "hours", true, "relax"),
	)

	got, err := repo.GetHabit(7)
	if err != nil {
		t.Fatalf("GetHabit: %v", err)
	}
	want := Habit{ID: 7, Name: "yoga", Unit: "hours", Pomodoro: true, Why: "relax"}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetHabitNotFound(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectQuery(selectHabit).WithArgs(99).WillReturnRows(sqlmock.NewRows(habitRowColumns))

	if _, err := repo.GetHabit(99); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("got error %v, want sql.ErrNoRows", err)
	}
}

func TestAddHabitLetsDatabaseAssignID(t *testing.T) {
	repo, mock := newMockRepository(t)
	// The ID from the request is ignored; habit_id is AUTO_INCREMENT.
	mock.ExpectExec(insertHabit).
		WithArgs("yoga", "hours", true, "relax").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.AddHabit(Habit{ID: 42, Name: "yoga", Unit: "hours", Pomodoro: true, Why: "relax"})
	if err != nil {
		t.Fatalf("AddHabit: %v", err)
	}
}

func TestAddHabitPassesUserInputAsParameters(t *testing.T) {
	repo, mock := newMockRepository(t)
	// Quotes in user input must reach the driver as plain values, not as SQL.
	name := "it's yoga'); DROP TABLE habit; --"
	mock.ExpectExec(insertHabit).
		WithArgs(name, "hours", false, "why's that").
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.AddHabit(Habit{Name: name, Unit: "hours", Why: "why's that"}); err != nil {
		t.Fatalf("AddHabit: %v", err)
	}
}

func TestAddHabitError(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectExec(insertHabit).WillReturnError(errDB)

	if err := repo.AddHabit(Habit{Name: "yoga"}); !errors.Is(err, errDB) {
		t.Errorf("got error %v, want %v", err, errDB)
	}
}

func TestAddHabitFromJSON(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectExec(insertHabit).
		WithArgs("yoga", "hours", true, "relax").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := []byte(`{"name":"yoga","unit":"hours","pomodoro":true,"why":"relax"}`)
	if err := repo.AddHabitFromJSON(req); err != nil {
		t.Fatalf("AddHabitFromJSON: %v", err)
	}
}

func TestAddHabitFromJSONInvalid(t *testing.T) {
	repo, _ := newMockRepository(t)
	// No query may be sent for malformed input; the cleanup check enforces that.
	if err := repo.AddHabitFromJSON([]byte(`{not json`)); err == nil {
		t.Error("expected an error for invalid JSON, got nil")
	}
}

func TestEditHabit(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectExec(updateHabit).
		WithArgs("yoga", "minutes", false, "calm", 3).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.EditHabit(Habit{ID: 3, Name: "yoga", Unit: "minutes", Pomodoro: false, Why: "calm"})
	if err != nil {
		t.Fatalf("EditHabit: %v", err)
	}
}

func TestEditHabitError(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectExec(updateHabit).WillReturnError(errDB)

	if err := repo.EditHabit(Habit{ID: 3}); !errors.Is(err, errDB) {
		t.Errorf("got error %v, want %v", err, errDB)
	}
}

func TestEditHabitFromJSON(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectExec(updateHabit).
		WithArgs("yoga", "minutes", true, "calm", 3).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := []byte(`{"id":3,"name":"yoga","unit":"minutes","pomodoro":true,"why":"calm"}`)
	if err := repo.EditHabitFromJSON(req); err != nil {
		t.Fatalf("EditHabitFromJSON: %v", err)
	}
}

func TestEditHabitFromJSONInvalid(t *testing.T) {
	repo, _ := newMockRepository(t)
	if err := repo.EditHabitFromJSON([]byte(`[]`)); err == nil {
		t.Error("expected an error for invalid JSON, got nil")
	}
}

func TestDeleteHabit(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectExec(deleteHabit).WithArgs(5).WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.DeleteHabit(5); err != nil {
		t.Fatalf("DeleteHabit: %v", err)
	}
}

func TestDeleteHabitError(t *testing.T) {
	repo, mock := newMockRepository(t)
	mock.ExpectExec(deleteHabit).WithArgs(5).WillReturnError(errDB)

	if err := repo.DeleteHabit(5); !errors.Is(err, errDB) {
		t.Errorf("got error %v, want %v", err, errDB)
	}
}
