package service

import (
	"context"
	"projects/Repairment_service/Repairment_user_service/config"
	"projects/Repairment_service/Repairment_user_service/genproto/user_service"
	"projects/Repairment_service/Repairment_user_service/grpc/client"
	"projects/Repairment_service/Repairment_user_service/models"
	"projects/Repairment_service/Repairment_user_service/pkg/logger"
	"projects/Repairment_service/Repairment_user_service/storage"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type UserService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	*user_service.UnimplementedUserServiceServer
}

func NewUserService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, services client.ServiceManagerI) *UserService {
	return &UserService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: services,
	}
}

func (u *UserService) Create(ctx context.Context, req *user_service.CreateUserRequest) (resp *user_service.User, err error) {
	u.log.Info("-------Create user---------->", logger.Any("req", req))

	pkey, err := u.strg.User().Create(ctx, req)
	if err != nil {
		u.log.Error("!!!CreateUser -> User -> Create", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err = u.strg.User().GetById(ctx, pkey)
	if err != nil {
		u.log.Error("!!!CreateUser -> User -> GetById", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (u *UserService) GetById(ctx context.Context, req *user_service.UserPrimaryKey) (resp *user_service.User, err error) {
	u.log.Info("-------GetById user---------->", logger.Any("req", req))

	resp, err = u.strg.User().GetById(ctx, req)
	if err != nil {
		u.log.Error("!!!GetById -> User -> GetById", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (u *UserService) GetList(ctx context.Context, req *user_service.GetListUserRequest) (resp *user_service.GetListUserResponse, err error) {
	u.log.Info("-------GetAll user---------->", logger.Any("req", req))

	resp, err = u.strg.User().GetList(ctx, req)
	if err != nil {
		u.log.Error("!!!GetAll -> User -> GetAll", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return
}

func (u *UserService) Update(ctx context.Context, req *user_service.UpdateUserRequest) (resp *user_service.User, err error) {
	u.log.Info("-------Update user---------->", logger.Any("req", req))

	rowsAffected, err := u.strg.User().Update(ctx, req)
	if err != nil {
		u.log.Error("!!!Update -> User -> Update", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	resp, err = u.strg.User().GetById(ctx, &user_service.UserPrimaryKey{Id: req.Id})
	if err != nil {
		u.log.Error("!!!Update -> User -> GetById", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, nil

	return
}

func (u *UserService) UpdatePatch(ctx context.Context, req *user_service.UpdatePatchUser) (resp *user_service.User, err error) {
	u.log.Info("-------UpdatePatch user---------->", logger.Any("req", req))

	updatePatchModel := models.UpdatePatchRequest {
		Id: req.GetId(),
		Fields: req.GetFields().AsMap(),
	}

	rowsAffected, err := u.strg.User().UpdatePatch(ctx, &updatePatchModel)
	if err != nil {
		u.log.Error("!!!UpdatePatch -> User -> UpdatePatch", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	resp, err = u.strg.User().GetById(ctx, &user_service.UserPrimaryKey{Id: req.Id})
	if err != nil {
		u.log.Error("!!!UpdatePatch -> User -> GetById", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, nil
}

func (b *UserService) Delete(ctx context.Context, req *user_service.UserPrimaryKey) (*empty.Empty, error) {
	b.log.Info("---DeleteUser--->", logger.Any("req", req))

	err := b.strg.User().Delete(ctx, req)
	if err != nil {
		b.log.Error("!!!DeleteOrder--->", logger.Error(err))
		return &empty.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &empty.Empty{}, nil
}