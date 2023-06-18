package client

import "projects/Repairment_service/Repairment_user_service/config"

type ServiceManagerI interface{}

type grpcClients struct{}

func NewGrpcClients(cfg config.Config) (ServiceManagerI,error) {

	return &grpcClients{}, nil
}