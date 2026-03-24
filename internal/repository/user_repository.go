package repository

import (
	"context"
	"fmt"
	"user-management-api/internal/db"
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

func (ur *SqlUserRepository) GetAllV2(ctx context.Context, search string, orderBy string, sort string, limit int32, offset int32,deleted bool) ([]sqlc.User, error) {
	query := `SELECT * FROM users WHERE user_deleted_at IS NULL AND (
    $1::text IS NULL OR
    $1::text = '' OR
    user_email ILIKE '%' || $1 || '%' OR
    user_fullname ILIKE '%' || $1 || '%'
)`

if deleted {
	query += " AND user_deleted_at IS NOT NULL"
} else {
	query += " AND user_deleted_at IS NULL"
}
order := "ASC"
if sort == "desc" {
	order = "DESC"
}

 switch orderBy {
	case "user_id", "user_created_at":
		query += fmt.Sprintf(" ORDER BY %s %s", orderBy, order)
	default:
		query += " ORDER BY user_id ASC"
 }

 query += " LIMIT $2 OFFSET $3 -- name: Get All Version 2"

 rows,err := db.DBPool.Query(ctx, query, search, limit, offset)

 if err != nil {
	return nil, err
}
defer rows.Close()
users := []sqlc.User{}
for rows.Next() {
	var i sqlc.User
	if err := rows.Scan(
		&i.UserID,
		&i.UserUuid,
		&i.UserEmail,
		&i.UserPassword,
		&i.UserFullname,
		&i.UserAge,
		&i.UserStatus,
		&i.UserLevel,
		&i.UserDeletedAt,
		&i.UserCreatedAt,
		&i.UserUpdatedAt,
	); err != nil {
		return nil, err
	}
	users = append(users, i)
}
if err := rows.Err(); err != nil {
	return nil, err
}
 return users, nil

}

func (ur *SqlUserRepository) CountUsers(ctx context.Context, search string, deleted bool) (int64, error) {
	total, err := ur.db.CountUsers(ctx, sqlc.CountUsersParams{
		Search: search,
		Deleted: &deleted,
	})
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

func (ur *SqlUserRepository) GetByUUID(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error) {
	user, err := ur.db.GetUser(ctx, userUuid)
	if err != nil {
		return sqlc.User{}, err
	}
	return user, nil
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

func (ur *SqlUserRepository) FindByEmail(ctx context.Context, email string) (sqlc.User, error) {
	user, err := ur.db.GetUserByEmail(ctx, email)
	if err != nil {
		return sqlc.User{}, err
	}
	return user, nil
}
