package tsmutils

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"gitlab.com/Blockdaemon/go-tsm-sdkv2/v64/tsm"
)

func GetPubkeyStringFromClient(client *tsm.Client, keyId string) string {
	ctx := context.Background()
	var derivationPath []uint32 = nil
	publicKey, err := client.ECDSA().PublicKey(ctx, keyId, derivationPath)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(publicKey)
}

func GetClientsFromConfigs(configs []*tsm.Configuration) []*tsm.Client {
	clients := make([]*tsm.Client, len(configs))
	for i, config := range configs {
		var err error
		if clients[i], err = tsm.NewClient(config); err != nil {
			panic(err)
		}
	}
	return clients
}

func GenerateSessionConfig(players []int, pubkeyStr string) *tsm.SessionConfig {
	player0PublicTenantKey, err := base64.StdEncoding.DecodeString(pubkeyStr)
	if err != nil {
		panic(err)
	}

	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}

	sessionID := tsm.GenerateSessionID()
	sessionConfig := tsm.NewSessionConfig(sessionID, players, dynamicPublicKeys)
	return sessionConfig
}

// func GenerateKey(threshold int, curveName string) *tsm.SessionConfig {
// 	player0PublicTenantKey, err := base64.StdEncoding.DecodeString("MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A==")
// 	if err != nil {
// 		panic(err)
// 	}

// 	dynamicPublicKeys := map[int][]byte{
// 		0: player0PublicTenantKey,
// 	}

// 	sessionID := tsm.GenerateSessionID()
// 	sessionConfig := tsm.NewSessionConfig(sessionID, players, dynamicPublicKeys)
// 	ctx := context.Background()

// 	keyIDs := make([]string, len(clients))
// 	var eg errgroup.Group
// 	for i, client := range clients {
// 		client, i := client, i
// 		eg.Go(func() error {
// 			var err error
// 			keyIDs[i], err = client.ECDSA().GenerateKey(ctx, sessionConfig, threshold, curveName, "")
// 			return err
// 		})
// 	}
// 	if err := eg.Wait(); err != nil {
// 		panic(err)
// 	}
// }

func GetClientFromConfig(config *tsm.Configuration) *tsm.Client {
	client, err := tsm.NewClient(config)
	if err != nil {
		panic(err)
	}
	return client
}

func KeyListing(configs []*tsm.Configuration) {
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
		fmt.Printf("node: %d, keyIDs %v\n", idx, keyIDs)
	}
}
