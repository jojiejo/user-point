package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"fleethub.shell.co.id/api/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//GenerateBearerCards => Generate Bearer Card
func (server *Server) GenerateBearerCards(c *gin.Context) {
	log.Printf("Begin => Generate Bearer Card")
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

	crb := models.CardRequestBearer{}
	err = json.Unmarshal(body, &crb)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	//Bank code, bank card number, country code, expiry date
	bankCode := os.Getenv("CARD_BANK_CODE")
	//bankCardNumber := os.Getenv("CARD_BANK_CARD_NUMBER")
	countryCode := os.Getenv("CARD_COUNTRY_CODE")
	//expiryDate := os.Getenv("CARD_EXPIRY_DATE")

	//Load Res Profile ID
	cg := models.CardGroup{}
	receivedCardGroup, err := cg.FindCardGroupByID(server.DB, crb.CardGroupID)
	if err != nil {
		log.Printf(err.Error())
		errList["no_card_group"] = "No card group found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	crb.ResProfileID = receivedCardGroup.ResProfileID

	//Load Card Prefix & Suffix
	ct := models.CardType{}
	receivedCardType, err := ct.FindCardTypeByID(server.DB, crb.CardTypeID)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			log.Printf(err.Error())
			errList["no_card_type"] = "No card type found"
			c.JSON(http.StatusNotFound, gin.H{
				"error": errList,
			})
			return
		}
	}

	crb.CardTypePrefix = receivedCardType.Prefix
	crb.CardTypeSuffix = receivedCardType.Code

	//Load Card Number
	type CardNumber struct {
		LastCardNo uint64 `json:"last_card_no"`
		NewCard    uint64 `json:"new_card"`
	}
	cardNumber := CardNumber{}
	err = server.DB.Debug().
		Table("mstCardNumber").
		Select("last_card_no").
		Where("card_type_id = ? AND country_code = ? AND bank_code = ?", crb.CardTypeID, countryCode, bankCode).
		Order("last_card_no desc").
		First(&cardNumber).
		Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			log.Printf(err.Error())
			errList["no_card_number_found"] = "No card number found"
			c.JSON(http.StatusNotFound, gin.H{
				"error": errList,
			})
			return
		}
	}

	//Card number conditioning
	if gorm.IsRecordNotFoundError(err) {
		cardNumber.NewCard = 1
	} else {
		cardNumber.NewCard = cardNumber.LastCardNo + 1
	}

	//Load Batch Inventory
	type BatchNumber struct {
		LastBatch uint64 `json:"last_batch"`
		NewBatch  uint64 `json:"new_batch"`
	}
	batchNumber := BatchNumber{}
	err = server.DB.Debug().
		Table("mstCardInventory").
		Select("last_batch").
		Where("cc_id = ?", crb.CCID).
		Order("last_batch desc").
		First(&batchNumber).
		Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			log.Printf(err.Error())
			errList["no_batch_found"] = "No batch found"
			c.JSON(http.StatusNotFound, gin.H{
				"error": errList,
			})
			return
		}
	}

	//Batch number conditioning
	if gorm.IsRecordNotFoundError(err) {
		batchNumber.NewBatch = 1
	} else {
		if crb.Batch == 0 {
			batchNumber.NewBatch = batchNumber.LastBatch + 1
		} else {
			batchNumber.NewBatch = crb.Batch
		}
	}

	/*	crb.Prepare()
		errorMessages := ac.Validate()
		if len(errorMessages) > 0 {
			log.Println(errorMessages)
			errList = errorMessages
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": errList,
			})
			return
		}

		createdAccountClass, err := ac.CreateAccountClass(server.DB)
		if err != nil {
			log.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"response": createdAccountClass,
		})*/

	log.Printf("End => Generate Bearer Card")
}

//GenerateVehicleCard => Generate Vehicle Card
func (server *Server) GenerateVehicleCard(c *gin.Context) {
	log.Printf("Begin => Generate Vehicle Card")
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

	ac := models.AccountClass{}
	err = json.Unmarshal(body, &ac)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	ac.Prepare()
	errorMessages := ac.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	createdAccountClass, err := ac.CreateAccountClass(server.DB)
	if err != nil {
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"response": createdAccountClass,
	})

	log.Printf("End => Generate Vehicle Class")
}
