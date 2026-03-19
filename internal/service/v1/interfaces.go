package v1service

import (
	"user-management-api/internal/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers(search string, page, limit int) 
	CreateUser(ctx *gin.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	GetUserByUUID(uuid string) 
	UpdateUser(ctx *gin.Context, userParams sqlc.UpdateUserParams) (sqlc.User, error) 
	SoftDeleteUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error)
	RestoreUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error) 
	DeleteUser(ctx *gin.Context, userUuid uuid.UUID) error
}
