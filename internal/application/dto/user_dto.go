package dto

import (
	"time"

	"hari-donganh/gin-learning/internal/domain/entities"
)

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
	Phone    string
	Role     entities.UserRole
}

type UpdateUserRequest struct {
	Name  string
	Phone string
}

type ChangePasswordRequest struct {
	CurrentPassword string
	NewPassword     string
	ConfirmPassword string
}

type ChangeRoleRequest struct {
	Role entities.UserRole
}

type LoginRequest struct {
	Email    string
	Password string
}

type UserResponse struct {
	ID        uint
	Name      string
	Email     string
	Phone     string
	Role      entities.UserRole
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserListResponse struct {
	Users      []UserResponse
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

type LoginResponse struct {
	User         UserResponse
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type UserQueryParams struct {
	Page     int
	Limit    int
	Role     entities.UserRole
	IsActive *bool
	Search   string
	SortBy   string
	Order    string
}

func ToUserResponse(user *entities.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserListResponse(users []*entities.User, total int64, page, limit int) UserListResponse {
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(user)
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

func ToLoginResponse(user *entities.User, token string, expiresAt time.Time) LoginResponse {
	return LoginResponse{
		User:        ToUserResponse(user),
		AccessToken: token,
		ExpiresAt:   expiresAt,
	}
}
