package storage

import (
	"context"
	"projects/Repairment_service/Repairment_user_service/genproto/user_service"
	"projects/Repairment_service/Repairment_user_service/models"
)

type StorageI interface {
	CloseDB()
	User() UserRepoI
}

type UserRepoI interface {
	Create(context.Context, *user_service.CreateUserRequest) (*user_service.UserPrimaryKey, error)
	GetById(context.Context, *user_service.UserPrimaryKey) (*user_service.User, error)
	GetList(context.Context, *user_service.GetListUserRequest) (*user_service.GetListUserResponse, error)
	Update(context.Context, *user_service.UpdateUserRequest) (int64, error)
	UpdatePatch(context.Context, *models.UpdatePatchRequest) (int64, error)
	Delete(context.Context, *user_service.UserPrimaryKey) error
}