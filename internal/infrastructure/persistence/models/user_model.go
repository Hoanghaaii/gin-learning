package models

import (
	"strings"

	"gorm.io/gorm"

	"hari-donganh/gin-learning/internal/domain/entities"
)

type UserModel struct {
	Basemodel
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Phone    string `gorm:"size:20;uniqueIndex;not null" json:"phone"`
	Role     string `gorm:"size:20;not null;default:'buyer'" json:"role"`
	IsActive bool   `gorm:"default:true" json:"is_active"`

	// Products     []ProductModel     `gorm:"foreignKey:SellerID;references:ID"`
	// Orders       []OrderModel       `gorm:"foreignKey:UserID;references:ID"`
	// Participants []EventParticipantModel `gorm:"foreignKey:UserID;references:ID"`
	// Notifications []NotificationModel `gorm:"foreignKey:UserID;references:ID"`
}

func (UserModel) TableName() string {
	return "users"
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Name = strings.TrimSpace(u.Name)
	return nil
}

func (u *UserModel) BeforeUpdate(tx *gorm.DB) error {
	if u.Name != "" {
		u.Name = strings.TrimSpace(u.Name)
	}
	return nil
}

func (u *UserModel) ToDomain() *entities.User {
	return &entities.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		Phone:     u.Phone,
		Role:      entities.UserRole(u.Role),
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *UserModel) FromDomain(user *entities.User) {
	u.ID = user.ID
	u.Name = user.Name
	u.Email = user.Email
	u.Password = user.Password
	u.Phone = user.Phone
	u.Role = string(user.Role)
	u.IsActive = user.IsActive
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
}

func NewUserModelFromDomain(user *entities.User) *UserModel {
	UserModel := &UserModel{}
	UserModel.FromDomain(user)
	return UserModel
}

func (UserModel) ScopeActive(db *gorm.DB) *gorm.DB {
	return db.Where("is_active = ?", true)
}

func (UserModel) ScopeInactive(db *gorm.DB) *gorm.DB {
	return db.Where("is_active = ?", false)
}

func (UserModel) ScopeByRole(role string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("role = ?", role)
	}
}
