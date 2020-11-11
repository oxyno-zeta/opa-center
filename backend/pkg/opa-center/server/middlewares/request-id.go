package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/gofrs/uuid"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

type contextKey struct {
	name string
}

var reqCtxKey = &contextKey{name: "request-id"}

const requestIDHeader = "X-Request-Id"
const requestIDContextKey = "RequestID"

func RequestID(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request id from request
		requestID := c.Request.Header.Get(requestIDHeader)

		// Check if request id exists
		if requestID == "" {
			// Generate uuid
			uuid, err := uuid.NewV4()
			if err != nil {
				// Log error
				logger.Errorln(err)
				// Send response
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

				return
			}
			// Save it in variable
			requestID = uuid.String()
		}

		// Store it in context
		c.Set(requestIDContextKey, requestID)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), reqCtxKey, requestID))

		// Put it on header
		c.Writer.Header().Set(requestIDHeader, requestID)

		// Next
		c.Next()
	}
}

func GetRequestIDFromGin(c *gin.Context) string {
	requestIDObj, requestIDExists := c.Get(requestIDContextKey)
	if requestIDExists {
		// return request id
		return requestIDObj.(string)
	}

	return ""
}

func GetRequestIDFromContext(ctx context.Context) string {
	requestIDObj := ctx.Value(reqCtxKey)
	if requestIDObj != nil {
		return requestIDObj.(string)
	}

	return ""
}
