package main

import (
	"database/sql"
	"time"
)

type DateRepository interface {
	GetAllDates(uint) ([]Date, error)
	GetTodaysCountForHabit(uint) (Date, error)
	AddCount(Date) error
}

// GetAllDates gets all dates for a given habit.
func (s MySQLRepository) GetAllDates(habitID uint) ([]Date, error) {
	results, err := s.DB.Query("SELECT date_id, date_date, date_count, habit_id FROM date WHERE habit_id = ? ORDER BY date_date", habitID)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	dates := []Date{}
	for results.Next() {
		var date Date
		if err := results.Scan(&date.ID, &date.Date, &date.Count, &date.HabitID); err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}
	return dates, results.Err()
}

// GetTodaysCountForHabit gets the entry for today (local time) on a given habit.
// If nothing was logged today, a Date with Count 0 is returned.
func (s MySQLRepository) GetTodaysCountForHabit(habitID uint) (Date, error) {
	date := NewDate(time.Now().Format(dateLayout), 0, habitID)
	err := s.DB.QueryRow("SELECT date_id, date_date, date_count, habit_id FROM date WHERE date_date = ? AND habit_id = ?", date.Date, habitID).
		Scan(&date.ID, &date.Date, &date.Count, &date.HabitID)
	if err == sql.ErrNoRows {
		return date, nil
	}
	return date, err
}

// AddCount adds date.Count to the entry for the given habit and day, creating
// the entry if it does not exist yet. This is a single atomic statement and
// relies on the unique key (habit_id, date_date), see schema.sql.
func (s MySQLRepository) AddCount(date Date) error {
	if err := date.Validate(); err != nil {
		return err
	}
	_, err := s.DB.Exec(
		"INSERT INTO date (date_date, date_count, habit_id) VALUES (?, ?, ?) "+
			"ON DUPLICATE KEY UPDATE date_count = date_count + VALUES(date_count)",
		date.Date, date.Count, date.HabitID,
	)
	return err
}

func (s MySQLRepository) AddCountFromJSON(req []byte) error {
	date, err := JSONToDate(req)
	if err != nil {
		return err
	}
	return s.AddCount(date)
}
