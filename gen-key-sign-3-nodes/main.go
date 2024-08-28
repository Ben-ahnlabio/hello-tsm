package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"

	"gitlab.com/Blockdaemon/go-tsm-sdkv2/v64/tsm"
	"golang.org/x/sync/errgroup"
)

func createKeyAndSign() {
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

	threshold := 1
	keyGenPlayers := []int{0, 1, 2}
	keyGenSessionConfig := tsm.NewSessionConfig(tsm.GenerateSessionID(), keyGenPlayers, nil)

	fmt.Println("Generating key using players", keyGenPlayers)
	ctx := context.Background()
	keyIDs := make([]string, len(clients))

	var eg errgroup.Group
	for i, client := range clients {
		client, i := client, i
		eg.Go(func() error {
			var err error
			keyIDs[i], err = client.ECDSA().GenerateKey(ctx, keyGenSessionConfig, threshold, "secp256k1", "")
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		panic(err)
	}

	for i := 1; i < len(keyIDs); i++ {
		if keyIDs[0] != keyIDs[i] {
			panic("key IDs do not match")
		}
	}

	keyID := keyIDs[0]
	fmt.Println("Generated key with ID:", keyID)

	var derivationPath []uint32 = nil

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

	message := []byte("This is a message to be signed")
	msgHash := sha256.Sum256(message)

	signPlayers := []int{0, 1}
	sessionID := tsm.GenerateSessionID()
	signSessionConfig := tsm.NewSessionConfig(sessionID, signPlayers, nil)

	fmt.Println("Creating signature using players", signPlayers)
	partialSignaturesLock := sync.Mutex{}
	var partialSignatures [][]byte
	for _, player := range signPlayers {
		player := player
		eg.Go(func() error {
			if partialSignResult, err := clients[player].ECDSA().Sign(ctx, signSessionConfig, keyID, derivationPath, msgHash[:]); err != nil {
				return err
			} else {
				partialSignaturesLock.Lock()
				partialSignatures = append(partialSignatures, partialSignResult.PartialSignature)
				partialSignaturesLock.Unlock()
				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		panic(err)
	}

	signature, err := tsm.ECDSAFinalizeSignature(msgHash[:], partialSignatures)
	if err != nil {
		panic(err)
	}

	if err = tsm.ECDSAVerifySignature(publicKey, msgHash[:], signature.ASN1()); err != nil {
		panic(err)
	}

	fmt.Println("Signature:", hex.EncodeToString(signature.ASN1()))
}

func main() {
	fmt.Println("hello tsm")

	createKeyAndSign()
	keyListing()
}

func keyListing() {
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

	ctx := context.Background()
	for idx, client := range clients {
		keyIDs, err := client.KeyManagement().ListKeys(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("node: %d, keyIDs %v", idx, keyIDs)
	}
}
