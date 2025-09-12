package usecase

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"hari-donganh/gin-learning/internal/application/dto"
	"hari-donganh/gin-learning/internal/application/repository"
	"hari-donganh/gin-learning/internal/domain/entities"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := uc.validateCreateUserRequest(req); err != nil {
		return nil, err
	}
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}
	exists, err = uc.userRepo.ExistsByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("phone number already exists")
	}

	hashedPassword, err := uc.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &entities.User{
		Name:      strings.TrimSpace(req.Name),
		Email:     strings.ToLower(strings.TrimSpace(req.Email)),
		Password:  hashedPassword,
		Phone:     uc.cleanPhoneNumber(req.Phone),
		Role:      req.Role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	response := dto.ToUserResponse(user)
	return &response, nil
}

func (uc *UserUseCase) GetUserById(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := dto.ToUserResponse(user)
	return &response, nil
}

func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	response := dto.ToUserResponse(user)
	return &response, nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, userID uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if req.Name != "" {
		if err := uc.validateName(req.Name); err != nil {
			return nil, err
		}
		user.Name = strings.TrimSpace(req.Name)
	}
	if req.Phone != "" {
		if err := uc.validatePhone(req.Phone); err != nil {
			return nil, err
		}
		// Check if phone exists for other users
		existingUser, _ := uc.userRepo.FindByID(ctx, userID) // current user
		if req.Phone != existingUser.Phone {
			exists, err := uc.userRepo.ExistsByPhone(ctx, req.Phone)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("phone number already exists")
			}
		}
		user.Phone = uc.cleanPhoneNumber(req.Phone)
	}
	user.UpdatedAt = time.Now()
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	response := dto.ToUserResponse(user)
	return &response, nil
}

func (uc *UserUseCase) ChangePassword(ctx context.Context, userID uint, req dto.ChangePasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("password confirmation does not match")
	}
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}
	if err := uc.validatePassword(req.NewPassword); err != nil {
		return err
	}
	hashedPassword, err := uc.hashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to process new password")
	}
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()
	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUseCase) ChangeRole(ctx context.Context, userID uint, req dto.ChangeRoleRequest, performedByUserID uint) error {
	canManage, err := uc.CanUserManageUsers(ctx, performedByUserID)
	if err != nil {
		return err
	}
	if !canManage {
		return errors.New("only admin can change role users")
	}
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	user.Role = req.Role
	user.UpdatedAt = time.Now()
	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUseCase) GetAllUsers(ctx context.Context, params dto.UserQueryParams) (*dto.UserListResponse, error) {
	offset := (params.Page - 1) * params.Limit
	var users []*entities.User
	var total int64
	var err error
	if params.Role != "" {
		users, err = uc.userRepo.FindByRole(ctx, params.Role, params.Limit, offset)
		if err != nil {
			return nil, err
		}
		total, err = uc.userRepo.CountByRole(ctx, params.Role)
		if err != nil {
			return nil, err
		}
	} else if params.IsActive != nil && *params.IsActive {
		users, err = uc.userRepo.FindActiveUser(ctx, params.Limit, offset)
		if err != nil {
			return nil, err
		}
		total, err = uc.userRepo.CountActiveUsers(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		users, total, err = uc.userRepo.FindAll(ctx, params.Limit, offset)
		if err != nil {
			return nil, err
		}
	}
	response := dto.ToUserListResponse(users, total, params.Page, params.Limit)
	return &response, nil
}

func (uc *UserUseCase) DeactivateUser(ctx context.Context, userID uint, performedByUserID uint) error {
	canManage, err := uc.CanUserManageUsers(ctx, performedByUserID)
	if err != nil {
		return err
	}
	if !canManage {
		return errors.New("only admin can deactivate users")
	}
	targetUser, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if !targetUser.IsActive {
		return errors.New("user is already deactivated")
	}
	return uc.userRepo.DeactivateUser(ctx, userID)
}

func (uc *UserUseCase) ActivateUser(ctx context.Context, userID uint, performedByUserID uint) error {
	canManage, err := uc.CanUserManageUsers(ctx, performedByUserID)
	if err != nil {
		return err
	}
	if !canManage {
		return errors.New("only admin can activate users")
	}
	targetUser, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if targetUser.IsActive {
		return errors.New("user is already activated")
	}
	return uc.userRepo.ActivateUser(ctx, userID)
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, userID uint, performedByUserID uint) error {
	canManage, err := uc.CanUserManageUsers(ctx, performedByUserID)
	if err != nil {
		return err
	}
	if !canManage {
		return errors.New("only admin can delete users")
	}

	// Check if user exists
	_, err = uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	return uc.userRepo.Delete(ctx, userID)
}

func (uc *UserUseCase) CanUserCreateProduct(ctx context.Context, userID uint) (bool, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return user.Role == entities.RoleBuyer || user.Role == entities.RoleAdmin, nil
}

func (uc *UserUseCase) CanUserManageUsers(ctx context.Context, userID uint) (bool, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}
	if user.Role != entities.RoleAdmin {
		return false, errors.New("only admin can activate users")
	}
	return true, nil
}

//Private helper methods (không cần ctx)

func (uc *UserUseCase) validateCreateUserRequest(req dto.CreateUserRequest) error {
	if err := uc.validateName(req.Name); err != nil {
		return err
	}
	if err := uc.validateEmail(req.Email); err != nil {
		return err
	}
	if err := uc.validatePassword(req.Password); err != nil {
		return err
	}
	if err := uc.validatePhone(req.Phone); err != nil {
		return err
	}
	return uc.validateRole(req.Role)
}

func (uc *UserUseCase) validateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if len(name) < 2 || len(name) > 100 {
		return errors.New("name must be between 2 and 100 characters")
	}
	nameRegex := regexp.MustCompile(`^[a-zA-ZÀ-ỹ\s\-'\.]+$`)
	if !nameRegex.MatchString(name) {
		return errors.New("name contains invalid characters")
	}
	return nil
}

func (uc *UserUseCase) validateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	if len(email) > 254 {
		return errors.New("email too long")
	}
	return nil
}

func (uc *UserUseCase) validatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if len(password) < 8 || len(password) > 128 {
		return errors.New("password must be in range 8 and 128 characters")
	}
	if !uc.isValidPasswordComplexity(password) {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character")
	}
	return nil
}

func (uc *UserUseCase) validatePhone(phone string) error {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return errors.New("phone cannot be empty")
	}
	cleanPhone := uc.cleanPhoneNumber(phone)
	phoneRegex := regexp.MustCompile(`^(\+84|84|0)[1-9][0-9]{8,9}$`)
	if !phoneRegex.MatchString(cleanPhone) {
		return errors.New("invalid phone number format")
	}
	return nil
}

func (uc *UserUseCase) validateRole(role entities.UserRole) error {
	switch role {
	case entities.RoleBuyer, entities.RoleSeller, entities.RoleAdmin:
		return nil
	default:
		return errors.New("invalid user role")
	}
}

func (uc *UserUseCase) isValidPasswordComplexity(password string) bool {
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}

func (uc *UserUseCase) cleanPhoneNumber(phone string) string {
	// Remove common phone formatting characters
	cleanPhone := regexp.MustCompile(`[\s\-\(\)\.]+`).ReplaceAllString(phone, "")
	return cleanPhone
}

func (uc *UserUseCase) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
