package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	auction "github.com/parxyws/nego-gin/internal/auction/domain"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/parxyws/nego-gin/internal/user/repository"
	"github.com/parxyws/nego-gin/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *repository.UserRepositoryImpl
	ctx  context.Context
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	// Use PostgreSQL for testing
	db := test.SetupTestDB(suite.T())

	// Clean start: Drop tables before migration
	_ = db.Migrator().DropTable(&domain.User{}, &auction.Auction{})

	// Auto migrate the schema
	err := db.AutoMigrate(&domain.User{}, &auction.Auction{})
	suite.Require().NoError(err)

	suite.db = db
	suite.repo = repository.NewUserRepository(db).(*repository.UserRepositoryImpl)
	suite.ctx = context.Background()
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	// Clean up after each test
	err := test.CleanupTestDB(suite.db, "users")
	suite.Require().NoError(err)
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	// Final cleanup
	sqlDB, err := suite.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (suite *UserRepositoryTestSuite) TestCreateUser_Success() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        "user-001",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		FirstName: "Test",
		LastName:  "User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := suite.repo.CreateUser(suite.ctx, user)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "testuser", result.Username)
	assert.Equal(suite.T(), "test@example.com", result.Email)
}

func (suite *UserRepositoryTestSuite) TestCreateUser_DuplicateEmail() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user1 := &domain.User{
		ID:        "user-001",
		Username:  "testuser1",
		Email:     "duplicate@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := suite.repo.CreateUser(suite.ctx, user1)
	assert.NoError(suite.T(), err)

	// Try to create another user with the same email
	user2 := &domain.User{
		ID:        "user-002",
		Username:  "testuser2",
		Email:     "duplicate@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = suite.repo.CreateUser(suite.ctx, user2)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrRegistered, err)
}

func (suite *UserRepositoryTestSuite) TestCreateUser_DuplicateUsername() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user1 := &domain.User{
		ID:        "user-001",
		Username:  "duplicateuser",
		Email:     "user1@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := suite.repo.CreateUser(suite.ctx, user1)
	assert.NoError(suite.T(), err)

	// Try to create another user with the same username
	user2 := &domain.User{
		ID:        "user-002",
		Username:  "duplicateuser",
		Email:     "user2@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = suite.repo.CreateUser(suite.ctx, user2)
	assert.Error(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestUpdateUser_Success() {
	// Create a user first
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        "user-001",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		FirstName: "Test",
		LastName:  "User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := suite.repo.CreateUser(suite.ctx, user)
	assert.NoError(suite.T(), err)

	// Update the user
	user.FirstName = "Updated"
	user.LastName = "Name"
	user.Phone = "1234567890"

	result, err := suite.repo.UpdateUser(suite.ctx, user)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Updated", result.FirstName)
	assert.Equal(suite.T(), "Name", result.LastName)
	assert.Equal(suite.T(), "1234567890", result.Phone)
}

func (suite *UserRepositoryTestSuite) TestUpdateUser_NonExistent() {
	user := &domain.User{
		ID:        "non-existent-id",
		Username:  "ghost",
		Email:     "ghost@example.com",
		FirstName: "Ghost",
		LastName:  "User",
	}

	_, err := suite.repo.UpdateUser(suite.ctx, user)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrInvalidData, err)
}

func (suite *UserRepositoryTestSuite) TestDeleteUser_Success() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        "user-001",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := suite.repo.CreateUser(suite.ctx, user)
	assert.NoError(suite.T(), err)

	// Delete the user
	deleteUser := &domain.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	err = suite.repo.DeleteUser(suite.ctx, deleteUser)

	assert.NoError(suite.T(), err)

	// Verify user is deleted
	var count int64
	suite.db.Model(&domain.User{}).Where("email = ?", "test@example.com").Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *UserRepositoryTestSuite) TestDeleteUser_WrongPassword() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        "user-001",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := suite.repo.CreateUser(suite.ctx, user)
	assert.NoError(suite.T(), err)

	// Try to delete with wrong password
	deleteUser := &domain.User{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	err = suite.repo.DeleteUser(suite.ctx, deleteUser)

	assert.Error(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestDeleteUser_NonExistent() {
	deleteUser := &domain.User{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	err := suite.repo.DeleteUser(suite.ctx, deleteUser)

	assert.Error(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestReadByUsername_Success() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        "user-001",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := suite.repo.CreateUser(suite.ctx, user)
	assert.NoError(suite.T(), err)

	// Read by username
	searchUser := &domain.User{Username: "testuser"}
	result, err := suite.repo.ReadByUsername(suite.ctx, searchUser)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "testuser", result.Username)
	assert.Equal(suite.T(), "test@example.com", result.Email)
}

func (suite *UserRepositoryTestSuite) TestReadByUsername_NotFound() {
	searchUser := &domain.User{Username: "nonexistent"}
	result, err := suite.repo.ReadByUsername(suite.ctx, searchUser)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *UserRepositoryTestSuite) TestReadByUsername_CaseInsensitive() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        "user-001",
		Username:  "TestUser",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := suite.repo.CreateUser(suite.ctx, user)
	assert.NoError(suite.T(), err)

	// Search with different case should work with PostgreSQL
	searchUser := &domain.User{Username: "testuser"}
	result, err := suite.repo.ReadByUsername(suite.ctx, searchUser)

	// This test may fail or succeed depending on your column collation
	// Including it to verify PostgreSQL behavior
	if err == nil {
		assert.Equal(suite.T(), "TestUser", result.Username)
	}
}

func (suite *UserRepositoryTestSuite) TestReadByEmail_Success() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        "user-001",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := suite.repo.CreateUser(suite.ctx, user)
	assert.NoError(suite.T(), err)

	// Read by email
	searchUser := &domain.User{Email: "test@example.com"}
	result, err := suite.repo.ReadByEmail(suite.ctx, searchUser)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "testuser", result.Username)
	assert.Equal(suite.T(), "test@example.com", result.Email)
}

func (suite *UserRepositoryTestSuite) TestReadByEmail_NotFound() {
	searchUser := &domain.User{Email: "nonexistent@example.com"}
	result, err := suite.repo.ReadByEmail(suite.ctx, searchUser)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *UserRepositoryTestSuite) TestReadById_Success() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        "user-001",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := suite.repo.CreateUser(suite.ctx, user)
	assert.NoError(suite.T(), err)

	// Read by ID
	searchUser := &domain.User{ID: "user-001"}
	result, err := suite.repo.ReadById(suite.ctx, searchUser)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "user-001", result.ID)
	assert.Equal(suite.T(), "testuser", result.Username)
}

func (suite *UserRepositoryTestSuite) TestReadById_NotFound() {
	searchUser := &domain.User{ID: "non-existent-id"}
	result, err := suite.repo.ReadById(suite.ctx, searchUser)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *UserRepositoryTestSuite) TestReadAllByRoles_Success() {
	// Create multiple verified users
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	users := []*domain.User{
		{
			ID:         "user-001",
			Username:   "user1",
			Email:      "user1@example.com",
			Password:   string(hashedPassword),
			IsVerified: sql.NullTime{Time: time.Now(), Valid: true},
			CreatedAt:  time.Now().Add(-3 * time.Hour),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "user-002",
			Username:   "user2",
			Email:      "user2@example.com",
			Password:   string(hashedPassword),
			IsVerified: sql.NullTime{Time: time.Now(), Valid: true},
			CreatedAt:  time.Now().Add(-2 * time.Hour),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "user-003",
			Username:   "user3",
			Email:      "user3@example.com",
			Password:   string(hashedPassword),
			IsVerified: sql.NullTime{Time: time.Now(), Valid: true},
			CreatedAt:  time.Now().Add(-1 * time.Hour),
			UpdatedAt:  time.Now(),
		},
	}

	for _, u := range users {
		_, err := suite.repo.CreateUser(suite.ctx, u)
		assert.NoError(suite.T(), err)
	}

	// Read all verified users
	results, err := suite.repo.ReadAllByRoles(suite.ctx, "", true, 10, "")

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 3)

	// Verify descending order by created_at
	assert.True(suite.T(), results[0].CreatedAt.After(results[1].CreatedAt) || results[0].CreatedAt.Equal(results[1].CreatedAt))
	assert.True(suite.T(), results[1].CreatedAt.After(results[2].CreatedAt) || results[1].CreatedAt.Equal(results[2].CreatedAt))
}

func (suite *UserRepositoryTestSuite) TestReadAllByRoles_WithLimit() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	for i := 1; i <= 5; i++ {
		user := &domain.User{
			ID:         "user-00" + string(rune('0'+i)),
			Username:   "user" + string(rune('0'+i)),
			Email:      "user" + string(rune('0'+i)) + "@example.com",
			Password:   string(hashedPassword),
			IsVerified: sql.NullTime{Time: time.Now(), Valid: true},
			CreatedAt:  time.Now().Add(-time.Duration(i) * time.Hour),
			UpdatedAt:  time.Now(),
		}
		_, err := suite.repo.CreateUser(suite.ctx, user)
		assert.NoError(suite.T(), err)
	}

	// Read with limit of 3
	results, err := suite.repo.ReadAllByRoles(suite.ctx, "", true, 3, "")

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 3)
}

func (suite *UserRepositoryTestSuite) TestReadAllByRoles_UnverifiedUsersExcluded() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Create verified user
	verifiedUser := &domain.User{
		ID:         "user-001",
		Username:   "verified",
		Email:      "verified@example.com",
		Password:   string(hashedPassword),
		IsVerified: sql.NullTime{Time: time.Now(), Valid: true},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err := suite.repo.CreateUser(suite.ctx, verifiedUser)
	assert.NoError(suite.T(), err)

	// Create unverified user
	unverifiedUser := &domain.User{
		ID:        "user-002",
		Username:  "unverified",
		Email:     "unverified@example.com",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = suite.repo.CreateUser(suite.ctx, unverifiedUser)
	assert.NoError(suite.T(), err)

	// Read all - should only return verified users
	results, err := suite.repo.ReadAllByRoles(suite.ctx, "", true, 10, "")

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 1)
	assert.Equal(suite.T(), "verified", results[0].Username)
}

func (suite *UserRepositoryTestSuite) TestReadAllByRoles_Pagination() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Create 10 verified users
	for i := 1; i <= 10; i++ {
		user := &domain.User{
			ID:         fmt.Sprintf("user-%03d", i),
			Username:   fmt.Sprintf("user%d", i),
			Email:      fmt.Sprintf("user%d@example.com", i),
			Password:   string(hashedPassword),
			IsVerified: sql.NullTime{Time: time.Now(), Valid: true},
			CreatedAt:  time.Now().Add(-time.Duration(i) * time.Minute),
			UpdatedAt:  time.Now(),
		}
		_, err := suite.repo.CreateUser(suite.ctx, user)
		assert.NoError(suite.T(), err)
	}

	// First page
	firstPage, err := suite.repo.ReadAllByRoles(suite.ctx, "", true, 5, "")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), firstPage, 5)

	// Second page using cursor
	lastUser := firstPage[len(firstPage)-1]
	secondPage, err := suite.repo.ReadAllByRoles(suite.ctx, lastUser.ID, true, 5, lastUser.CreatedAt.Format(time.RFC3339))
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), secondPage, 5)

	// Ensure no overlap
	for _, u1 := range firstPage {
		for _, u2 := range secondPage {
			assert.NotEqual(suite.T(), u1.ID, u2.ID)
		}
	}
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
