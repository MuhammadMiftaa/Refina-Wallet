package server

import (
	"net"

	"refina-wallet/config/db"
	"refina-wallet/config/env"
	"refina-wallet/interface/grpc/client"
	"refina-wallet/interface/queue"
	"refina-wallet/internal/repository"
	"refina-wallet/internal/service"

	wpb "github.com/MuhammadMiftaa/Refina-Protobuf/wallet"
	"google.golang.org/grpc"
)

func SetupGRPCServer(dbInstance db.DatabaseClient, queueInstance queue.RabbitMQClient) (*grpc.Server, *net.Listener, error) {
	lis, err := net.Listen("tcp", ":"+env.Cfg.Server.GRPCPort)
	if err != nil {
		return nil, nil, err
	}

	s := grpc.NewServer()

	txManager := repository.NewTxManager(dbInstance.GetDB())
	walletsRepo := repository.NewWalletRepository(dbInstance.GetDB())
	walletTypesRepo := repository.NewWalletTypesRepository(dbInstance.GetDB())
	outboxRepo := repository.NewOutboxRepository(dbInstance.GetDB())
	transactionClient := client.NewTransactionClient(client.GetManager().GetTransactionClient())

	walletService := service.NewWalletService(
		txManager,
		walletsRepo,
		walletTypesRepo,
		outboxRepo,
		transactionClient,
		queueInstance,
	)
	walletTypesService := service.NewWalletTypesService(txManager, walletTypesRepo)

	walletServer := &walletServer{
		walletService:      walletService,
		walletTypesService: walletTypesService,
		walletsRepository:  walletsRepo,
	}
	wpb.RegisterWalletServiceServer(s, walletServer)

	return s, &lis, nil
}
