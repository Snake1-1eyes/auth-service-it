package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Snake1-1eyes/auth-service-it/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (int64, error)
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context, sessionID string, userID int64, duration time.Duration) error
}

type UseCase struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func NewUseCase(userRepo UserRepository, sessionRepo SessionRepository) *UseCase {
	return &UseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *UseCase) SignUp(ctx context.Context, login, password string) (int64, error) {
	// Check if user exists
	existingUser, err := uc.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return 0, err
	}
	if existingUser != nil {
		return 0, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Login:        login,
		PasswordHash: string(hashedPassword),
	}

	return uc.userRepo.CreateUser(ctx, user)
}

func (uc *UseCase) Login(ctx context.Context, login, password string) (string, error) {
	user, err := uc.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	sessionID := uuid.New().String()
	err = uc.sessionRepo.CreateSession(ctx, sessionID, user.ID, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}
