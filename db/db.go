package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "api.db")

	if err != nil {
		panic("Could not connect to database.")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		available_tickets INTEGER NOT NULL DEFAULT 0,
		owner_id INTEGER
	)
	`

	_, err := DB.Exec(createEventsTable)

	if err != nil {
		panic("Could not create events table.")
	}

	createBookingsTable := `
	CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER NOT NULL REFERENCES events(id),
		user_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`

	_, err = DB.Exec(createBookingsTable)

	if err != nil {
		panic("Could not create bookings table.")
	}
}
