package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine) {
	v1 := server.Group("/api/v1")

	v1.GET("/events", getEvents)
	v1.GET("/events/:id", getEvent)
	v1.POST("/events", createEvent)
	v1.PUT("/events/:id", updateEvent)
	v1.DELETE("/events/:id", deleteEvent)

	v1.POST("/events/:id/bookings", createBooking)
	v1.GET("/events/:id/bookings/count", getEventBookingsCount)
	v1.GET("/bookings", getMyBookings)

	// TODO: support versioning via content negotiation for breaking changes.
}
