package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/princetheprogrammerbtw/nanoci/internal/config"
	"github.com/princetheprogrammerbtw/nanoci/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type AuthService struct {
	oauthConfig *oauth2.Config
	userRepo    domain.UserRepository
}

func NewAuthService(cfg *config.Config, userRepo domain.UserRepository) *AuthService {
	return &AuthService{
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GithubClientID,
			ClientSecret: cfg.GithubSecret,
			Endpoint:     github.Endpoint,
			Scopes:       []string{"user:email", "repo"},
		},
		userRepo: userRepo,
	}
}

func (s *AuthService) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state)
}

func (s *AuthService) HandleCallback(ctx context.Context, code string) (*domain.User, error) {
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth exchange failed: %w", err)
	}

	client := s.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get github user: %w", err)
	}
	defer resp.Body.Close()

	var ghUser GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return nil, fmt.Errorf("failed to decode github user: %w", err)
	}

	githubID := fmt.Sprintf("%d", ghUser.ID)
	user, err := s.userRepo.GetByGithubID(ctx, githubID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user = &domain.User{
			GithubID:  githubID,
			Username:  ghUser.Login,
			Email:     ghUser.Email,
			AvatarURL: ghUser.AvatarURL,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
