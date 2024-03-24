package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Leopold1975/url_shortener/internal/pkg/config"
	"github.com/Leopold1975/url_shortener/internal/shortener/domain/urls"
	"github.com/Leopold1975/url_shortener/internal/shortener/repository"
	"github.com/redis/go-redis/v9"
)

type ShortenerRepo struct {
	rdb     *redis.Client
	expTime time.Duration
}

func New(ctx context.Context, cfg config.RedisDB) (ShortenerRepo, error) {
	rdb := redis.NewClient(&redis.Options{ //nolint:exhaustruct
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := connect(ctx, rdb); err != nil {
		return ShortenerRepo{}, err
	}

	return ShortenerRepo{
		rdb:     rdb,
		expTime: cfg.ExpTime,
	}, nil
}

func connect(ctx context.Context, rdb *redis.Client) error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)

		defaultDelay := time.Second

		for {
			if err := rdb.Ping(ctx).Err(); err != nil {
				time.Sleep(defaultDelay)
				defaultDelay += time.Second

				if defaultDelay > time.Second*10 {
					errCh <- fmt.Errorf("cannot ping redis db error: %w", err)

					return
				}

				continue
			}

			break
		}
	}()
	select {
	case <-ctx.Done():
		return fmt.Errorf("context error: %w", ctx.Err())
	case err := <-errCh:
		return err
	}
}

func (s ShortenerRepo) CreateURL(ctx context.Context, u urls.URL) (string, error) {
	exists, err := s.rdb.Exists(ctx, u.ShortURL).Result()
	if err != nil {
		return "", fmt.Errorf("exists operation error: %w", err)
	}

	if exists > 0 {
		return "", repository.ErrAleradyExists
	}

	dataURL, err := json.Marshal(u)
	if err != nil {
		return "", fmt.Errorf("url marshal error: %w", err)
	}

	err = s.rdb.Set(ctx, u.ShortURL, dataURL, s.expTime).Err()
	if err != nil {
		return "", fmt.Errorf("create url set error: %w", err)
	}

	return u.ShortURL, nil
}

func (s ShortenerRepo) GetURL(ctx context.Context, shortURL string) (urls.URL, error) {
	var u urls.URL

	dataURL, err := s.rdb.Get(ctx, shortURL).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return u, repository.ErrNotFound
		}

		return u, fmt.Errorf("get url error %w", err)
	}

	err = json.Unmarshal([]byte(dataURL), &u)
	if err != nil {
		return u, fmt.Errorf("unmarshal url error %w", err)
	}

	return u, nil
}

func (s ShortenerRepo) DeleteURL(ctx context.Context, shortURL string) error {
	affected, err := s.rdb.Del(ctx, shortURL).Result()
	if err != nil {
		return fmt.Errorf("delete url error %w", err)
	}

	if affected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (s ShortenerRepo) UpdateURL(ctx context.Context, u urls.URL) (urls.URL, error) {
	dataURL, err := json.Marshal(u)
	if err != nil {
		return urls.URL{}, fmt.Errorf("url marshal error: %w", err)
	}

	exists, err := s.rdb.Exists(ctx, u.ShortURL).Result()
	if err != nil {
		return urls.URL{}, fmt.Errorf("exists operation error: %w", err)
	}

	if exists == 0 {
		return urls.URL{}, repository.ErrNotFound
	}

	err = s.rdb.Set(ctx, u.ShortURL, dataURL, s.expTime).Err()
	if err != nil {
		return urls.URL{}, fmt.Errorf("create url set error: %w", err)
	}

	return u, nil
}
