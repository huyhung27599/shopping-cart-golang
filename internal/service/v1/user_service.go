package v1service

import (
	"database/sql"
	"errors"
	"strconv"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAllUsers(ctx *gin.Context, search string, orderBy string, sort string, page int32, limit int32) ([]sqlc.User, int64, error)  {
	context := ctx.Request.Context()

	if sort == "" {
		sort = "desc"
	}
	if orderBy == "" {
		orderBy = "user_created_at"
	}
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		envLimit := utils.GetEnv("LIMIT_ITEMS_ON_PER_PAGE", "10")
		limitInt, err := strconv.Atoi(envLimit)
		if err != nil || limitInt <= 0 {
			limitInt=10
		}
		limit = int32(limitInt)
	}

	offset := (page - 1) * limit

	users, err := us.repo.GetAll(context, search, orderBy, sort, limit, offset)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError(err, "Failed to get users", utils.ErrCodeInternal)
	}

	total, err := us.repo.CountUsers(context, search)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError(err, "Failed to count users", utils.ErrCodeInternal)
	}
	return users, total, nil
	
	
}

func (us *userService) CreateUser(ctx *gin.Context, userParams sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	userParams.UserEmail = utils.NormalizeString(userParams.UserEmail)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userParams.UserPassword), bcrypt.DefaultCost)
	if err != nil {
	
		return sqlc.User{}, utils.WrapError(err, "Failed to hash password", utils.ErrCodeInternal)
	}
	userParams.UserPassword = string(hashedPassword)

 	user, err := us.repo.Create(context, userParams)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return sqlc.User{}, utils.NewError( "User already exists", utils.ErrCodeConflict)
			}
		}
		return sqlc.User{}, utils.WrapError(err, "Failed to create user", utils.ErrCodeInternal)
	}
	return user,nil
	
}

func (us *userService) GetUserByUUID(uuid string)  {
	
}

func (us *userService) UpdateUser(ctx *gin.Context, userParams sqlc.UpdateUserParams) (sqlc.User, error)  {
	context := ctx.Request.Context()

	if userParams.UserPassword != nil && *userParams.UserPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*userParams.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			return sqlc.User{}, utils.WrapError(err, "Failed to hash password", utils.ErrCodeInternal)
		}
		hashedPasswordStr := string(hashedPassword)
		userParams.UserPassword = &hashedPasswordStr
	}

	user, err := us.repo.Update(context, userParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "Failed to update user", utils.ErrCodeInternal)
	}
	return user, nil
}

func (us *userService) DeleteUser(ctx *gin.Context, userUuid uuid.UUID)  error  {
	context := ctx.Request.Context()
	err := us.repo.Delete(context, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("User not found", utils.ErrCodeNotFound)
		}
		return utils.WrapError(err, "Failed to delete user", utils.ErrCodeInternal)
	}
	return nil
}

func (us *userService) SoftDeleteUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error)  {
	context := ctx.Request.Context()

	user, err := us.repo.SoftDelete(context, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "Failed to soft delete user", utils.ErrCodeInternal)
	}
	return user, nil

	
}

func (us *userService) RestoreUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error)  {
	context := ctx.Request.Context()

	user, err := us.repo.Restore(context, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "Failed to restore user", utils.ErrCodeInternal)
	}
	return user, nil
}
