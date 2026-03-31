package repository

import (
	"context"
	"user-management-api/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll(ctx context.Context, search string, orderBy string, sort string, limit int32, offset int32) ([]sqlc.User, error)
	GetAllV2(ctx context.Context, search string, orderBy string, sort string, limit int32, offset int32, deleted bool) ([]sqlc.User, error)
	CountUsers(ctx context.Context, search string, deleted bool) (int64, error)
	Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	Update(ctx context.Context, userParams sqlc.UpdateUserParams) (sqlc.User, error)
	SoftDelete(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
	Delete(ctx context.Context, userUuid uuid.UUID) error
	Restore(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
	FindByEmail(ctx context.Context, email string) (sqlc.User, error)
	UpdatePassword(ctx context.Context, userUUID uuid.UUID, newPassword string) (sqlc.User, error)
	GetByUUID(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
}
