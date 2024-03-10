package repository

import (
	"context"
	"sync"

	"github.com/Leopold1975/url_shortener/internal/shortener/domain/urls"
)

type FakeRepo struct {
	r  map[string]urls.URL
	mu sync.RWMutex
}

func NewFakeRepo() FakeRepo {
	return FakeRepo{
		r: make(map[string]urls.URL, 100),
	}
}
func (f *FakeRepo) CreateURL(_ context.Context, u urls.URL) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.r[u.ShortURL]; ok {
		return "", ErrAleradyExists
	}
	f.r[u.ShortURL] = u

	return u.ShortURL, nil
}

func (f *FakeRepo) GetURL(_ context.Context, shortURL string) (urls.URL, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	u, ok := f.r[shortURL]
	if !ok {
		return urls.URL{}, ErrNotFound
	}

	return u, nil
}

func (f *FakeRepo) DeleteURL(_ context.Context, shortURL string) error {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if _, ok := f.r[shortURL]; !ok {
		return ErrNotFound
	}

	delete(f.r, shortURL)
	return nil
}

func (f *FakeRepo) UpdateURL(_ context.Context, u urls.URL) (urls.URL, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.r[u.ShortURL]; !ok {
		return urls.URL{}, ErrNotFound
	}

	f.r[u.ShortURL] = u

	return u, nil
}
