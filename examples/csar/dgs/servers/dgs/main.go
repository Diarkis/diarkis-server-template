// © 2019-2024 Diarkis Inc. All rights reserved.

package main

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Diarkis/diarkis"
	pufferDgs "github.com/Diarkis/diarkis-server-template/examples/csar/dgs/puffer/go/dgs"
	dgsapp "github.com/Diarkis/diarkis/dgs/app"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/encryption"
	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/mesh"
	"github.com/Diarkis/diarkis/server"
)

var logger = log.New("DGS")

func main() {
	rootpath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	logConfigPath := "configs/shared/log.json"
	meshConfigPath := "configs/shared/mesh.json"

	log.Setup(filepath.Join(rootpath, logConfigPath))
	mesh.Setup(filepath.Join(rootpath, meshConfigPath))
	server.SetupAsUDPServer(filepath.Join(rootpath, "configs/dgs/main.json")) // define as custom UDP server type: DGS

	dgsapp.SetOnInstanceStart(handleDGSStart)
	dgsapp.SetOnInstanceRPC(handleDGSRPC)
	if err := dgsapp.Setup(); err != nil {
		panic(err)
	}

	mesh.SetNodeType("DGS")
	mesh.SetNodeRole("DGS")

	diarkis.Start()
}

func handleDGSStart(payload []byte) ([]byte, error) {
	proto := pufferDgs.NewMeshInstanceCreateRequest()
	err := proto.Unpack(payload)
	if err != nil {
		logger.Errorf("failed to parse DGS start payload", "error", err)
		return nil, err
	}

	serverHost := mesh.GetMyNodeEndPoint()
	serverHost, serverPort, err := net.SplitHostPort(serverHost)
	if err != nil {
		logger.Errorf("failed to parse current node endpoint", "error", err)
		return nil, err
	}

	serverPortInt, err := strconv.Atoi(serverPort)
	if err != nil {
		logger.Errorf("failed to parse current node endpoint", "error", err)
		return nil, err
	}

	response := pufferDgs.NewMeshInstanceCreateResponse()
	for _, memberID := range proto.MemberIDs {
		sid, err := encryption.CreateKey()
		if err != nil {
			logger.Errorf("failed to generate sid", "error", err)
			return nil, err
		}
		key, err := encryption.CreateKey()
		if err != nil {
			logger.Errorf("failed to generate credentials key", "error", err)
			return nil, err
		}
		iv, err := encryption.CreateKey()
		if err != nil {
			logger.Errorf("failed to generate credentials iv", "error", err)
			return nil, err
		}
		macKey, err := encryption.CreateKey()
		if err != nil {
			logger.Errorf("failed to generate credentials macKey", "error", err)
			return nil, err
		}

		// dummy response
		response.Credentials = append(response.Credentials, &pufferDgs.MeshInstanceCreateResponseUserCredentials{
			SID:        sid,
			IV:         iv,
			Key:        key,
			MacKey:     macKey,
			ServerType: diarkisexec.GetServerType(),
			UID:        memberID,
			ServerPort: uint16(serverPortInt),
			ServerHost: serverHost,
		})
	}

	// Simulate a DGS creation that takes some time.
	// Uncomment to test the behavior.
	// time.Sleep(time.Second * 15)

	return response.Pack(), nil
}

func handleDGSRPC(payload []byte) ([]byte, error) {
	proto := pufferDgs.NewMeshInstanceRPCRequest()
	err := proto.Unpack(payload)
	if err != nil {
		logger.Errorf("failed to parse DGS RPC payload", "error", err)
		return nil, err
	}

	// Only support backfill type
	if proto.Type != "backfill" {
		return nil, errors.New("unsupported type")
	}

	backfillProto := pufferDgs.NewMeshInstanceCreateRequest()
	err = backfillProto.Unpack(proto.Data)
	if err != nil {
		logger.Errorf("failed to parse DGS backfill payload", "error", err)
		return nil, err
	}

	if len(backfillProto.MemberIDs) != 1 {
		logger.Errorf("backfill does not support multiple members")
		return nil, errors.New("does not support multiple members")
	}

	serverHost := mesh.GetMyNodeEndPoint()
	serverHost, serverPort, err := net.SplitHostPort(serverHost)
	if err != nil {
		logger.Errorf("failed to parse current node endpoint", "error", err)
		return nil, err
	}

	serverPortInt, err := strconv.Atoi(serverPort)
	if err != nil {
		logger.Errorf("failed to parse current node endpoint", "error", err)
		return nil, err
	}

	memberID := backfillProto.MemberIDs[0]

	sid, err := encryption.CreateKey()
	if err != nil {
		logger.Errorf("failed to generate sid", "error", err)
		return nil, err
	}
	key, err := encryption.CreateKey()
	if err != nil {
		logger.Errorf("failed to generate credentials key", "error", err)
		return nil, err
	}
	iv, err := encryption.CreateKey()
	if err != nil {
		logger.Errorf("failed to generate credentials iv", "error", err)
		return nil, err
	}
	macKey, err := encryption.CreateKey()
	if err != nil {
		logger.Errorf("failed to generate credentials macKey", "error", err)
		return nil, err
	}

	response := &pufferDgs.MeshInstanceCreateResponseUserCredentials{
		SID:        sid,
		IV:         iv,
		Key:        key,
		MacKey:     macKey,
		ServerType: diarkisexec.GetServerType(),
		UID:        memberID,
		ServerPort: uint16(serverPortInt),
		ServerHost: serverHost,
	}

	logger.Info("DGS backfill response", "UserID", memberID, "Data", response.String())

	return response.Pack(), nil
}
