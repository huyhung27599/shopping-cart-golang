package repository

import (
	"context"
	"user-management-api/internal/db/sqlc"

	"github.com/google/uuid"
)

type SqlUserRepository struct {
	db sqlc.Querier
}

func NewSqlUserRepository(db sqlc.Querier) UserRepository {
	return &SqlUserRepository{db: db}
}

func (ur *SqlUserRepository) FindAll() {
	
}

func (ur *SqlUserRepository) Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.db.CreateUser(ctx, userParams)
	if err != nil {
		return sqlc.User{}, err
	}
	return user, nil
}

func (ur *SqlUserRepository) FindByUUID(uuid string) {
	

	
}

func (ur *SqlUserRepository) Update(ctx context.Context, userParams sqlc.UpdateUserParams) (sqlc.User, error) {
	user, err := ur.db.UpdateUser(ctx, userParams)
	if err != nil {
		return sqlc.User{}, err
	}
	return user, nil
	
}

func (ur *SqlUserRepository) Delete(ctx context.Context, userUuid uuid.UUID)  error {
	_, err := ur.db.TrashUser(ctx, userUuid)
	if err != nil {
		return err
	}
	return nil
}

func (ur *SqlUserRepository) SoftDelete(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.SoftDeleteUser(ctx, userUuid)
	if err != nil {
		return sqlc.User{}, err
	}
	return user, nil
}

func (ur *SqlUserRepository) Restore(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.RestoreUser(ctx, userUuid)
	if err != nil {
		return sqlc.User{}, err
	}
	return user, nil
}

func (ur *SqlUserRepository) FindByEmail(email string)  {
	
}
