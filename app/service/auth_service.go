package service

import (
	"fmt"
	"gemini-poc/app/adapter"
	"gemini-poc/app/dto"
	"gemini-poc/utils/config"
	"sync"

	"go.uber.org/zap"
)

type AuthService struct {
	da *adapter.DestinationAdapter

	token Token

	conf *config.AuthConfig

	log *zap.Logger
}

// Token struct with mutex
type Token struct {
	dto.AuthResponse
	mu sync.Mutex
}

func NewAuthService(
	da *adapter.DestinationAdapter,

	conf *config.AuthConfig,

	log *zap.Logger,
) *AuthService {
	a := AuthService{
		da: da,

		token: Token{},

		conf: conf,
		log:  log,
	}

	err := a.FetchServiceToken()
	if err != nil {
		log.Panic("Error getting service token", zap.Error(err))
	}

	return &a
}

func (a *AuthService) FetchServiceToken() error {
	a.token.mu.Lock()
	defer a.token.mu.Unlock()

	res, err := a.da.Login(dto.NewAuthRequest(a.conf.Username, a.conf.Password))
	if err != nil {
		a.log.Error("Error getting service token", zap.Error(err))
		return err
	}

	a.token.AuthResponse = *res

	return nil
}

func (a *AuthService) RefreshServiceToken() error {
	a.token.mu.Lock()
	defer a.token.mu.Unlock()

	res, err := a.da.RefreshToken(dto.NewRefreshTokenRequest(a.token.RefreshToken))
	if err != nil {
		a.log.Error("Error refreshing service token", zap.Error(err))
		return err
	}

	a.token.AuthResponse = *res

	return nil
}

func (a *AuthService) getServiceToken() *dto.AuthResponse {
	a.token.mu.Lock()
	defer a.token.mu.Unlock()

	return &a.token.AuthResponse
}

func (a *AuthService) GetAuthorizationHeader() string {
	return fmt.Sprintf("Bearer %s", a.getServiceToken().AccessToken)
}
