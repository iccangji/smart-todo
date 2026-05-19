package auth

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(
		ctx context.Context,
		req RegisterRequest,
	) (*User, error)
	Login(
		ctx context.Context,
		req LoginRequest,
	) (string, string, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	return string(bytes), err
}

func comparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

func (s *service) Register(
	ctx context.Context,
	req RegisterRequest,
) (*User, error) {

	hash, err := hashPassword(req.Password)

	if err != nil {
		return nil, err
	}

	user := &User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	err = s.repository.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Login(
	ctx context.Context,
	req LoginRequest,
) (string, string, error) {

	user, err := s.repository.FindByEmail(
		ctx,
		req.Email,
	)

	if err != nil {
		return "", "", err
	}

	err = comparePassword(
		user.PasswordHash,
		req.Password,
	)

	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := GenerateAccessToken(*user)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(*user)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
