package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/SeifeldenAtia/events-app/db"
)

var ErrEventNotFound = errors.New("event not found")
var ErrNotEnoughTickets = errors.New("not enough tickets available")

type Booking struct {
	ID       int64 `json:"id"`
	EventID  int64 `json:"eventId"`
	UserID   int64 `json:"userId"`
	Quantity int   `json:"quantity"`
}

type BookingWithEvent struct {
	ID        int64     `json:"id"`
	EventID   int64     `json:"eventId"`
	EventName string    `json:"eventName"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}

// CreateBooking atomically decrements available tickets and inserts the
// booking in one transaction to prevent overselling.
func CreateBooking(eventID, userID int64, quantity int) (*Booking, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	result, err := tx.Exec(
		`UPDATE events SET available_tickets = available_tickets - ? WHERE id = ? AND available_tickets >= ?`,
		quantity, eventID, quantity,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if rowsAffected == 0 {
		tx.Rollback()

		var exists bool
		err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM events WHERE id = ?)", eventID).Scan(&exists)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrEventNotFound
		}
		return nil, ErrNotEnoughTickets
	}

	insertResult, err := tx.Exec(
		`INSERT INTO bookings(event_id, user_id, quantity) VALUES (?, ?, ?)`,
		eventID, userID, quantity,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	bookingID, err := insertResult.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Booking{ID: bookingID, EventID: eventID, UserID: userID, Quantity: quantity}, nil
}

func GetBookingsByUser(userID int64) ([]BookingWithEvent, error) {
	query := `
	SELECT b.id, b.event_id, e.name, b.quantity, b.created_at
	FROM bookings b
	JOIN events e ON e.id = b.event_id
	WHERE b.user_id = ?
	`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := []BookingWithEvent{}
	for rows.Next() {
		var b BookingWithEvent
		if err := rows.Scan(&b.ID, &b.EventID, &b.EventName, &b.Quantity, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func CountBookingsForEvent(eventID int64) (int, error) {
	var total sql.NullInt64
	query := "SELECT SUM(quantity) FROM bookings WHERE event_id = ?"
	err := db.DB.QueryRow(query, eventID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return int(total.Int64), nil
}
