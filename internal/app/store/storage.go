package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(dbUrl string) (*Storage, error) {

	conn, err := pgx.Connect(context.Background(), dbUrl)

	if err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	storage := &Storage{
		conn: conn,
	}

	err = storage.PrepareUrlTable()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) PrepareUrlTable() error {

	query := `
		DROP TABLE IF EXISTS url;
		CREATE TABLE IF NOT EXISTS url(
			id BIGSERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`

	_, err := s.conn.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to prepare url table in storage: %w", err)
	}

	return nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {

	query := `
		INSERT INTO url(alias, url)
		VALUES ($1,$2)
		RETURNING id
	`

	var id int64

	err := s.conn.QueryRow(context.Background(), query, alias, urlToSave).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to save url in storage: %w", err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {

	query := `
		SELECT url
		FROM url
		WHERE alias = $1
	`

	var url string

	err := s.conn.QueryRow(context.Background(), query, alias).Scan(&url)

	if err != nil {
		return "", fmt.Errorf("failed to get url by alias in storage: %w", err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {

	query := `
		DELETE url
		FROM url
		WHERE alias = $1
	`

	var url string

	err := s.conn.QueryRow(context.Background(), query, alias).Scan(&url)

	if err != nil {
		return fmt.Errorf("failed to delete url by alias in storage: %w", err)
	}

	return nil
}
