package main

import (
	"encoding/json"
)

func JSONToHabit(req []byte) (Habit, error) {
	var h Habit
	err := json.Unmarshal(req, &h)
	return h, err
}

func HabitToJSON(habit Habit) ([]byte, error) {
	return json.Marshal(habit)
}

func JSONToDate(req []byte) (Date, error) {
	var d Date
	err := json.Unmarshal(req, &d)
	return d, err
}

func DateToJSON(date Date) ([]byte, error) {
	return json.Marshal(date)
}
