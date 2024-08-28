package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"example.com/tsmutils"
	"gitlab.com/Blockdaemon/go-tsm-sdkv2/v64/tsm"
)

var tsmDynamicMob0 = tsm.Configuration{URL: "http://localhost:8510"}.WithAPIKeyAuthentication("apikey0")
var tsmDynamicMob1 = tsm.Configuration{URL: "http://localhost:8511"}.WithAPIKeyAuthentication("apikey0")

type GetKeyResult struct {
	KeyId         string `json:"keyId"`
	UserPublicKey string `json:"userPublicKey"`
}

type CopyKeyResult struct {
	NewKeyId      string `json:"keyId"`
	UserPublicKey string `json:"userPublicKey"`
}

func main() {
	fmt.Printf("Hello, tsm client.\n")

	mobile0PublicKey := "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="
	mobile1PublicKey := "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEkAzm+8yn+d0ypywEwtgNnjisUkXBH17HpOd9YqRDybobqmCuaZA8cqAyLFS/qlu6j7lKCDWBwTElXJgvG9nywQ=="

	genKeyResult := client0GenKey(mobile0PublicKey)
	log.Printf("genKeyResult: %v\n", genKeyResult)
	copyKeyResult := client1CopyKey(mobile1PublicKey, genKeyResult.KeyId)
	log.Printf("copyKeyResult: %v\n", copyKeyResult)

	if genKeyResult.UserPublicKey != copyKeyResult.UserPublicKey {
		panic("User public key mismatch")
	}
}

func client0GenKey(nodePubKey string) *GetKeyResult {
	sessionId := generateKey(nodePubKey)
	player0PublicTenantKey, err := base64.StdEncoding.DecodeString(nodePubKey)
	if err != nil {
		panic(err)
	}

	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}
	players := []int{0, 1, 2} // The players (nodes) that should generate a sharing of the key
	sessionConfig := tsm.NewSessionConfig(sessionId, players, dynamicPublicKeys)
	ctx := context.Background()

	client := tsmutils.GetClientFromConfig(tsmDynamicMob0)
	threshold := 1
	keyId, err := client.ECDSA().GenerateKey(ctx, sessionConfig, threshold, "secp256k1", "")
	if err != nil {
		panic(err)
	}
	log.Printf("keyId: %s\n", keyId)
	userPubkey := tsmutils.GetPubkeyStringFromClient(client, keyId)
	log.Printf("userPubkey: %s\n", userPubkey)
	return &GetKeyResult{KeyId: keyId, UserPublicKey: userPubkey}
}

func client1CopyKey(nodePubKey string, keyId string) *CopyKeyResult {
	sessionId := copyKey(nodePubKey, keyId)
	player0PublicTenantKey, err := base64.StdEncoding.DecodeString(nodePubKey)
	if err != nil {
		panic(err)
	}

	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}

	client := tsmutils.GetClientFromConfig(tsmDynamicMob1)
	newThreshold := 1
	newPlayers := []int{0, 1, 2} // The players (nodes) that should generate a sharing of the key
	keyCopySessionConfig := tsm.NewSessionConfig(sessionId, newPlayers, dynamicPublicKeys)

	ctx := context.Background()
	curveName := "secp256k1"
	newKeyId, err := client.ECDSA().CopyKey(ctx, keyCopySessionConfig, "", curveName, newThreshold, "")
	if err != nil {
		panic(err)
	}
	userPubkey := tsmutils.GetPubkeyStringFromClient(client, newKeyId)
	return &CopyKeyResult{
		NewKeyId:      newKeyId,
		UserPublicKey: userPubkey,
	}
}

type GenerateKeyRequestBody struct {
	PublicKey string `json:"publicKey"`
}

type GenerateKeyResponse struct {
	SessionId string `json:"sessionId"`
}

func generateKey(publicKey string) string {

	url := "http://localhost:3000/generateKey"
	addrReqBody := GenerateKeyRequestBody{
		PublicKey: publicKey,
	}
	value, _ := json.Marshal(addrReqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(value))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "ABC")

	client := &http.Client{Timeout: time.Duration(3000) * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("failed to get session id. status code: %d", resp.StatusCode))
	}

	var resObj GenerateKeyResponse
	err = json.Unmarshal(body, &resObj)
	if err != nil {
		panic(err)
	}

	return resObj.SessionId
}

type CopyKeyRequestBody struct {
	PublicKey string `json:"publicKey"`
	KeyId     string `json:"keyId"`
}

type CopyKeyResponse struct {
	SessionId string `json:"sessionId"`
}

func copyKey(publicKey string, existingKeyId string) string {
	url := "http://localhost:3000/copyKey"
	addrReqBody := CopyKeyRequestBody{
		PublicKey: publicKey,
		KeyId:     existingKeyId,
	}
	value, _ := json.Marshal(addrReqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(value))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "ABC")

	client := &http.Client{Timeout: time.Duration(3000) * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("failed to get session id. status code: %d", resp.StatusCode))
	}

	var resObj CopyKeyResponse
	err = json.Unmarshal(body, &resObj)
	if err != nil {
		panic(err)
	}

	return resObj.SessionId
}
