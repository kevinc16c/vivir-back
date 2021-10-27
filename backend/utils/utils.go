package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/skip2/go-qrcode"
)

func ParseInt(str string) uint {
	id, _ := strconv.ParseUint(str, 10, 32)

	return uint(id)
}

// func GetNow() time.Time {
// 	location, err := time.LoadLocation("America/Buenos_Aires")
// 	if err != nil {
// 		log.Print(err)
// 	}

// 	now := time.Now()
// 	t := now.In(location)

// 	log.Print(t)

// 	return t
// }

func GetFormatoImagen(imagen string) string {

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imagen))
	_, formatString, err := image.Decode(reader)
	if err != nil {
		return "error"
	}
	return formatString
}

func SavePng(imagen string, nombreImagen string, urlImagen string) string {
	dec, err := base64.StdEncoding.DecodeString(imagen)
	if err != nil {
		return "error de imagen"
	}

	r := bytes.NewReader(dec)
	im, err := png.Decode(r)
	if err != nil {
		return "Bad png " + err.Error()
	}

	f, err := os.OpenFile(urlImagen, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return "No se puede abrir el archivo " + err.Error()
	}

	png.Encode(f, im)

	return "ok"
}

func SaveJpg(imagen string, nombreImagen string, urlImagen string) string {
	dec, err := base64.StdEncoding.DecodeString(imagen)
	if err != nil {
		return "error de imagen"
	}

	r := bytes.NewReader(dec)
	im, err := jpeg.Decode(r)
	if err != nil {
		return "Bad jpg " + err.Error()
	}

	f, err := os.OpenFile(urlImagen, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return "No se puede abrir el archivo " + err.Error()
	}

	jpeg.Encode(f, im, &jpeg.Options{Quality: 100})

	return "ok"
}

func CreateQr(content string, size int, urlQr string) {
	err := qrcode.WriteFile(content, qrcode.Medium, size, urlQr)
	if err != nil {
		log.Print(err)
	}
}

func Encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	//key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher([]byte(keyString))
	if err != nil {
		log.Print(fmt.Sprintf("%v", err))
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Print(fmt.Sprintf("%v", err))
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Print(fmt.Sprintf("%v", err))
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func Decrypt(encryptedString string, keyString string) (decryptedString string) {

	//key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	block, err := aes.NewCipher([]byte(keyString))
	if err != nil {
		log.Print(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Print(err)
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Print(err)
	}

	return fmt.Sprintf("%s", plaintext)
}
