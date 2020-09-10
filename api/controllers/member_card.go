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
			log.Printf(err.Error())
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
	var generatedCards []string
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

		err = server.DB.Debug().
			Table("mstMemberCard").
			Where("card_id = ?", encryptedCard).
			Order("card_id desc").
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
			convertedBankCode, _ := strconv.Atoi(bankCode)
			convertedCountryCode, _ := strconv.Atoi(countryCode)
			memberCard := models.MemberCard{
				CardID:           validCard,
				ExpDate:          crb.ExpDate,
				CVV:              validCvv,
				BankCode:         convertedBankCode,
				CountryCode:      convertedCountryCode,
				Status:           "INACTIVE",
				Batch:            nextBatch,
				CardGroupID:      crb.CardGroupID,
				CardHolderTypeID: 1,
				CardTypeID:       crb.CardTypeID,
				CardProfileID:    crb.ResProfileID,
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

			//If save card error => DB rollback, break
			generatedCards = append(generatedCards, cardID)
			nextCardNumber++

		} else {
			nextCardNumber++
		}
	}

	cardInsertionTrx.Commit()

	if gorm.IsRecordNotFoundError(cardNumberErr) {
		//Make new card number data with latest nextCardNumber -1
	} else {
		//Update card number with latest nextCardNumber
	}

	if gorm.IsRecordNotFoundError(batchNumberErr) {
		//Make new batch data with latest nextCardNumber -1
	} else {
		//Update batch with latest nextCardNumber
	}

	log.Printf("End => Generate Bearer Card")
}
