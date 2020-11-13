package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/business"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/log"
)

func AddStatusEndpoints(router gin.IRouter, busiServices *business.Services) {
	router.POST("/api/status/:partitionid", func(c *gin.Context) {
		// Get logger from request
		logger := log.GetLoggerFromGin(c)
		// Get partition id
		partitionID := c.Param("partitionid")
		// Read all input
		bb, err := ioutil.ReadAll(c.Request.Body)
		// Check error
		if err != nil {
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})

			return
		}

		// Create data for input parsing
		var mm map[string]interface{}

		err = json.Unmarshal(bb, &mm)
		// Check error
		if err != nil {
			logger.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})

			return
		}

		// Call service
		err = busiServices.StatusSvc.UnsecureCreate(partitionID, mm)
		// Check error
		if err != nil {
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})

			return
		}

		// Answer ok
		c.JSON(http.StatusOK, gin.H{"answer": "ok"})
	})
}
