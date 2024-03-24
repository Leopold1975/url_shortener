package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Leopold1975/url_shortener/internal/pkg/config"
	"github.com/Leopold1975/url_shortener/internal/shortener/domain/urls"
	"github.com/Leopold1975/url_shortener/internal/shortener/repository"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // driver for migrations
	"github.com/pressly/goose/v3"
)

type ShortenerRepo struct {
	// db *pgx.ConnPool // TODO: run bench to know, which one is faster
	db *pgxpool.Pool
}

func NewPostgresRepo(ctx context.Context, cfg config.DB) (ShortenerRepo, error) {
	connString := "postgres://" + cfg.Username + ":" + cfg.Password + "@" +
		cfg.Addr + "/" + cfg.DB + "?" + "sslmode=" + cfg.SSLmode + "&pool_max_conns=" + cfg.MaxConns

	dbP := new(pgxpool.Pool)

	if err := connect(ctx, connString, &dbP); err != nil {
		return ShortenerRepo{}, err
	}

	if err := applyMigration(cfg); err != nil {
		return ShortenerRepo{}, err
	}

	return ShortenerRepo{
		db: dbP,
	}, nil
}

func (s ShortenerRepo) CreateURL(ctx context.Context, url urls.URL) (short string, err error) { //nolint:nonamedreturns
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("cannot begin transection error: %w", err)
	}

	defer func() {
		err = commitOrRollback(ctx, tx, err, "create")
	}()

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqB := psql.Insert("urls").
		Columns("uuid", "url_long", "url_short", "created_at", "clicks").
		Values(url.UUID, url.LongURL, url.ShortURL, url.CreatedAt, url.Clicks)

	query, args, err := sqB.ToSql()
	if err != nil {
		return "", fmt.Errorf("to sql error: %w", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return "", fmt.Errorf("exec error: %w", err)
	}

	return url.ShortURL, nil
}

func (s ShortenerRepo) GetURL(ctx context.Context, shortURL string) (urls.URL, error) {
	var url urls.URL

	errCh := make(chan error)
	done := make(chan struct{})

	go func() {
		defer close(errCh)
		defer close(done)

		psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
		sqB := psql.Select("uuid", "url_long", "url_short", "created_at", "clicks").Where(
			squirrel.Eq{"url_short": shortURL}).From("urls")

		query, args, err := sqB.ToSql()
		if err != nil {
			errCh <- fmt.Errorf("to sql error: %w", err)

			return
		}

		row := s.db.QueryRow(ctx, query, args...)

		err = row.Scan(&url.UUID, &url.LongURL, &url.ShortURL, &url.CreatedAt, &url.Clicks)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				errCh <- repository.ErrNotFound

				return
			}
			errCh <- fmt.Errorf("scan error: %w", err)
		}
	}()
	select {
	case <-ctx.Done():
		return urls.URL{}, fmt.Errorf("context error: %w", ctx.Err())
	case err := <-errCh:
		return urls.URL{}, err
	case <-done:
		return url, nil
	}
}

func (s ShortenerRepo) DeleteURL(ctx context.Context, shortURL string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cannot begin transection error: %w", err)
	}

	defer func() {
		err = commitOrRollback(ctx, tx, err, "delete")
	}()

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqB := psql.Delete("urls").Where(
		squirrel.Eq{"url_short": shortURL})

	query, args, err := sqB.ToSql()
	if err != nil {
		return fmt.Errorf("to sql error: %w", err)
	}

	ct, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (s ShortenerRepo) UpdateURL(ctx context.Context, u urls.URL) (urls.URL, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return urls.URL{}, fmt.Errorf("cannot begin transaction error: %w", err)
	}

	defer func() {
		err = commitOrRollback(ctx, tx, err, "update")
	}()

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqB := psql.Update("urls").SetMap(squirrel.Eq{
		"url_long":  u.LongURL,
		"url_short": u.ShortURL,
		"clicks":    u.Clicks,
	}).Where(
		squirrel.Eq{"uuid": u.UUID})

	query, args, err := sqB.ToSql()
	if err != nil {
		return urls.URL{}, fmt.Errorf("to sql error: %w", err)
	}

	ct, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return urls.URL{}, fmt.Errorf("exec error: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return urls.URL{}, repository.ErrNotFound
	}

	return u, nil
}

func connect(ctx context.Context, connString string, dbP **pgxpool.Pool) error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)

		dbc, err := pgxpool.New(ctx, connString)
		if err != nil {
			errCh <- fmt.Errorf("cannot create db pool error: %w", err)

			return
		}

		defaultDelay := time.Second

		for {
			if err := dbc.Ping(ctx); err != nil {
				time.Sleep(defaultDelay)
				defaultDelay += time.Second

				if defaultDelay > time.Second*10 {
					errCh <- fmt.Errorf("cannot ping db error: %w", err)

					return
				}

				continue
			}

			break
		}

		*dbP = dbc
	}()
	select {
	case <-ctx.Done():
		return fmt.Errorf("context error: %w", ctx.Err())
	case err := <-errCh:
		return err
	}
}

func applyMigration(cfg config.DB) error {
	migrationsDir := "./migrations"
	defaultVersion := 0

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose set dialect error: %w", err)
	}

	connString := "postgres://" + cfg.Username + ":" + cfg.Password + "@" +
		cfg.Addr + "/" + cfg.DB

	dbM, err := goose.OpenDBWithDriver("pgx", connString)
	if err != nil {
		return fmt.Errorf("goose open pgx db error: %w", err)
	}
	defer dbM.Close()

	if cfg.Reload {
		if err := goose.DownTo(dbM, migrationsDir, int64(defaultVersion)); err != nil {
			return fmt.Errorf("goose down error: %w", err)
		}
	}

	if err := goose.UpTo(dbM, migrationsDir, int64(cfg.Version)); err != nil {
		return fmt.Errorf("goose up error: %w", err)
	}

	return nil
}

func commitOrRollback(ctx context.Context, tx pgx.Tx, err error, where string) error {
	if err == nil {
		if errT := tx.Commit(ctx); errT != nil {
			err = fmt.Errorf("commit error: %w", errT)
		}
	} else {
		if errT := tx.Rollback(ctx); errT != nil {
			err = fmt.Errorf("%s error: %w rollback error: %w", where, err, errT)
		} else {
			err = fmt.Errorf("%s error: %w", where, err)
		}
	}

	return err
}
