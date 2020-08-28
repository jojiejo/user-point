package controllers

import (
	"encoding/json"
	"log"
	"net/http"
    "errors"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetMonthlySalesReport(c *gin.Context) {
	log.Printf("Begin => Get Monthly Sales Report Data")

	report_date			:= c.Param("report_date")
	report_data 		:= models.MonthlySalesReportResponse{}
	retrieved_data, err := report_data.MonthlySalesFile(report_date)

	// INPUT VALIDATION
    var error_messages = make(map[string]string)

    if len(report_date) != 6 {
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
	log.Printf("Get Monthly Sales Report Data : ", string(stringified_retrieved_data))

	c.JSON(http.StatusOK, retrieved_data)

	log.Printf("End => Get Monthly Sales Report Data")
}