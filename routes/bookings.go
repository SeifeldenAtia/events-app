package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SeifeldenAtia/events-app/models"
	"github.com/gin-gonic/gin"
)

type bookTicketInput struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

func createBooking(context *gin.Context) {
	userID, ok := getRequiredHeaderID(context, "X-User-Id")
	if !ok {
		return
	}

	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	var input bookTicketInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data.", "error": err.Error()})
		return
	}

	booking, err := models.CreateBooking(eventId, userID, input.Quantity)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEventNotFound):
			context.JSON(http.StatusNotFound, gin.H{"message": "Event not found."})
		case errors.Is(err, models.ErrNotEnoughTickets):
			context.JSON(http.StatusConflict, gin.H{"message": "Not enough tickets available."})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create booking. Try again later."})
		}
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Booking created!", "booking": booking})
}

func getMyBookings(context *gin.Context) {
	userID, ok := getRequiredHeaderID(context, "X-User-Id")
	if !ok {
		return
	}

	bookings, err := models.GetBookingsByUser(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch bookings. Try again later."})
		return
	}

	context.JSON(http.StatusOK, bookings)
}

func getEventBookingsCount(context *gin.Context) {
	ownerID, ok := getRequiredHeaderID(context, "X-Owner-Id")
	if !ok {
		return
	}

	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		if errors.Is(err, models.ErrEventNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": "Event not found."})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the event."})
		return
	}

	if event.OwnerID != ownerID {
		context.JSON(http.StatusForbidden, gin.H{"message": "You do not own this event."})
		return
	}

	count, err := models.CountBookingsForEvent(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not count bookings. Try again later."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"eventId": eventId, "ticketsBooked": count})
}
