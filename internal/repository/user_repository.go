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

func (ur *SqlUserRepository) GetAll(ctx context.Context, search string, orderBy string, sort string, limit int32, offset int32) ([]sqlc.User, error) {
	var (
		users []sqlc.User
		err error
	)
	
	switch {
		case orderBy == "user_id" && sort == "asc":
			users, err = ur.db.ListUsersUserIdAsc(ctx, sqlc.ListUsersUserIdAscParams{
				Limit: limit,
				Offset: offset,
				Search: search,
			})
		case orderBy == "user_id" && sort == "desc":
			users, err = ur.db.ListUsersUserIdDesc(ctx, sqlc.ListUsersUserIdDescParams{
				Limit: limit,
				Offset: offset,
				Search: search,
			})
		case orderBy == "user_created_at" && sort == "asc":
			users, err = ur.db.ListUsersCreatedAtAsc(ctx, sqlc.ListUsersCreatedAtAscParams{
				Limit: limit,
				Offset: offset,
				Search: search,
			})
		case orderBy == "user_created_at" && sort == "desc":
			users, err = ur.db.ListUsersCreatedAtDesc(ctx, sqlc.ListUsersCreatedAtDescParams{
				Limit: limit,
				Offset: offset,
				Search: search,
			})
		default:
			users, err = ur.db.ListUsersUserIdAsc(ctx, sqlc.ListUsersUserIdAscParams{
				Limit: limit,
				Offset: offset,
				Search: search,
			})
		
	}
	if err != nil {
		return []sqlc.User{}, err
	}
	return users, nil
}

func (ur *SqlUserRepository) CountUsers(ctx context.Context, search string) (int64, error) {
	total, err := ur.db.CountUsers(ctx, search)
	if err != nil {
		return 0, err
	}
	return total, nil
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
