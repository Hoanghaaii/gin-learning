package repository

import (
	"context"

	"hari-donganh/gin-learning/internal/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id uint) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uint) error

	FindAll(ctx context.Context, limit, offset int) ([]*entities.User, int64, error)
	FindByRole(ctx context.Context, role entities.UserRole, limit, offset int) ([]*entities.User, error)
	FindActiveUser(ctx context.Context, limit, offset int) ([]*entities.User, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.User, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)

	CountByRole(ctx context.Context, role entities.UserRole) (int64, error)
	CountActiveUsers(ctx context.Context) (int64, error)

	ActivateUser(ctx context.Context, userID uint) error
	DeactivateUser(ctx context.Context, userID uint) error
}
