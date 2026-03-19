package repository

import (
	"context"
	"user-management-api/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindAll()
	Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	FindByUUID(uuid string)
	Update(ctx context.Context, userParams sqlc.UpdateUserParams) (sqlc.User, error)
	SoftDelete(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
	Delete(ctx context.Context, userUuid uuid.UUID) error
	Restore(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
	FindByEmail(email string)
}
