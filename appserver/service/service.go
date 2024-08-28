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
			}
			log.Printf("Generated key with ID: %v, player: %d", keyIDs[i], i)
			return err
		}()
	}
	return sessionID, nil
}
