package main

import (
	_ "embed"

	"github.com/wailsapp/wails"
)

func main() {
	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "habit_tracker",
		Colour: "#131313",
	})
	app.Bind(NewMySQLConnection())
	app.Bind(NewHabit)
	app.Bind(NewDate)
	app.Bind(JSONToHabit)
	app.Run()
}
