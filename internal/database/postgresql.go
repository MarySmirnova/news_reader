package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/MarySmirnova/news_reader/internal/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ctx context.Context = context.Background()

type Store struct {
	db *pgxpool.Pool
}

func NewPostgresDB(cfg config.Postgres) (*Store, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(ctx); err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

func (s *Store) GetPGXPool() *pgxpool.Pool {
	return s.db
}

func (s *Store) WriteNews(posts []*Post) error {
	query := `
	INSERT INTO news.posts (
		title,
		content,
		pubTime,
		link)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (link) DO NOTHING;`

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	for _, post := range posts {
		_, err = tx.Exec(ctx, query, post.Title, post.Content, post.PubTime, post.Link)
		if err != nil {
			return err
		}
	}

	tx.Commit(ctx)
	return nil
}

func (s *Store) GetLastNews(n int) ([]*Post, error) {
	query := `
	SELECT 
		id,
		title,
		content,
		pubTime,
		link
	FROM news.posts
	ORDER BY pubTime DESC
	LIMIT $1;`

	var posts []*Post

	rows, err := s.db.Query(ctx, query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post

		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return posts, nil
}
