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
	defaultLen := 100

	return FakeRepo{
		r:  make(map[string]urls.URL, defaultLen),
		mu: sync.RWMutex{},
	}
}

func (f *FakeRepo) CreateURL(_ context.Context, url urls.URL) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.r[url.ShortURL]; ok {
		return "", ErrAleradyExists
	}

	f.r[url.ShortURL] = url

	return url.ShortURL, nil
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

func (f *FakeRepo) UpdateURL(_ context.Context, url urls.URL) (urls.URL, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.r[url.ShortURL]; !ok {
		return urls.URL{}, ErrNotFound
	}

	f.r[url.ShortURL] = url

	return url, nil
}
