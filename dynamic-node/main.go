package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"gitlab.com/Blockdaemon/go-tsm-sdkv2/v64/tsm"
	"golang.org/x/sync/errgroup"

	"example.com/tsmutils"
)

func main() {
	//keyListing()
	keyID := mob0_key()
	mob1_key_copy(keyID)
}

var tsmDynamic0 = tsm.Configuration{URL: "http://localhost:8510"}.WithAPIKeyAuthentication("apikey0")
var tsmStatic1 = tsm.Configuration{URL: "http://localhost:8501"}.WithAPIKeyAuthentication("apikey1")
var tsmStatic2 = tsm.Configuration{URL: "http://localhost:8502"}.WithAPIKeyAuthentication("apikey2")
var tsmDynamicMob0 = tsm.Configuration{URL: "http://localhost:8511"}.WithAPIKeyAuthentication("apikey0")
var tsmDynamicMob1 = tsm.Configuration{URL: "http://localhost:8512"}.WithAPIKeyAuthentication("apikey0")

func mob0_key() string {
	configs := []*tsm.Configuration{
		tsmDynamic0,
		tsmStatic1,
		tsmStatic2,
	}

	clients := make([]*tsm.Client, len(configs))
	for i, config := range configs {
		var err error
		if clients[i], err = tsm.NewClient(config); err != nil {
			panic(err)
		}
	}

	// Generate a key, with MPC Node 0 dynamically configured

	threshold := 1            // The security threshold of the key
	players := []int{0, 1, 2} // The players (nodes) that should generate a sharing of the key
	curveName := "secp256k1"

	// Provide Node 0 public key dynamically
	player0PublicTenantKey, err := base64.StdEncoding.DecodeString("MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A==")
	if err != nil {
		panic(err)
	}
	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}

	sessionID := tsm.GenerateSessionID()
	sessionConfig := tsm.NewSessionConfig(sessionID, players, dynamicPublicKeys)
	ctx := context.Background()

	keyIDs := make([]string, len(clients))
	var eg errgroup.Group
	for i, client := range clients {
		client, i := client, i
		eg.Go(func() error {
			var err error
			keyIDs[i], err = client.ECDSA().GenerateKey(ctx, sessionConfig, threshold, curveName, "")
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		panic(err)
	}

	fmt.Println("Generated key with ID:", keyIDs[0])

	var derivationPath []uint32 = nil

	keyID := keyIDs[0]
	publicKeys := make([][]byte, len(clients))
	for i, client := range clients {
		var err error
		publicKeys[i], err = client.ECDSA().PublicKey(ctx, keyID, derivationPath)
		if err != nil {
			panic(err)
		}
	}

	publicKey := publicKeys[0]
	fmt.Println("Public key:", hex.EncodeToString(publicKey))

	return keyID
}

func mob1_key_copy(keyID string) {
	configs := []*tsm.Configuration{
		tsmDynamicMob0, tsmStatic1, tsmStatic2,
	}

	clients := make([]*tsm.Client, len(configs))
	for i, config := range configs {
		var err error
		if clients[i], err = tsm.NewClient(config); err != nil {
			panic(err)
		}
	}

	player0PublicTenantKey, err := base64.StdEncoding.DecodeString("MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEkAzm+8yn+d0ypywEwtgNnjisUkXBH17HpOd9YqRDybobqmCuaZA8cqAyLFS/qlu6j7lKCDWBwTElXJgvG9nywQ==")
	if err != nil {
		panic(err)
	}
	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}

	newPlayers := []int{0, 1, 2} // The new set of players
	newThreshold := 1            // The security threshold for the new copy of the key
	keyCopySessionConfig := tsm.NewSessionConfig(tsm.GenerateSessionID(), newPlayers, dynamicPublicKeys)

	fmt.Println("Copying key using players", newPlayers, "and threshold", newThreshold)
	newKeyIDs := make([]string, len(clients))
	var eg errgroup.Group
	ctx := context.Background()
	for i, client := range clients {
		client, i := client, i
		eg.Go(func() error {
			var err error
			var existingKeyID, curveName string
			if i == 0 {
				existingKeyID = ""
				curveName = "secp256k1"
			} else {
				existingKeyID = keyID
				curveName = ""
			}
			newKeyIDs[i], err = client.ECDSA().CopyKey(ctx, keyCopySessionConfig, existingKeyID, curveName, newThreshold, "")
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		panic(err)
	}
	newKeyID := newKeyIDs[0]
	fmt.Println("CopyKey completed; new key ID:", newKeyID)
	publicKey := tsmutils.GetPubkeyStringFromClient(clients[0], newKeyID)
	fmt.Println("Public key:", publicKey)
}

func generateKey() {
	// Create clients for each of the nodes

	configs := []*tsm.Configuration{
		tsm.Configuration{URL: "http://localhost:8500"}.WithAPIKeyAuthentication("apikey0"),
		tsm.Configuration{URL: "http://localhost:8501"}.WithAPIKeyAuthentication("apikey1"),
		tsm.Configuration{URL: "http://localhost:8502"}.WithAPIKeyAuthentication("apikey2"),
	}

	clients := make([]*tsm.Client, len(configs))
	for i, config := range configs {
		var err error
		if clients[i], err = tsm.NewClient(config); err != nil {
			panic(err)
		}
	}

	// Generate a key, with MPC Node 0 dynamically configured

	threshold := 1            // The security threshold of the key
	players := []int{0, 1, 2} // The players (nodes) that should generate a sharing of the key
	curveName := "secp256k1"

	// Provide Node 0 public key dynamically
	player0PublicTenantKey, err := base64.StdEncoding.DecodeString("MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE0OqvUD8ezIIHktmgrDIRh7bwQ3k9G8HZochWXovvQjCm4wQiJBHunl82I9pbeVLD9fa/40Fv8/NRcYiGh/cyUw==")
	if err != nil {
		panic(err)
	}
	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}

	sessionID := tsm.GenerateSessionID()
	sessionConfig := tsm.NewSessionConfig(sessionID, players, dynamicPublicKeys)
	ctx := context.Background()

	keyIDs := make([]string, len(clients))
	var eg errgroup.Group
	for i, client := range clients {
		client, i := client, i
		eg.Go(func() error {
			var err error
			keyIDs[i], err = client.ECDSA().GenerateKey(ctx, sessionConfig, threshold, curveName, "")
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		panic(err)
	}

	fmt.Println("Generated key with ID:", keyIDs[0])
	var derivationPath []uint32 = nil
	publicKey, err := clients[0].ECDSA().PublicKey(ctx, keyIDs[0], derivationPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Public key:", hex.EncodeToString(publicKey))
}
