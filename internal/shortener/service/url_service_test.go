package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Leopold1975/url_shortener/internal/pkg/config"
	"github.com/Leopold1975/url_shortener/internal/shortener/domain/urls"
	"github.com/Leopold1975/url_shortener/internal/shortener/repository"
	"github.com/Leopold1975/url_shortener/internal/shortener/service"
	"github.com/Leopold1975/url_shortener/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var testURL = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

func TestBasicWithCache(t *testing.T) {
	r1 := repository.NewFakeRepo()
	r2 := repository.NewFakeRepo()
	l, err := logger.New(config.Logger{Level: "info"})
	assert.NoError(t, err, "expected %v		actual: %v", nil, err)
	ctx := context.Background()

	s := service.New(&r1, &r2, l)

	sh, err := s.CreateShortURL(ctx, testURL)
	assert.NoError(t, err, "expected %v		actual: %v", nil, err)

	u, err := s.GetURL(ctx, sh)
	assert.NoError(t, err, "expected %v		actual: %v", nil, err)
	assert.Equal(t, testURL, u.LongURL, "expected %v		actual: %v", testURL, u.LongURL)
	assert.Equal(t, sh, u.ShortURL, "expected %v		actual: %v", sh, u.ShortURL)

	u, err = s.GetURLWithInc(ctx, sh)
	assert.NoError(t, err, "expected %v		actual: %v", nil, err)
	assert.Equal(t, testURL, u.LongURL, "expected %v		actual: %v", testURL, u.LongURL)
	assert.Equal(t, sh, u.ShortURL, "expected %v		actual: %v", sh, u.ShortURL)
	assert.Equal(t, int64(1), u.Clicks)

	u, err = s.GetURL(ctx, sh)
	assert.NoError(t, err, "expected %v		actual: %v", nil, err)
	assert.Equal(t, testURL, u.LongURL, "expected %v		actual: %v", testURL, u.LongURL)
	assert.Equal(t, sh, u.ShortURL, "expected %v		actual: %v", sh, u.ShortURL)
	assert.Equal(t, int64(1), u.Clicks)

	err = s.DeleteURL(ctx, u.ShortURL)
	assert.NoError(t, err, "expected %v		actual: %v", nil, err)

	u, err = s.GetURL(ctx, sh)
	assert.ErrorIs(t, err, repository.ErrNotFound, "expected %v		actual: %v", repository.ErrNotFound, err)
	assert.Equal(t, urls.URL{}, u)

	u, err = s.GetURLWithInc(ctx, sh)
	assert.ErrorIs(t, err, repository.ErrNotFound, "expected %v		actual: %v", repository.ErrNotFound, err)
	assert.Equal(t, urls.URL{}, u)
}

func TestBasicWithoutCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := repository.NewMockRepository(ctrl)

	r1 := repository.NewFakeRepo()
	l, err := logger.New(config.Logger{Level: "info"})
	assert.NoError(t, err, "expected %v		actual: %v", nil, err)
	ctx := context.Background()
	mockErr := errors.New("mock err")

	s := service.New(&r1, m, l)

	m.EXPECT().
		CreateURL(ctx, gomock.AssignableToTypeOf(urls.URL{})).
		Return("", mockErr).
		Times(1)

	sh, err := s.CreateShortURL(ctx, testURL)
	assert.NoError(t, err, "expected %v		actual: %v", nil, err)

	m.EXPECT().
		GetURL(ctx, sh).
		Return(urls.URL{}, mockErr).
		Times(1)

	u, err := s.GetURL(ctx, sh)
	assert.NoError(t, err, "expected %v		actual: %v", nil, err)
	assert.Equal(t, testURL, u.LongURL, "expected %v		actual: %v", testURL, u.LongURL)
	assert.Equal(t, sh, u.ShortURL, "expected %v		actual: %v", sh, u.ShortURL)

	m.EXPECT().
		GetURL(ctx, sh).
		Return(urls.URL{}, mockErr).
		Times(1)

	m.EXPECT().
		UpdateURL(ctx, gomock.AssignableToTypeOf(urls.URL{})).
		Return(urls.URL{}, mockErr).
		Times(1)

	u, err = s.GetURLWithInc(ctx, sh)
	assert.NoError(t, err, "expected %v		actual: %v", nil, err)
	assert.Equal(t, testURL, u.LongURL, "expected %v		actual: %v", testURL, u.LongURL)
	assert.Equal(t, sh, u.ShortURL, "expected %v		actual: %v", sh, u.ShortURL)
	assert.Equal(t, int64(1), u.Clicks)

	m.EXPECT().
		GetURL(ctx, sh).
		Return(urls.URL{}, mockErr).
		Times(1)

	u, err = s.GetURL(ctx, sh)
	assert.NoError(t, err, "expected %v		actual: %v", nil, err)
	assert.Equal(t, testURL, u.LongURL, "expected %v		actual: %v", testURL, u.LongURL)
	assert.Equal(t, sh, u.ShortURL, "expected %v		actual: %v", sh, u.ShortURL)
	assert.Equal(t, int64(1), u.Clicks)

	m.EXPECT().
		DeleteURL(ctx, sh).
		Return(mockErr).
		Times(1)

	err = s.DeleteURL(ctx, u.ShortURL)
	assert.NoError(t, err, "expected %v		actual: %v", nil, err)

	m.EXPECT().
		GetURL(ctx, sh).
		Return(urls.URL{}, mockErr).
		Times(1)

	u, err = s.GetURL(ctx, sh)
	assert.ErrorIs(t, err, repository.ErrNotFound, "expected %v		actual: %v", repository.ErrNotFound, err)
	assert.Equal(t, urls.URL{}, u)

	m.EXPECT().
		GetURL(ctx, sh).
		Return(urls.URL{}, mockErr).
		Times(1)

	u, err = s.GetURLWithInc(ctx, sh)
	assert.ErrorIs(t, err, repository.ErrNotFound, "expected %v		actual: %v", repository.ErrNotFound, err)
	assert.Equal(t, urls.URL{}, u)
}
