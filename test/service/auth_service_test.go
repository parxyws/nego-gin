package service_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/internal/user/service"
	"github.com/parxyws/nego-gin/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, entity *domain.User) (*domain.User, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, entity *domain.User) (*domain.User, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, entity *domain.User) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockUserRepository) ReadByUsername(ctx context.Context, entity *domain.User) (*domain.User, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) ReadByEmail(ctx context.Context, entity *domain.User) (*domain.User, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) ReadById(ctx context.Context, entity *domain.User) (*domain.User, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) ReadAllByRoles(ctx context.Context, id string, sortOrderDesc bool, limit int, createdAt string) ([]domain.User, error) {
	args := m.Called(ctx, id, sortOrderDesc, limit, createdAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

// MockUserRoleRepository is a mock implementation of UserRoleRepository
type MockUserRoleRepository struct {
	mock.Mock
}

func (m *MockUserRoleRepository) CreateUserRole(ctx context.Context, entity *domain.UserRole) (*domain.UserRole, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserRole), args.Error(1)
}

func (m *MockUserRoleRepository) DeleteUserRole(ctx context.Context, entity []domain.UserRole) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockUserRoleRepository) ReadAllUserRole(ctx context.Context) ([]domain.UserRole, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserRole), args.Error(1)
}

func (m *MockUserRoleRepository) CountUserRoleByRoleID(ctx context.Context, entity []int32) (map[int32]int64, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[int32]int64), args.Error(1)
}

func (m *MockUserRoleRepository) ReadUserRoleByUserID(ctx context.Context, entity *domain.UserRole) ([]domain.UserRole, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserRole), args.Error(1)
}

func TestAuthService_Register_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	mockRoleRepo := new(MockUserRoleRepository)
	mockRedis := test.SetupTestRedis(t)
	defer test.CleanupTestRedis(mockRedis)

	cfg, _ := test.InitEnv()

	authService := service.NewAuthService(cfg, mockRepo, mockRoleRepo, mockRedis)

	ctx := context.Background()
	request := &dto.UserRegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Email:     "john@example.com",
		Password:  "password123",
	}

	// Mock expectations
	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*domain.User")).Return(&domain.User{
		ID:        "test-ulid-123",
		Username:  request.Username,
		Email:     request.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	// Execute
	response, err := authService.Register(ctx, request)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.ReferenceID)

	expectedIdentifier := base64.StdEncoding.EncodeToString([]byte(request.Email))
	assert.Equal(t, expectedIdentifier, response.ReferenceID)

	// Verify Redis data
	otpKey := "otp" + response.ReferenceID
	userRefKey := "user-ref" + response.ReferenceID

	val, err := mockRedis.Get(ctx, otpKey).Result()
	assert.NoError(t, err)
	assert.NotEmpty(t, val)

	val, err = mockRedis.Get(ctx, userRefKey).Result()
	assert.NoError(t, err)
	assert.Equal(t, request.Email, val)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRoleRepo := new(MockUserRoleRepository)
	mockRedis := test.SetupTestRedis(t)
	defer test.CleanupTestRedis(mockRedis)

	cfg, _ := test.InitEnv()
	authService := service.NewAuthService(cfg, mockRepo, mockRoleRepo, mockRedis)

	ctx := context.Background()
	email := "test@example.com"
	refID := base64.StdEncoding.EncodeToString([]byte(email))
	otp := "123456"

	// Preset Redis
	mockRedis.Set(ctx, "otp"+refID, otp, 10*time.Minute)
	mockRedis.Set(ctx, "user-ref"+refID, email, 10*time.Minute)

	req := &dto.UserValidateAccRequest{
		ReferenceID: refID,
		OTP:         otp,
	}

	// Mock expectations
	mockRepo.On("ReadByEmail", ctx, mock.MatchedBy(func(u *domain.User) bool {
		return u.Email == email
	})).Return(&domain.User{ID: "user-123", Email: email}, nil)

	mockRepo.On("UpdateUser", ctx, mock.MatchedBy(func(u *domain.User) bool {
		return u.ID == "user-123" && u.IsVerified.Valid == true
	})).Return(&domain.User{ID: "user-123"}, nil)

	// Execute
	err := authService.ValidateUser(ctx, req)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_PasswordHashing(t *testing.T) {
	t.Run("Password is hashed correctly", func(t *testing.T) {
		password := "testpassword123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, string(hashedPassword))

		// Verify the hashed password can be compared
		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
		assert.NoError(t, err)
	})

	t.Run("Wrong password fails comparison", func(t *testing.T) {
		password := "correctpassword"
		wrongPassword := "wrongpassword"

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(wrongPassword))

		assert.Error(t, err)
	})
}

func TestAuthService_EdgeCases(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRoleRepo := new(MockUserRoleRepository)
	mockRedis := test.SetupTestRedis(t)
	defer test.CleanupTestRedis(mockRedis)

	cfg, _ := test.InitEnv()
	authService := service.NewAuthService(cfg, mockRepo, mockRoleRepo, mockRedis)
	ctx := context.Background()

	t.Run("ValidateUser_InvalidOTP", func(t *testing.T) {
		email := "test@example.com"
		refID := base64.StdEncoding.EncodeToString([]byte(email))
		mockRedis.Set(ctx, "otp"+refID, "111111", 10*time.Minute)

		req := &dto.UserValidateAccRequest{
			ReferenceID: refID,
			OTP:         "222222",
		}

		err := authService.ValidateUser(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("Register_RepoError", func(t *testing.T) {
		req := &dto.UserRegisterRequest{
			Username: "erroruser",
			Email:    "error@example.com",
			Password: "password123",
		}

		mockRepo.On("CreateUser", ctx, mock.Anything).Return(nil, fmt.Errorf("database fault")).Once()

		resp, err := authService.Register(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "database fault")
	})
}
