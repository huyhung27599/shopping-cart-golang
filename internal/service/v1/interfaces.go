package v1service

import (
	"user-management-api/internal/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers(ctx *gin.Context, search string, orderBy string, sort string, page int32, limit int32, deleted bool) ([]sqlc.User, int64, error)
	CreateUser(ctx *gin.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	GetUserByUUID(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error)
	UpdateUser(ctx *gin.Context, userParams sqlc.UpdateUserParams) (sqlc.User, error) 
	SoftDeleteUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error)
	RestoreUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error) 
	DeleteUser(ctx *gin.Context, userUuid uuid.UUID) error
}

type AuthService interface {
	Login(ctx *gin.Context, email string, password string)  (string, string, int, error)
	Logout(ctx *gin.Context, refreshToken string) error
	RefreshToken(ctx *gin.Context, refreshToken string) (string, string, int, error)
	RequestForgotPassword(ctx *gin.Context, email string) error
	ResetPassword(ctx *gin.Context, token string, newPassword string) error
}