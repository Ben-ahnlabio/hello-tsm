package service

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/ahnlabio/tsm-appserver/tsmutils"
	"gitlab.com/Blockdaemon/go-tsm-sdkv2/v64/tsm"
)

func GenerateKey(publicKeyStr string) (string, error) {
	player0PublicTenantKey, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		log.Printf("GenerateKey Service Error decoding public key: %v", err)
		return "", err
	}

	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}

	threshold := 1            // The security threshold of the key
	players := []int{0, 1, 2} // The players (nodes) that should generate a sharing of the key
	curveName := "secp256k1"

	sessionID := tsm.GenerateSessionID()
	sessionConfig := tsm.NewSessionConfig(sessionID, players, dynamicPublicKeys)
	ctx := context.Background()

	clients := tsmutils.GetClients()

	keyIDs := make([]string, len(clients))
	for i, client := range clients {
		client, i := client, i
		go func() error {
			var err error
			keyIDs[i], err = client.ECDSA().GenerateKey(ctx, sessionConfig, threshold, curveName, "")
			if err != nil {
				log.Printf("Error generating key: %v", err)
				return err
			}
			log.Printf("Generated key with ID: %v, player: %d", keyIDs[i], i)
			return err
		}()
	}
	return sessionID, nil
}

func CopyKey(publicKeyStr string, existingKeyID string) (string, error) {
	player0PublicTenantKey, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		log.Printf("GenerateKey Service Error decoding public key: %v", err)
		return "", err
	}

	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}

	newThreshold := 1            // The security threshold of the key
	newPlayers := []int{0, 1, 2} // The players (nodes) that should generate a sharing of the key
	curveName := "secp256k1"

	sessionID := tsm.GenerateSessionID()
	keyCopySessionConfig := tsm.NewSessionConfig(sessionID, newPlayers, dynamicPublicKeys)
	ctx := context.Background()

	clients := tsmutils.GetClients()
	for i, client := range clients {
		client, i := client, i
		go func() error {
			var err error
			newKeyId, err := client.ECDSA().CopyKey(ctx, keyCopySessionConfig, existingKeyID, curveName, newThreshold, "")
			if err != nil {
				log.Printf("Error generating key: %v", err)
				return err
			}
			log.Printf("Copied existingKeyID: %s, newKeyId: %s, player: %d", existingKeyID, newKeyId, i)
			return err
		}()
	}
	return sessionID, nil
}

func PreSign(publicKeyStr string, keyId string) (string, error) {
	player0PublicTenantKey, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		log.Printf("GenerateKey Service Error decoding public key: %v", err)
		return "", err
	}

	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}

	presignSessionID := tsm.GenerateSessionID()
	players := []int{0, 1, 2}
	presigSessionConfig := tsm.NewSessionConfig(presignSessionID, players, dynamicPublicKeys)

	clients := tsmutils.GetClients()

	presignatureIDs := make([][]string, len(clients))
	for i, client := range clients {
		i, client := i, client

		go func() error {
			var err error
			presignatureIDs[i], err = client.ECDSA().GeneratePresignatures(context.TODO(), presigSessionConfig, keyId, 1)
			if err != nil {
				log.Printf("Error generating presignature: %v", err)
			}

			log.Printf("Generated presignature with ID: %v, player: %d", presignatureIDs[i], i)
			return err
		}()
	}

	return presignSessionID, nil
}

func FinalizeSign(preSignatureId string, messageHash string, keyId string) ([]string, error) {
	clients := tsmutils.GetClients()
	partialSignatures := make([]string, 0)

	messageHashBytes, err := base64.StdEncoding.DecodeString(messageHash)
	if err != nil {
		return partialSignatures, err
	}

	for _, client := range clients {
		client := client
		partialSignResult, err := client.ECDSA().SignWithPresignature(context.TODO(), keyId, preSignatureId, nil, messageHashBytes[:])
		if err != nil {
			return partialSignatures, err
		}
		partialSignatureStr := base64.StdEncoding.EncodeToString(partialSignResult.PartialSignature)
		partialSignatures = append(partialSignatures, partialSignatureStr)
	}

	return partialSignatures, nil
}
