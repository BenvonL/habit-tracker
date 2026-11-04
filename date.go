package main

import (
	"errors"
	"time"
)

// dateLayout is the format used for dates exchanged with the frontend and
// stored in the date_date column.
const dateLayout = "2006-01-02"

type Date struct {
	ID      uint   `json:"id"`
	Date    string `json:"date"`
	Count   uint   `json:"count"`
	HabitID uint   `json:"habit_id"`
}

func NewDate(Today string, Count uint, HabitID uint) Date {
	return Date{Date: Today, Count: Count, HabitID: HabitID}
}

// Validate checks that a date entry can be stored.
func (d Date) Validate() error {
	if d.HabitID == 0 {
		return errors.New("missing habit id")
	}
	if d.Count == 0 {
		return errors.New("count must be greater than 0")
	}
	if _, err := time.Parse(dateLayout, d.Date); err != nil {
		return errors.New("date must be in the format YYYY-MM-DD")
	}
	return nil
}
