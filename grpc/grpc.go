package grpc

import (
	"projects/Repairment_service/Repairment_user_service/config"
	"projects/Repairment_service/Repairment_user_service/genproto/user_service"
	"projects/Repairment_service/Repairment_user_service/grpc/client"
	"projects/Repairment_service/Repairment_user_service/grpc/service"
	"projects/Repairment_service/Repairment_user_service/pkg/logger"
	"projects/Repairment_service/Repairment_user_service/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI, svcs client.ServiceManagerI) (grpcServer *grpc.Server) {

	grpcServer = grpc.NewServer()

	user_service.RegisterUserServiceServer(grpcServer, service.NewUserService(cfg, log, strg, svcs))

	reflection.Register(grpcServer)
	return
}