package v1service

import (
	"errors"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
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

func (us *userService) GetAllUsers(search string, page, limit int)  {
	
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

func (us *userService) UpdateUser(uuid string)  {
	
}

func (us *userService) DeleteUser(uuid string)  {
	
}
