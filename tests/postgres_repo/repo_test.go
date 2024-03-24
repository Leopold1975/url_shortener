package postgres_test

import (
	"context"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/Leopold1975/url_shortener/internal/pkg/config"
	"github.com/Leopold1975/url_shortener/internal/shortener/domain/urls"
	"github.com/Leopold1975/url_shortener/internal/shortener/repository"
	"github.com/Leopold1975/url_shortener/internal/shortener/repository/postgres"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"
)

type ShortenerRepoTestSuite struct {
	suite.Suite
	repo postgres.ShortenerRepo
}

func (s *ShortenerRepoTestSuite) SetupSuite() {
	cmd := exec.Command("docker", "compose", "-f", "./test_compose.yaml", "up", "--build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	if err := cmd.Start(); err != nil {
		s.T().Fatalf("cannot start docker compose error: %v", err)
	}

	cfg := config.DB{
		Addr:     "127.0.0.1:3333",
		Username: "postgres",
		Password: "postgres",
		DB:       "urls",
		SSLmode:  "disable",
		MaxConns: "10",
		Reload:   true,
		Version:  1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	repo, err := postgres.NewPostgresRepo(ctx, cfg)
	s.Require().NoError(err)

	s.repo = repo
}

func (s *ShortenerRepoTestSuite) TearDownSuite() {
	cmd := exec.Command("docker", "compose", "-f", "./test_compose.yaml", "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	if err := cmd.Run(); err != nil {
		s.T().Fatalf("cannot start docker compose error: %v", err)
	}
}

func (s *ShortenerRepoTestSuite) TestBasic() {
	//s.T().Skip()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	testURL := "https://test132.com"

	u, err := urls.PrepareURL(testURL)
	s.Require().NoError(err, "expected %v	actual %v", nil, err)

	short, err := s.repo.CreateURL(ctx, u)
	s.Require().Equal(u.ShortURL, short, "expected %v	actual %v", u.ShortURL, short)
	s.Require().NoError(err, "expected %v	actual %v", nil, err)

	url, err := s.repo.GetURL(ctx, short)
	s.Require().Equal(testURL, url.LongURL, "expected %v	actual %v", testURL, u.LongURL)
	s.Require().Equal(u.ShortURL, url.ShortURL, "expected %v	actual %v", u.ShortURL, url.ShortURL)
	s.Require().Equal(u.Clicks, url.Clicks, "expected %v	actual %v", u.Clicks, url.Clicks)
	s.Require().NoError(err, "expected %v	actual %v", nil, err)

	u.Clicks++
	url, err = s.repo.UpdateURL(ctx, u)
	s.Require().Equal(testURL, url.LongURL, "expected %v	actual %v", testURL, u.LongURL)
	s.Require().Equal(u.ShortURL, url.ShortURL, "expected %v	actual %v", u.ShortURL, url.ShortURL)
	s.Require().Equal(u.Clicks, url.Clicks, "expected %v	actual %v", u.Clicks, url.Clicks)
	s.Require().NoError(err, "expected %v	actual %v", nil, err)

	err = s.repo.DeleteURL(ctx, url.ShortURL)
	s.Require().NoError(err, "expected %v	actual %v", nil, err)

	url, err = s.repo.GetURL(ctx, url.ShortURL)
	s.Require().ErrorIs(err, repository.ErrNotFound, "expected %v	actual %v", repository.ErrNotFound, err)
	s.Require().Equal("", url.UUID, "expected %v	actual %v", "", url.UUID)
	s.Require().Equal("", url.ShortURL, "expected %v	actual %v", "", url.ShortURL)
	s.Require().Equal("", url.LongURL, "expected %v	actual %v", "", url.LongURL)
	s.Require().Equal(int64(0), url.Clicks, "expected %v	actual %v", 0, url.Clicks)
	s.Require().Equal(time.Time{}, url.CreatedAt, "expected %v	actual %v", time.Time{}, url.CreatedAt.UnixNano())

	u.Clicks++
	url, err = s.repo.UpdateURL(ctx, u)
	s.Require().ErrorIs(err, repository.ErrNotFound, "expected %v	actual %v", repository.ErrNotFound, err)
	s.Require().Equal("", url.UUID, "expected %v	actual %v", "", url.UUID)
	s.Require().Equal("", url.ShortURL, "expected %v	actual %v", "", url.ShortURL)
	s.Require().Equal("", url.LongURL, "expected %v	actual %v", "", url.LongURL)
	s.Require().Equal(int64(0), url.Clicks, "expected %v	actual %v", 0, url.Clicks)
	s.Require().Equal(time.Time{}, url.CreatedAt, "expected %v	actual %v", time.Time{}, url.CreatedAt.UnixNano())

	err = s.repo.DeleteURL(ctx, u.ShortURL)
	s.Require().ErrorIs(err, repository.ErrNotFound, "expected %v	actual %v", repository.ErrNotFound, err)

	err = s.repo.DeleteURL(ctx, testURL)
	s.Require().ErrorIs(err, repository.ErrNotFound, "expected %v	actual %v", repository.ErrNotFound, err)
}

func (s *ShortenerRepoTestSuite) TestRepoWithLoad() {
	gofakeit.Seed(123)
	m := make(map[string]urls.URL, 10000)
	var u string

	for i := 0; i < 10000; i++ {
		u = gofakeit.URL()

		if _, ok := m[u]; ok {
			continue
		}

		ur, err := urls.PrepareURL(u)
		s.Require().NoError(err, "expected %v	actual %v", nil, err)
		m[u] = ur
	}

	var wg sync.WaitGroup
	urlsC := make(chan urls.URL, 1000)

	go func() {
		for _, v := range m {
			urlsC <- v
		}
		close(urlsC)
	}()
	time.Sleep(time.Second)
	semaphore := make(chan struct{}, 100)

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			for u := range urlsC {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

				short, err := s.repo.CreateURL(ctx, u)
				s.Require().NoError(err, "expected %v	actual %v", nil, err)
				s.Require().Equal(u.ShortURL, short, "expected %v	actual %v", u.ShortURL, short)

				url, err := s.repo.GetURL(ctx, short)
				s.Require().Equal(u.ShortURL, url.ShortURL, "expected %v	actual %v", u.ShortURL, url.ShortURL)
				s.Require().Equal(u.Clicks, url.Clicks, "expected %v	actual %v", u.Clicks, url.Clicks)
				s.Require().NoError(err, "expected %v	actual %v", nil, err)

				u.Clicks++
				url, err = s.repo.UpdateURL(ctx, u)
				s.Require().Equal(u.ShortURL, url.ShortURL, "expected %v	actual %v", u.ShortURL, url.ShortURL)
				s.Require().Equal(u.Clicks, url.Clicks, "expected %v	actual %v", u.Clicks, url.Clicks)
				s.Require().NoError(err, "expected %v	actual %v", nil, err)

				url, err = s.repo.GetURL(ctx, short)
				s.Require().Equal(u.ShortURL, url.ShortURL, "expected %v	actual %v", u.ShortURL, url.ShortURL)
				s.Require().Equal(u.Clicks, url.Clicks, "expected %v	actual %v", u.Clicks, url.Clicks)
				s.Require().NoError(err, "expected %v	actual %v", nil, err)
				cancel()
			}
		}()
	}

	wg.Wait()
}

func TestShortenerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(ShortenerRepoTestSuite))
}
