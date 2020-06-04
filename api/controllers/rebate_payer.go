package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetRebatePayerRelations(c *gin.Context) {
	log.Printf("Begin => Get Rebate Payer Relations")

	rebatePayerRelation := models.RebatePayer{}
	rebatePayerRelations, err := rebatePayerRelation.FindRebatePayerRelations(server.DB)
	if err != nil {
		log.Printf(err.Error())
		errList["no_rebate_payer"] = "No rebate payer relation found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedRebatePayerRelations, _ := json.Marshal(rebatePayerRelations)
	log.Printf("Get Rebate Payer Relations : ", string(stringifiedRebatePayerRelations))
	c.JSON(http.StatusOK, gin.H{
		"response": rebatePayerRelations,
	})

	log.Printf("End => Get Rebate Payer Relations")
}

func (server *Server) GetRebatePayerRelationByID(c *gin.Context) {
	log.Printf("Begin => Get Rebate Payer Relation By ID")
	relationID := c.Param("id")
	convertedRelationID, err := strconv.ParseUint(relationID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	rebatePayerRelation := models.RebatePayer{}
	rebatePayerRelationReceived, err := rebatePayerRelation.FindRebatePayerRelationByID(server.DB, convertedRelationID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_relation"] = "No relation found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	stringifiedReceivedRebatePayerRelation, _ := json.Marshal(rebatePayerRelationReceived)
	log.Printf("Get Rebate Payer Relation By ID : ", string(stringifiedReceivedRebatePayerRelation))
	c.JSON(http.StatusOK, gin.H{
		"response": rebatePayerRelationReceived,
	})

	log.Printf("End => Get Rebate Payer Relation By ID")
}

func (server *Server) CreateMainRebatePayer(c *gin.Context) {
	log.Printf("Begin => Create Rebate Payer Relation")
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	prp := models.PostedRebatePayer{}
	err = json.Unmarshal(body, &prp)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	//Count whether there is still a relation active
	type ActiveRelation struct {
		RppCount int `json:"rpp_count"`
	}

	activeRelation := ActiveRelation{}
	var numberOfRelation = 0
	if len(prp.CCID) > 0 {
		for i, _ := range prp.CCID {
			err = server.DB.Debug().Raw("EXEC spAPI_RebatePayer_CountActive ?, ?, ?", prp.CCID[i], "Main", prp.StartedAt).Scan(&activeRelation).Error
			if err != nil {
				log.Printf(err.Error())
				errList["unmarshal_error"] = "Cannot unmarshal body"
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": errList,
				})
				return
			}

			if activeRelation.RppCount > 0 {
				numberOfRelation++
			}
		}
	}

	if numberOfRelation > 0 {
		errList["linked_retailer"] = strconv.Itoa(numberOfRelation) + " of selected account still have an active main rebate."
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	prp.Prepare()
	errorMessages := prp.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	prpCreated, errorMessages := prp.CreateRebatePayerRelation(server.DB)
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		c.JSON(http.StatusNotFound, gin.H{
			"error": errorMessages,
		})
		return
	} else {
		stringifiedPrpCreated, _ := json.Marshal(prpCreated)
		log.Printf("Get Main Rebate Payer: ", string(stringifiedPrpCreated))
		c.JSON(http.StatusOK, gin.H{
			"response": prpCreated,
		})
	}

	log.Printf("End => Create Rebate Payer Relation")
}

func (server *Server) CreatePromotionalRebatePayer(c *gin.Context) {
	log.Printf("Begin => Create Rebate Payer Relation")
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	prp := models.PostedRebatePayer{}
	err = json.Unmarshal(body, &prp)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	//Count whether there is still a relation active
	type ActiveRelation struct {
		RppCount int `json:"rpp_count"`
	}

	activeRelation := ActiveRelation{}
	var numberOfRelation = 0
	if len(prp.CCID) > 0 {
		for i, _ := range prp.CCID {
			err = server.DB.Debug().Raw("EXEC spAPI_RebatePayer_CountActive ?, ?, ?", prp.CCID[i], "Promotional", prp.StartedAt).Scan(&activeRelation).Error
			if err != nil {
				log.Printf(err.Error())
				errList["unmarshal_error"] = "Cannot unmarshal body"
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": errList,
				})
				return
			}

			if activeRelation.RppCount > 0 {
				numberOfRelation++
			}
		}
	}

	if numberOfRelation > 1 {
		errList["linked_retailer"] = strconv.Itoa(numberOfRelation) + " of selected account still have an active promotional rebate."
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	prp.Prepare()
	errorMessages := prp.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	prpCreated, errorMessages := prp.CreateRebatePayerRelation(server.DB)
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		c.JSON(http.StatusNotFound, gin.H{
			"error": errorMessages,
		})
		return
	} else {
		stringifiedPrpCreated, _ := json.Marshal(prpCreated)
		log.Printf("Get Promotional Rebate Payer : ", string(stringifiedPrpCreated))
		c.JSON(http.StatusOK, gin.H{
			"response": prpCreated,
		})
	}

	log.Printf("End => Create Rebate Payer Relation")
}

func (server *Server) UpdateRebatePayerRelation(c *gin.Context) {
	log.Printf("Begin => Update Rebate Payer Relation")

	errList = map[string]string{}
	relationID := c.Param("id")
	relationid, err := strconv.ParseUint(relationID, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_request"] = "Invalid request"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errList,
		})
		return
	}

	originalRelation := models.RebatePayer{}
	err = server.DB.Debug().Model(models.RebatePayer{}).Where("id = ?", relationid).Order("id desc").Take(&originalRelation).Error
	if err != nil {
		log.Printf(err.Error())
		errList["no_relation_found"] = "No relation found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	relation := models.RebatePayer{}
	err = json.Unmarshal(body, &relation)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	relation.ID = originalRelation.ID
	relationUpdated, err := relation.UpdateRebatePayerRelation(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": relationUpdated,
	})

	log.Printf("End => Update Update Rebate Payer Relation")
}

func (server *Server) CheckBulkAssignRebateToPayer(c *gin.Context) {
	log.Printf("Begin => Check Bulk Assign Rebate to Payer")
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	bulkCheckRebate := models.BulkAssignRebate{}
	err = json.Unmarshal(body, &bulkCheckRebate)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	bulkCheckRebate.Prepare()
	errorMessages := bulkCheckRebate.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	checkedField, errorMessages := bulkCheckRebate.BulkCheckAssignRebate(server.DB)
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errorMessages,
		})
		return
	} else {
		stringifiedCheckedField, _ := json.Marshal(checkedField)
		log.Printf("Get Bulk Assign Rebate to Payer : ", string(stringifiedCheckedField))
		c.JSON(http.StatusOK, gin.H{
			"response": checkedField,
		})
	}

	log.Printf("End => Check Bulk Assign Rebate to Payer")
}

func (server *Server) BulkAssignRebateToPayer(c *gin.Context) {
	log.Printf("Begin => Bulk Assign Rebate To Payer")
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf(err.Error())
		errList["invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	bulkAssignRebate := models.BulkAssignRebate{}
	err = json.Unmarshal(body, &bulkAssignRebate)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	bulkAssignRebate.Prepare()
	errorMessages := bulkAssignRebate.ChargeValidate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	bulkAssignField, errorMessages := bulkAssignRebate.BulkAssignRebate(server.DB)
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		c.JSON(http.StatusNotFound, gin.H{
			"error": errorMessages,
		})
		return
	} else {
		stringifiedBulkAssignField, _ := json.Marshal(bulkAssignField)
		log.Printf("Get Bulk Assign Rebate To Payer : ", string(stringifiedBulkAssignField))
		c.JSON(http.StatusOK, gin.H{
			"response": bulkAssignField,
		})
	}

	log.Printf("End => Bulk Assign Rebate To Payer")
}
