package main

import (
	"database/sql"
)

type HabitRepository interface {
	GetAllHabits() ([]Habit, error)
	GetHabit(uint) (Habit, error)
	AddHabit(Habit) error
	EditHabit(Habit) error
	DeleteHabit(uint) error
}

type MySQLRepository struct {
	DB *sql.DB
}

func NewMySQLConnection() *MySQLRepository {
	s := MySQLRepository{DB: DBConnect()}
	return &s
}

const habitColumns = "habit_id, habit_name, habit_unit, habit_pomodoro, habit_why"

// GetAllHabits retrieves the list of Habits.
func (s MySQLRepository) GetAllHabits() ([]Habit, error) {
	results, err := s.DB.Query("SELECT " + habitColumns + " FROM habit ORDER BY habit_id")
	if err != nil {
		return nil, err
	}
	defer results.Close()
	habits := []Habit{}
	for results.Next() {
		var habit Habit
		if err := results.Scan(&habit.ID, &habit.Name, &habit.Unit, &habit.Pomodoro, &habit.Why); err != nil {
			return nil, err
		}
		habits = append(habits, habit)
	}
	return habits, results.Err()
}

// GetHabit gets a habit with the given ID.
func (s MySQLRepository) GetHabit(id uint) (Habit, error) {
	var result Habit
	err := s.DB.QueryRow("SELECT "+habitColumns+" FROM habit WHERE habit_id = ?", id).
		Scan(&result.ID, &result.Name, &result.Unit, &result.Pomodoro, &result.Why)
	return result, err
}

// AddHabitFromJSON adds a habit given a JSON value of a habit object
func (s *MySQLRepository) AddHabitFromJSON(req []byte) error {
	habit, err := JSONToHabit(req)
	if err != nil {
		return err
	}
	return s.AddHabit(habit)
}

// AddHabit adds a new habit to the database. The ID is assigned by the database.
func (s MySQLRepository) AddHabit(habit Habit) error {
	_, err := s.DB.Exec(
		"INSERT INTO habit (habit_name, habit_unit, habit_pomodoro, habit_why) VALUES (?, ?, ?, ?)",
		habit.Name, habit.Unit, habit.Pomodoro, habit.Why,
	)
	return err
}

func (s MySQLRepository) EditHabitFromJSON(req []byte) error {
	habit, err := JSONToHabit(req)
	if err != nil {
		return err
	}
	return s.EditHabit(habit)
}

func (s MySQLRepository) EditHabit(habit Habit) error {
	_, err := s.DB.Exec(
		"UPDATE habit SET habit_name = ?, habit_unit = ?, habit_pomodoro = ?, habit_why = ? WHERE habit_id = ?",
		habit.Name, habit.Unit, habit.Pomodoro, habit.Why, habit.ID,
	)
	return err
}

// DeleteHabit deletes a habit. Its logged dates are removed by the
// ON DELETE CASCADE foreign key, see schema.sql.
func (s MySQLRepository) DeleteHabit(id uint) error {
	_, err := s.DB.Exec("DELETE FROM habit WHERE habit_id = ?", id)
	return err
}
