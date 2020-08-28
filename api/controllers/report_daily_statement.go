package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetDailyStatementReport(c *gin.Context) {
	log.Printf("Begin => Get Daily Statement Report Data")

	reportDate 		:= c.Param("report_date")
	soldToNumber	:= c.Param("sold_to_number")
	shipToNumber	:= c.Param("ship_to_number")

	// INPUT VALIDATION
    if len(reportDate) != 8 {
		errList["date_format"] = "Report Date Format not Valid"
        c.JSON(http.StatusOK, gin.H{
        	"error" : errList,
        })
        log.Println(errList)
        return
    }

    reportData 			:= models.DailyStatementReportResponse{}
    reportDataParam 	:= models.DailyStatementReportFileParam{
    	ReportFulldate	: reportDate,
    	SoldToNumber 	: soldToNumber,
    	ShipToNumber 	: shipToNumber,
    	FtpConn			: server.FtpConn,
    }

    retrievedData, err := reportData.DailyStatementFile(reportDataParam)

	// INPUT HANDLER
	if err != nil {
		errList["ftp_error"] = err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
        log.Println(errList)
		return
	}

	stringifiedRetrievedData, _ := json.Marshal(retrievedData)
	log.Printf("Get Daily Statement Report Data : ", string(stringifiedRetrievedData))

	c.JSON(http.StatusOK, retrievedData)

	log.Printf("End => Get Daily Statement Report Data")
}