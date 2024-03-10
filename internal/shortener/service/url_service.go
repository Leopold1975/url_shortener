package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leopold1975/url_shortener/internal/shortener/domain/urls"
	"github.com/Leopold1975/url_shortener/internal/shortener/repository"
	"github.com/Leopold1975/url_shortener/pkg/logger"
)

var ErrInvalidURL = errors.New("invalid URL")

type Repository interface {
	CreateURL(ctx context.Context, u urls.URL) (string, error)
	GetURL(ctx context.Context, shortURL string) (urls.URL, error)
	DeleteURL(ctx context.Context, shortURL string) error
	UpdateURL(ctx context.Context, u urls.URL) (urls.URL, error)
}

type ShortenerService struct {
	str   Repository
	cache Repository
	lg    logger.Logger
}

func New(repo Repository, cache Repository, l logger.Logger) ShortenerService {
	return ShortenerService{
		str:   repo,
		cache: cache,
		lg:    l,
	}
}

func (s ShortenerService) CreateShortURL(ctx context.Context, longURL string) (string, error) {
	if !urls.Validate(longURL) {
		return "", ErrInvalidURL
	}

	u, err := urls.PrepareURL(longURL)
	if err != nil {
		return "", err
	}

	if _, err := s.cache.CreateURL(ctx, u); err != nil {
		s.lg.Error("create URL cache error: %s", err.Error())
	}

	return s.str.CreateURL(ctx, u)
}

func (s ShortenerService) GetURL(ctx context.Context, shortURL string) (urls.URL, error) {
	u, err := s.cache.GetURL(ctx, shortURL)
	switch {
	case err == nil:
		return u, nil
	case errors.Is(err, repository.ErrNotFound):
		s.lg.Info("cache missed")
		if _, err := s.cache.UpdateURL(ctx, u); err != nil {
			s.lg.Error("update URL cache error: %s", err.Error())
		}
	default:
		s.lg.Error("get URL cache error: ", err.Error())
	}

	u, err = s.str.GetURL(ctx, shortURL)

	return u, err
}

func (s ShortenerService) GetURLWithInc(ctx context.Context, shortURL string) (urls.URL, error) {
	u, err := s.GetURL(ctx, shortURL)
	if err != nil {
		return urls.URL{}, err
	}

	u.Clicks++
	if _, err := s.cache.UpdateURL(ctx, u); err != nil {
		s.lg.Error("update URL cache error: %s", err.Error())
	}
	return s.str.UpdateURL(ctx, u)
}

func (s ShortenerService) DeleteURL(ctx context.Context, shortURL string) error {
	err := s.cache.DeleteURL(ctx, shortURL)

	err2 := s.str.DeleteURL(ctx, shortURL)

	switch {
	case err != nil && err2 != nil:
		return fmt.Errorf("cache error: %w    repo error: %w", err, err2)
	case err2 != nil:
		return err2
	case err != nil:
		s.lg.Error("get cache error: ", err.Error())
		return nil
	default:
		return nil
	}
}
