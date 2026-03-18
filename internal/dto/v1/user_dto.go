package v1dto

import (
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/utils"
)

type UserDTO struct {
	UUID     string `json:"uuid"`
	Name     string `json:"full_name"`
	Email    string `json:"email_address"`
	Age      *int    `json:"age"`
	Status   string `json:"status"`
	Level    string `json:"level"`
	CreatedAt string `json:"created_at"`
	// DeletedAt string `json:"deleted_at"`
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email,email_advanced"`
	Age      int    `json:"age" binding:"gt=0"`
	Password string `json:"password" binding:"required,min=8,password_strong"`
	Status   int    `json:"status" binding:"required,oneof=1 2"`
	Level    int    `json:"level" binding:"required,oneof=1 2"`
}

type UpdateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email,email_advanced"`
	Age      int    `json:"age" binding:"required,gt=0"`
	Password string `json:"password" binding:"omitempty,min=8,password_strong"`
	Status   int    `json:"status" binding:"required,oneof=1 2"`
	Level    int    `json:"level" binding:"required,oneof=1 2"`
}

func (input *CreateUserInput) MapCreateInputToModel() sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		UserEmail: input.Email,
		UserPassword: input.Password,
		UserFullname: input.Name,
		UserAge: utils.ConvertToInt32Pointer(input.Age),
		UserStatus: int32(input.Status),
		UserLevel: int32(input.Level),
	}
	
}



func (input *UpdateUserInput) MapUpdateInputToModel()  {
	
}

func MapUserToDTO(user sqlc.User) *UserDTO {
	dto := &UserDTO{
		UUID: user.UserUuid.String(),
		Name: user.UserFullname,
		Email: user.UserEmail,
		Status: mapStatusText(int(user.UserStatus)),
		Level: mapLevelText(int(user.UserLevel)),
        CreatedAt: user.UserCreatedAt.Format(time.RFC3339),
	}

	if user.UserAge != nil {
		age := int(*user.UserAge)
		dto.Age = &age
	}

	// if user.UserDeletedAt.Valid {
	// 	dto.DeletedAt = user.UserDeletedAt.Time.Format(time.RFC3339)
	// }

	return dto
}

func mapStatusText(status int) string {
	switch status {
	case 1:
		return "Active"
	case 2:
		return "Inactive"
	case 3:
		return "Banned"
	default:
		return "Unknown"
	}
}

func mapLevelText(status int) string {
	switch status {
	case 1:
		return "Administrator"
	case 2:
		return "Moderator"
	case 3:
		return "Member"
	default:
		return "Unknown"
	}
}
