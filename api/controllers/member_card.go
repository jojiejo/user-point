package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"fleethub.shell.co.id/api/models"
	"fleethub.shell.co.id/api/security"
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

	errorMessages := crb.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
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
	receivedCardGroup, err := cg.FindCardGroupByID(server.DB, uint64(crb.CardGroupID))
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
	receivedCardType, err := ct.FindCardTypeByID(server.DB, uint64(crb.CardTypeID))
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
		LastCardNo int `json:"last_card_no"`
	}
	cardNumber := CardNumber{}
	cardNumberErr := server.DB.Debug().
		Table("mstCardNumber").
		Select("last_card_no").
		Where("card_type_id = ? AND country_code = ? AND bank_code = ?",
			crb.CardTypeID,
			countryCode,
			bankCode).
		Order("last_card_no desc").
		First(&cardNumber).
		Error
	if cardNumberErr != nil {
		if !gorm.IsRecordNotFoundError(cardNumberErr) {
			log.Printf(cardNumberErr.Error())
			errList["no_card_number_found"] = "No card number found"
			c.JSON(http.StatusNotFound, gin.H{
				"error": errList,
			})
			return
		}
	}

	//Card number conditioning
	var nextCardNumber int
	if gorm.IsRecordNotFoundError(err) {
		nextCardNumber = 1
	} else {
		nextCardNumber = cardNumber.LastCardNo + 1
	}

	//Load Batch Inventory
	type BatchNumber struct {
		LastBatch int `json:"last_batch"`
	}
	batchNumber := BatchNumber{}
	batchNumberErr := server.DB.Debug().
		Table("mstCardInventory").
		Select("last_batch").
		Where("cc_id = ?", crb.CCID).
		Order("last_batch desc").
		First(&batchNumber).
		Error
	if batchNumberErr != nil {
		if !gorm.IsRecordNotFoundError(batchNumberErr) {
			log.Printf(batchNumberErr.Error())
			errList["no_batch_found"] = "No batch found"
			c.JSON(http.StatusNotFound, gin.H{
				"error": errList,
			})
			return
		}
	}

	//Batch number conditioning
	var nextBatch int
	if gorm.IsRecordNotFoundError(batchNumberErr) {
		nextBatch = 1
	} else {
		if crb.Batch == 0 {
			nextBatch = batchNumber.LastBatch + 1
		} else {
			nextBatch = crb.Batch
		}
	}

	cardInsertionTrx := server.DB.Begin()
	for i := 1; i <= crb.CardCount; i++ {
		cardID := crb.CardTypePrefix +
			crb.CardTypeSuffix +
			countryCode +
			bankCode +
			security.PadLeft(strconv.Itoa(nextCardNumber), "0", 8)

		validCard := security.GenerateLuhn(cardID)
		encryptedCard, _ := security.Encrypt(validCard)
		validCard = strings.ToUpper(security.Bin2hex(encryptedCard)[0:32])

		cvv := rand.Intn(999-100) + 100
		validCvv := strconv.Itoa(cvv)

		originalCardNumber := models.MemberCard{}
		err = server.DB.Debug().
			Model(models.MemberCard{}).
			Where("card_id = ?", validCard).
			Order("card_id desc").
			First(&originalCardNumber).
			Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				log.Printf(err.Error())
				errList["no_card"] = "No card found"
				c.JSON(http.StatusNotFound, gin.H{
					"error": errList,
				})
				return
			}
		}

		if gorm.IsRecordNotFoundError(err) {
			memberCard := models.MemberCard{
				CardID:           validCard,
				ExpDate:          crb.ExpDate,
				CVV:              validCvv,
				BankCode:         bankCode,
				CountryCode:      countryCode,
				Status:           "INACTIVE",
				Batch:            nextBatch,
				CardGroupID:      crb.CardGroupID,
				CardHolderTypeID: 1,
				CardTypeID:       crb.CardTypeID,
				CardProfileID:    crb.ResProfileID,
			}

			memberCard.Prepare()
			errorMessages := memberCard.Validate()
			if len(errorMessages) > 0 {
				log.Println(errorMessages)
				errList = errorMessages
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": errList,
				})
				return
			}

			_, err := memberCard.CreateMemberCard(server.DB)
			if err != nil {
				cardInsertionTrx.Rollback()
				log.Printf(err.Error())
				errList["card_generation"] = "Card generation failed"
				c.JSON(http.StatusNotFound, gin.H{
					"error": errList,
				})
				return
			}

			// Insert Perso
			cardPerso := models.CardPerso{
				CCID:   crb.CCID,
				CardID: validCard,
			}

			cardPerso.Prepare()
			errorMessages = cardPerso.Validate()
			if len(errorMessages) > 0 {
				log.Println(errorMessages)
				errList = errorMessages
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": errList,
				})
				return
			}

			_, err = cardPerso.CreateCardPerso(server.DB)
			if err != nil {
				cardInsertionTrx.Rollback()
				log.Printf(err.Error())
				errList["card_generation"] = "Card generation failed. Perso Problem."
				c.JSON(http.StatusNotFound, gin.H{
					"error": errList,
				})
				return
			}

			nextCardNumber++

		} else {
			nextCardNumber++
		}
	}

	cardNumberTracker := models.CardNumber{
		CardTypeID:  crb.CardTypeID,
		CountryCode: countryCode,
		BankCode:    bankCode,
		LastCardNo:  nextCardNumber - 1,
	}

	if gorm.IsRecordNotFoundError(cardNumberErr) {
		_, err := cardNumberTracker.CreateCardNumber(server.DB)
		if err != nil {
			cardInsertionTrx.Rollback()
			log.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	} else {
		_, err := cardNumberTracker.UpdateCardNumber(server.DB)
		if err != nil {
			cardInsertionTrx.Rollback()
			log.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	}

	cardBatchNumberTracker := models.CardBatchNumber{
		CCID:      crb.CCID,
		LastBatch: nextBatch,
	}

	if gorm.IsRecordNotFoundError(batchNumberErr) {
		_, err := cardBatchNumberTracker.CreateCardBatchNumber(server.DB)
		if err != nil {
			cardInsertionTrx.Rollback()
			log.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	} else {
		_, err := cardBatchNumberTracker.UpdateCardBatchNumber(server.DB)
		if err != nil {
			cardInsertionTrx.Rollback()
			log.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	}

	cardInsertionTrx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"response": "Bearer cards have been generated",
	})
	log.Printf("End => Generate Bearer Card")
}

//GenerateDriverCards => Generate Driver Cards
func (server *Server) GenerateDriverCards(c *gin.Context) {
	log.Printf("Begin => Generate Driver Cards")
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

	crd := models.CardRequestDriver{}
	err = json.Unmarshal(body, &crd)
	if err != nil {
		log.Printf(err.Error())
		errList["unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errList,
		})
		return
	}

	crd.Prepare()
	errorMessages := crd.Validate()
	if len(errorMessages) > 0 {
		log.Println(errorMessages)
		errList = errorMessages
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
	receivedCardGroup, err := cg.FindCardGroupByID(server.DB, uint64(crd.CardGroupID))
	if err != nil {
		log.Printf(err.Error())
		errList["no_card_group"] = "No card group found"
		c.JSON(http.StatusNotFound, gin.H{
			"error": errList,
		})
		return
	}

	crd.ResProfileID = receivedCardGroup.ResProfileID

	//Load Card Prefix & Suffix
	ct := models.CardType{}
	receivedCardType, err := ct.FindCardTypeByID(server.DB, uint64(crd.CardTypeID))
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

	crd.CardTypePrefix = receivedCardType.Prefix
	crd.CardTypeSuffix = receivedCardType.Code

	//Load Card Number
	type CardNumber struct {
		LastCardNo int `json:"last_card_no"`
	}
	cardNumber := CardNumber{}
	cardNumberErr := server.DB.Debug().
		Table("mstCardNumber").
		Select("last_card_no").
		Where("card_type_id = ? AND country_code = ? AND bank_code = ?",
			crd.CardTypeID,
			countryCode,
			bankCode).
		Order("last_card_no desc").
		First(&cardNumber).
		Error
	if cardNumberErr != nil {
		if !gorm.IsRecordNotFoundError(cardNumberErr) {
			log.Printf(cardNumberErr.Error())
			errList["no_card_number_found"] = "No card number found"
			c.JSON(http.StatusNotFound, gin.H{
				"error": errList,
			})
			return
		}
	}

	//Card number conditioning
	var nextCardNumber int
	if gorm.IsRecordNotFoundError(err) {
		nextCardNumber = 1
	} else {
		nextCardNumber = cardNumber.LastCardNo + 1
	}

	//Load Batch Inventory
	type BatchNumber struct {
		LastBatch int `json:"last_batch"`
	}
	batchNumber := BatchNumber{}
	batchNumberErr := server.DB.Debug().
		Table("mstCardInventory").
		Select("last_batch").
		Where("cc_id = ?", crd.CCID).
		Order("last_batch desc").
		First(&batchNumber).
		Error
	if batchNumberErr != nil {
		if !gorm.IsRecordNotFoundError(batchNumberErr) {
			log.Printf(batchNumberErr.Error())
			errList["no_batch_found"] = "No batch found"
			c.JSON(http.StatusNotFound, gin.H{
				"error": errList,
			})
			return
		}
	}

	//Batch number conditioning
	var nextBatch int
	if gorm.IsRecordNotFoundError(batchNumberErr) {
		nextBatch = 1
	} else {
		if crd.Batch == 0 {
			nextBatch = batchNumber.LastBatch + 1
		} else {
			nextBatch = crd.Batch
		}
	}

	cardInsertionTrx := server.DB.Begin()
	for i := range crd.Drivers {
		cardID := crd.CardTypePrefix +
			crd.CardTypeSuffix +
			countryCode +
			bankCode +
			security.PadLeft(strconv.Itoa(nextCardNumber), "0", 8)

		validCard := security.GenerateLuhn(cardID)
		encryptedCard, _ := security.Encrypt(validCard)
		validCard = strings.ToUpper(security.Bin2hex(encryptedCard)[0:32])

		cvv := rand.Intn(999-100) + 100
		validCvv := strconv.Itoa(cvv)

		originalCardNumber := models.MemberCard{}
		err = server.DB.Debug().
			Model(models.MemberCard{}).
			Where("card_id = ?", validCard).
			Order("card_id desc").
			First(&originalCardNumber).
			Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				log.Printf(err.Error())
				errList["no_card"] = "No card found"
				c.JSON(http.StatusNotFound, gin.H{
					"error": errList,
				})
				return
			}
		}

		if gorm.IsRecordNotFoundError(err) {
			memberCard := models.MemberCard{
				CardID:           validCard,
				ExpDate:          crd.ExpDate,
				CVV:              validCvv,
				BankCode:         bankCode,
				CountryCode:      countryCode,
				Status:           "INACTIVE",
				Batch:            nextBatch,
				CardGroupID:      crd.CardGroupID,
				CardHolderTypeID: 2,
				CardTypeID:       crd.CardTypeID,
				CardProfileID:    crd.ResProfileID,
			}

			memberCard.Prepare()
			errorMessages = memberCard.Validate()
			if len(errorMessages) > 0 {
				log.Println(errorMessages)
				errList = errorMessages
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": errList,
				})
				return
			}

			_, err := memberCard.CreateMemberCard(server.DB)
			if err != nil {
				cardInsertionTrx.Rollback()
				log.Printf(err.Error())
				errList["card_generation"] = "Card generation failed"
				c.JSON(http.StatusNotFound, gin.H{
					"error": errList,
				})
				return
			}

			// Insert to Driver to Member Card Relation
			cardDriverRelation := models.CardDriverRelation{
				DriverID: crd.Drivers[i].ID,
				CardID:   validCard,
			}

			cardDriverRelation.Prepare()
			errorMessages = cardDriverRelation.Validate()
			if len(errorMessages) > 0 {
				log.Println(errorMessages)
				errList = errorMessages
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": errList,
				})
				return
			}

			_, err = cardDriverRelation.CreateCardDriverRelation(server.DB)
			if err != nil {
				cardInsertionTrx.Rollback()
				log.Printf(err.Error())
				errList["card_generation"] = "Card generation failed. Driver relation Problem."
				c.JSON(http.StatusNotFound, gin.H{
					"error": errList,
				})
				return
			}

			// Insert Perso
			cardPerso := models.CardPerso{
				CCID:   crd.CCID,
				CardID: validCard,
			}

			cardPerso.Prepare()
			errorMessages = cardPerso.Validate()
			if len(errorMessages) > 0 {
				log.Println(errorMessages)
				errList = errorMessages
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": errList,
				})
				return
			}

			_, err = cardPerso.CreateCardPerso(server.DB)
			if err != nil {
				cardInsertionTrx.Rollback()
				log.Printf(err.Error())
				errList["card_generation"] = "Card generation failed. Perso Problem."
				c.JSON(http.StatusNotFound, gin.H{
					"error": errList,
				})
				return
			}

			nextCardNumber++

		} else {
			nextCardNumber++
		}
	}

	cardNumberTracker := models.CardNumber{
		CardTypeID:  crd.CardTypeID,
		CountryCode: countryCode,
		BankCode:    bankCode,
		LastCardNo:  nextCardNumber - 1,
	}

	if gorm.IsRecordNotFoundError(cardNumberErr) {
		_, err := cardNumberTracker.CreateCardNumber(server.DB)
		if err != nil {
			cardInsertionTrx.Rollback()
			log.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	} else {
		_, err := cardNumberTracker.UpdateCardNumber(server.DB)
		if err != nil {
			cardInsertionTrx.Rollback()
			log.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	}

	cardBatchNumberTracker := models.CardBatchNumber{
		CCID:      crd.CCID,
		LastBatch: nextBatch,
	}

	if gorm.IsRecordNotFoundError(batchNumberErr) {
		_, err := cardBatchNumberTracker.CreateCardBatchNumber(server.DB)
		if err != nil {
			cardInsertionTrx.Rollback()
			log.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}
	}

	_, err = cardBatchNumberTracker.UpdateCardBatchNumber(server.DB)
	if err != nil {
		cardInsertionTrx.Rollback()
		log.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	cardInsertionTrx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"response": "Driver cards have been generated",
	})

	log.Println("End => Generate Driver Cards")
}
