package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SeifeldenAtia/events-app/models"
	"github.com/gin-gonic/gin"
)

// TODO: add pagination for events list.
func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events. Try again later."})
		return
	}
	context.JSON(http.StatusOK, events)
}

func getEvent(context *gin.Context) {
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
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		return
	}

	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {
	ownerID, ok := getRequiredHeaderID(context, "X-Owner-Id")
	if !ok {
		return
	}

	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data.", "error": err.Error()})
		return
	}

	event.OwnerID = ownerID

	err = event.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create event. Try again later."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Event created!", "event": event})
}

func updateEvent(context *gin.Context) {
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

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data.", "error": err.Error()})
		return
	}

	updatedEvent.ID = eventId
	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event updated successfully!"})
}

func deleteEvent(context *gin.Context) {
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

	bookedCount, err := models.CountBookingsForEvent(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not check existing bookings for the event."})
		return
	}
	if bookedCount > 0 {
		context.JSON(http.StatusConflict, gin.H{"message": "Cannot delete an event that already has tickets booked."})
		return
	}

	err = event.Delete()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the event."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully!"})
}
