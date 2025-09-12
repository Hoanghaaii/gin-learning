package repository

import "hari-donganh/gin-learning/internal/domain/entities"

type UserRepository interface {
	Create(user *entities.User) error
	FindByID(id uint) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id uint) error
	FindAll(limit, offset int) ([]*entities.User, int64, error)
	FindByRole(role entities.UserRole, limit, offset int) ([]*entities.User, error)
	FindActiveUser(limit, offset int) ([]*entities.User, error)
	Search(query string, limit, offset int) ([]*entities.User, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByPhone(phone string) (bool, error)
	CountByRole(role entities.UserRole) (int64, error)
	CountActiveUsers() (int64, error)
	ActivateUser(userID uint) error
	DeactivateUser(userID uint) error
}
