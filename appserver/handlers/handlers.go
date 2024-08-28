package handlers

import (
	"log"
	"net/http"

	"github.com/ahnlabio/tsm-appserver/service"
	"github.com/gin-gonic/gin"
)

type GenerateKeyRequestBody struct {
	PublicKey string `json:"publicKey"`
}

type GenerateKeyResponseBody struct {
	SessionId string `json:"sessionId"`
}

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
	c.JSON(http.StatusOK, GenerateKeyResponseBody{SessionId: sessionId})
}
