package handlers

import (
	"log"
	"net/http"

	"github.com/ahnlabio/tsm-appserver/service"
	"github.com/gin-gonic/gin"
)

type GenerateKeyRequestBody struct {
	PublicKey string `json:"publicKey" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
}

type GenerateKeyResponseBody struct {
	SessionId string `json:"sessionId" binding:"required" exaple:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
}

// GenerateKeyHandler godoc
// @Summary Generate a session key
// @Description Generate a session key
// @Tags session
// @Accept json
// @Produce json
// @Param body body GenerateKeyRequestBody true "Public key"
// @Success 200 {object} GenerateKeyResponseBody
// @Router /generateKey [post]
func GenerateKeyHandler(c *gin.Context) {
	var requestBody GenerateKeyRequestBody
	err := c.ShouldBind(&requestBody)
	if err != nil {
		log.Printf("[GenerateKeyHandler] c.ShouldBind Error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionId, err := service.GenerateKey(requestBody.PublicKey)
	if err != nil {
		log.Printf("[GenerateKeyHandler] service.GenerateKey Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[GenerateKeyHandler] session id: %s", sessionId)
	c.JSON(http.StatusOK, GenerateKeyResponseBody{SessionId: sessionId})
}

type CopyKeyRequestBody struct {
	PublicKey string `json:"publicKey" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
	KeyId     string `json:"keyId" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
}

type CopyResponseBody struct {
	SessionId string `json:"sessionId" binding:"required" exaple:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
}

// CopyKeyHandler godoc
// @Summary Copy a session key
// @Description Copy a session key
// @Tags session
// @Accept json
// @Produce json
// @Param body body CopyKeyRequestBody true "Public key and key ID"
// @Success 200 {object} CopyResponseBody
// @Router /copyKey [post]
func CopyKeyHandler(c *gin.Context) {
	var requestBody CopyKeyRequestBody
	err := c.ShouldBind(&requestBody)
	if err != nil {
		log.Printf("[GenerateKeyHandler] c.ShouldBind Error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionId, err := service.CopyKey(requestBody.PublicKey, requestBody.KeyId)
	if err != nil {
		log.Printf("[GenerateKeyHandler] service.GenerateKey Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[CopyKeyHandler] session id: %s", sessionId)
	c.JSON(http.StatusOK, GenerateKeyResponseBody{SessionId: sessionId})
}
