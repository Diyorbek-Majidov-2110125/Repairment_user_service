package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"projects/Repairment_service/Repairment_user_service/genproto/user_service"
	"projects/Repairment_service/Repairment_user_service/models"
	"projects/Repairment_service/Repairment_user_service/pkg/helper"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)



type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(ctx context.Context, req *user_service.CreateUserRequest) (resp *user_service.UserPrimaryKey, err error) {
	id := uuid.New().String()

	query := `
		INSERT INTO "user" (
			id,
			full_name,
			phone_number, 
			updated_at
		)
		VALUES ($1, $2, $3, NOW())
	`

	_, err = u.db.Exec(ctx, query, id, req.FullName, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	return &user_service.UserPrimaryKey{Id: id}, nil
}

func (u *userRepo) GetById(ctx context.Context, req *user_service.UserPrimaryKey) (resp *user_service.User, err error) {
	query := `
		Select 
			id,
			full_name,
			phone_number,
			created_at,
			updated_at
		from "user"
		where id = $1
	`

	var (
		id         sql.NullString
		full_name sql.NullString
		phone_number  sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)

	err = u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&full_name,
		&phone_number,
		&created_at,
		&updated_at,
	)

	if err != nil {
		return resp, err
	}

	resp = &user_service.User{
		Id:        id.String,
		FullName: full_name.String,
		PhoneNumber:  phone_number.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}
	return
}

func (u *userRepo) GetList(ctx context.Context, req *user_service.GetListUserRequest) (resp *user_service.GetListUserResponse, err error) {
	resp = &user_service.GetListUserResponse{}

	var (
		query  string
		limit  = ""
		offset = " OFFSET 0"
		params = make(map[string]interface{})
		filter = " WHERE TRUE"
		sort   = " ORDER BY created_at DESC"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			full_name,
			phone_number,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM "users"
		`

	if len(req.GetSearch()) > 0 {
		filter += " AND (full_name || ' ' phone_number) ILIKE '%' || '" + req.Search + " || '%' "
	}
	if req.GetLimit() > 0 {
		limit = " LIMIT :limit"
		params["limit"] = req.Limit
	}
	if req.GetOffset() > 0 {
		offset = " OFFSET :offset"
		params["offset"] = req.Offset
	}

	query += filter + sort + offset + limit

	query, args := helper.ReplaceQueryParams(query, params)
	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return resp, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id         sql.NullString
			full_name sql.NullString
			phone_number  sql.NullString
			created_at sql.NullString
			updated_at sql.NullString
		)

		err := rows.Scan(
			&resp.Total,
			&id,
			&full_name,
			&phone_number,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return resp, err
		}

		resp.Users = append(resp.Users, &user_service.User{
			Id:        id.String,
			FullName: full_name.String,
			PhoneNumber:  phone_number.String,
			CreatedAt: created_at.String,
			UpdatedAt: updated_at.String,
		})
	}

	return
}

func (u *userRepo) Update(ctx context.Context, req *user_service.UpdateUserRequest) (rowsAffected int64, err error) {
	var (
		query  string
		params = make(map[string]interface{})
	)

	query = `
		UPDATE "user"
		SET
			full_name = :full_name,
			phone_number = :phone_number,
			updated_at = NOW()
		WHERE id = :id
	`
	params = map[string]interface{}{
		"full_name": req.GetFullName(),
		"id":         req.GetId(),
		"phone_number":  req.GetPhoneNumber(),
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (u *userRepo) UpdatePatch(ctx context.Context, req *models.UpdatePatchRequest) (rowsAffected int64, err error) {
	var (
		set   = " SET "
		ind   = 0
		query string
	)

	if len(req.Fields) == 0 {
		err = errors.New("no updates provided")
		return
	}

	req.Fields["id"] = req.Id

	for key := range req.Fields {
		set += fmt.Sprintf(" %s = :%s ", key, key)
		if ind != len(req.Fields)-1 {
			set += ", "
		}
		ind++
	}

	query = `
		UPDATE "user"
		    ` + set + ` , updated_at = NOW()
		WHERE id = :id
	`

	query, args := helper.ReplaceQueryParams(query, req.Fields)

	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (u *userRepo) Delete(ctx context.Context, req *user_service.UserPrimaryKey) error {
	query := `
		DELETE FROM "user"
		WHERE id = $1
	`

	_, err := u.db.Exec(ctx, query, req.Id)
	if err != nil {
		return err
	}

	return nil
}
