package security

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/andreburgaud/crypt2go/ecb"
)

//Encrypt => Encrypt given string using AES ECB
func Encrypt(plainText, key string) ([]byte, error) {
	packedKey, err := hex.DecodeString(os.Getenv("CARD_KEY"))
	if err != nil {
		return nil, err
	}

	plainText = PadRight(plainText, "F", 32)
	packedPlainText, err := hex.DecodeString(plainText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(packedKey)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBEncrypter(block)
	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, packedPlainText)

	return cipherText, nil
}

//Decrypt => Decrypt given string using AES ECB
func decrypt(cipherText, key string) (string, error) {
	packedKey, err := hex.DecodeString(os.Getenv("CARD_KEY"))
	if err != nil {
		return "", err
	}

	packedCipherText, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(packedKey)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBDecrypter(block)
	plainText := make([]byte, len(packedCipherText))
	mode.CryptBlocks(plainText, packedCipherText)
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	replacedPad := reg.ReplaceAllString(Bin2hex(plainText), "")
	replacedLuhn := fmt.Sprintf("%16.16s", replacedPad)

	return replacedLuhn, nil
}

//PadRight => Give padding on the right side
func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

//PadLeft => Give padding on the left side
func PadLeft(str, pad string, lenght int) string {
	for {
		str = pad + str
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

//Bin2hex => Convert bin value to hex
func Bin2hex(str []byte) string {
	return hex.EncodeToString([]byte(str))
}

//GenerateLuhn => Generate Luhn
func GenerateLuhn(cardNumber string) string {
	stack := 0
	digits := strings.Split(reverseString(cardNumber), "")

	for key, value := range digits {
		loopValue, _ := strconv.Atoi(value)

		if key%2 == 0 {
			for _, secondValue := range strings.Split(strconv.Itoa(loopValue*2), "") {
				value += secondValue
			}
		}

		stack += loopValue
	}

	if stack %= 10; stack != 0 {
		stack -= 10
	}

	stringifiedAbsoluteNumber := fmt.Sprintf("%.0f", math.Abs(float64(stack)))
	return cardNumber + stringifiedAbsoluteNumber
}

func reverseString(text string) string {
	rns := []rune(text) // convert to rune
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}

	return string(rns)
}
