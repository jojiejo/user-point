package controllers

import (
	"encoding/json"
	"log"
	"net/http"
    "errors"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetDailyTransactionShellReport(c *gin.Context) {
	log.Printf("Begin => Get Daily Transaction Shell Report Data")

	report_date			:= c.Param("report_date")
	report_data 		:= models.DailyTransactionShellReportResponse{}
	retrieved_data, err := report_data.DailyTransactionShellFile(report_date)

	// INPUT VALIDATION
    var error_messages = make(map[string]string)

    if len(report_date) != 8 {
        err = errors.New("Report Date Format not Valid")
        error_messages["error_message"] = err.Error()
    }

	// INPUT HANDLER
	if len(error_messages) > 0 {
		log.Println(error_messages)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": error_messages,
		})
		return
	}

	stringified_retrieved_data, _ := json.Marshal(retrieved_data)
	log.Printf("Get Monthly Daily Transaction Shell Report Data : ", string(stringified_retrieved_data))

	c.JSON(http.StatusOK, retrieved_data)

	log.Printf("End => Get Daily Transaction Shell Report Data")
}