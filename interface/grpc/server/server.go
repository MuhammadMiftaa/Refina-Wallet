package server

import (
	"net"

	"refina-wallet/config/db"
	"refina-wallet/config/env"
	"refina-wallet/internal/repository"

	wpb "github.com/MuhammadMiftaa/Refina-Protobuf/wallet"
	"google.golang.org/grpc"
)

func SetupGRPCServer(dbInstance db.DatabaseClient) (*grpc.Server, *net.Listener, error) {
	lis, err := net.Listen("tcp", ":"+env.Cfg.Server.GRPCPort)
	if err != nil {
		return nil, nil, err
	}

	s := grpc.NewServer()

	walletServer := &walletServer{
		walletsRepository: repository.NewWalletRepository(dbInstance.GetDB()),
	}
	wpb.RegisterWalletServiceServer(s, walletServer)

	return s, &lis, nil
}
