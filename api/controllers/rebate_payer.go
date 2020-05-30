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
