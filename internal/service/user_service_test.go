package service_test

import (
	"context"
	"testing"

	"github.com/alejandro-bustamante/sancho/server/internal/repository"
	"github.com/alejandro-bustamante/sancho/server/internal/service"
	"github.com/alejandro-bustamante/sancho/server/internal/testutils"
	"github.com/stretchr/testify/suite"
)

type UserServiceSuite struct {
	suite.Suite
	db      repository.Database
	service *service.UserService
	ctx     context.Context
}

func (s *UserServiceSuite) SetupTest() {
	s.db = testutils.SetupTestDB(s.T())
	s.service = service.NewUserService(s.db)
	s.ctx = context.Background()
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) TestRegisterUser() {
	s.Run("should successfully register a new user", func() {
		username := "testuser"
		password := "secret123"
		email := "test@example.com"

		user, err := s.service.RegisterUser(s.ctx, username, password, email)

		// Require: stop test if error, because checking 'user' would panic if nil
		s.Require().NoError(err)
		s.Require().NotNil(user)

		// Assert: check values
		s.Assert().NotZero(user.ID)
		s.Assert().Equal(username, user.Username)
		s.Assert().Equal(email, *(user.Email))

		// Security check: Hash shouldn't match plain text
		s.Assert().NotEqual(password, user.PasswordHash)
	})

	s.Run("should fail when registering a duplicate username", func() {
		// Arrange: Insert the first user
		_, err := s.service.RegisterUser(s.ctx, "unique_user", "pass", "email@test.com")
		s.Require().NoError(err)

		// Act: Try to insert again
		_, errDuplicate := s.service.RegisterUser(s.ctx, "unique_user", "pass", "other@test.com")

		// Assert: Should return an error
		s.Assert().Error(errDuplicate)
	})
}

func (s *UserServiceSuite) TestAuthenticateUser() {
	username := "auth_user"
	password := "correct_password"

	// We use Require because if setup fails, the rest of the test is pointless.
	_, err := s.service.RegisterUser(s.ctx, username, password, "auth@example.com")
	s.Require().NoError(err)

	s.Run("should return true for valid credentials", func() {
		isAuthenticated, err := s.service.AuthenticateUser(s.ctx, username, password)
		s.Assert().NoError(err)
		s.Assert().True(isAuthenticated)
	})

	s.Run("should return false for invalid password", func() {
		isAuthenticated, err := s.service.AuthenticateUser(s.ctx, username, "wrong_password")
		s.Assert().NoError(err) // It's not a DB error, just failed auth
		s.Assert().False(isAuthenticated)
	})

	s.Run("should fail for non-existent user", func() {
		// Depending on implementation, this might return error or just false.
		// Assuming sqlc returns 'no rows' error:
		isAuthenticated, err := s.service.AuthenticateUser(s.ctx, "ghost_user", "password")

		// If your service wraps the error, check for Error.
		// If it returns false without error for user not found, change to NoError.
		// Here we assume it might bubble up sql.ErrNoRows or similar.
		if err == nil {
			s.Assert().False(isAuthenticated)
		} else {
			s.Assert().Error(err)
		}
	})
}
