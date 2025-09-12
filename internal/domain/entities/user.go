package entities

import "time"

type UserRole string

const (
	RoleBuyer  UserRole = "buyer"
	RoleSeller UserRole = "seller"
	RoleAdmin  UserRole = "admin"
)

type User struct {
	ID        uint
	Name      string
	Email     string
	Password  string
	Phone     string
	Role      UserRole
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
