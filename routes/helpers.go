package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getRequiredHeaderID reads an int64 id from the given request header.
// TODO: replace with JWT auth + authorization.
func getRequiredHeaderID(context *gin.Context, headerName string) (id int64, ok bool) {
	raw := context.GetHeader(headerName)
	if raw == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": headerName + " header is required."})
		return 0, false
	}

	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse " + headerName + " header."})
		return 0, false
	}

	return id, true
}
