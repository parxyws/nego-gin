package repository_test

import (
	"context"
	"testing"

	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/parxyws/nego-gin/internal/user/repository"
	"github.com/parxyws/nego-gin/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RoleRepositoryTestSuite struct {
	suite.Suite
	db       *gorm.DB
	roleRepo *repository.RoleRepositoryImpl
	urRepo   *repository.UserRoleRepositoryImpl
	ctx      context.Context
}

func (suite *RoleRepositoryTestSuite) SetupSuite() {
	db := test.SetupTestDB(suite.T())

	// Clean start
	_ = db.Migrator().DropTable(&domain.UserRole{}, &domain.Role{}, &domain.User{})

	// Auto migrate
	err := db.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.UserRole{})
	suite.Require().NoError(err)

	suite.db = db
	suite.roleRepo = repository.NewRoleRepository(db).(*repository.RoleRepositoryImpl)
	suite.urRepo = repository.NewUserRoleRepository(db).(*repository.UserRoleRepositoryImpl)
	suite.ctx = context.Background()
}

func (suite *RoleRepositoryTestSuite) TearDownTest() {
	// Clean up tables
	test.CleanupTestDB(suite.db, "user_roles", "roles", "users")
}

func (suite *RoleRepositoryTestSuite) TestCreateRole_Success() {
	role := &domain.Role{
		RoleName: "ADMIN",
	}

	result, err := suite.roleRepo.CreateRole(suite.ctx, role)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "ADMIN", result.RoleName)
}

func (suite *RoleRepositoryTestSuite) TestCreateRole_FirstOrCreate() {
	role1 := &domain.Role{RoleName: "USER"}
	_, err := suite.roleRepo.CreateRole(suite.ctx, role1)
	assert.NoError(suite.T(), err)

	// Try to create the same role again
	role2 := &domain.Role{RoleName: "USER"}
	result, err := suite.roleRepo.CreateRole(suite.ctx, role2)

	assert.NoError(suite.T(), err) // FirstOrCreate shouldn't error on duplicate name
	assert.Equal(suite.T(), role1.ID, result.ID)
}

func (suite *RoleRepositoryTestSuite) TestReadAllRole() {
	roles := []domain.Role{
		{RoleName: "ROLE1"},
		{RoleName: "ROLE2"},
	}
	for _, r := range roles {
		_, _ = suite.roleRepo.CreateRole(suite.ctx, &r)
	}

	result, err := suite.roleRepo.ReadAllRole(suite.ctx)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
}

func (suite *RoleRepositoryTestSuite) TestUserRole_Lifecycle() {
	// 1. Create a User and a Role
	u := &domain.User{ID: "user-1", Username: "u1", Email: "e1@test.com"}
	suite.db.Create(u)

	r := &domain.Role{RoleName: "SUPER"}
	suite.db.FirstOrCreate(r)

	// 2. Create UserRole
	ur := &domain.UserRole{
		UserID: u.ID,
		RoleID: r.ID,
	}
	res, err := suite.urRepo.CreateUserRole(suite.ctx, ur)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res)

	// 3. Read back
	results, err := suite.urRepo.ReadUserRoleByUserID(suite.ctx, &domain.UserRole{UserID: u.ID})
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 1)
	for _, result := range results {
		assert.Equal(suite.T(), r.RoleName, result.Role.RoleName)
	}

	// 4. Count by Role ID
	counts, err := suite.urRepo.CountUserRoleByRoleID(suite.ctx, []int32{r.ID})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), counts[r.ID])

	// 5. Delete
	err = suite.urRepo.DeleteUserRole(suite.ctx, results)
	assert.NoError(suite.T(), err)

	// 6. Verify deleted
	finalResults, _ := suite.urRepo.ReadUserRoleByUserID(suite.ctx, &domain.UserRole{UserID: u.ID})
	assert.Len(suite.T(), finalResults, 0)
}

func TestRoleRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RoleRepositoryTestSuite))
}
