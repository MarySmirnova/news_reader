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

//NewPostgresDB creates a new instance Store for PostgresDB.
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

//GetPGXPool returns the pgxpool.Pool.
func (s *Store) GetPGXPool() *pgxpool.Pool {
	return s.db
}

//WriteNews adds posts to the database, checking links for uniqueness.
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

//GetLastNews returns the latest n news, sorted by publication date.
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

//NewsAmount returns the number of news by filter.
func (s *Store) NewsAmount(filter string) (int, error) {
	query := `
	SELECT count(*)
	FROM news.posts
	WHERE title ILIKE '%$1%';`

	var amount int

	row := s.db.QueryRow(ctx, query, filter)
	err := row.Scan(&amount)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

//GetNews returns the specified page with news by filter
func (s *Store) GetNewsPage(filter string, page int, ipemsPerPage int) ([]*Post, error) {
	query := `
	SELECT 
		id,
		title,
		content,
		pubTime,
		link
	FROM news.posts
	WHERE title ILIKE '%$1%'
	ORDER BY pubTime DESC
	LIMIT $2
	OFFSET $3;`

	offset := (page - 1) * ipemsPerPage

	var posts []*Post

	rows, err := s.db.Query(ctx, query, filter, ipemsPerPage, offset)
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

//GetNewsByID returns one post by its id
func (s *Store) GetNewsByID(id int) (*Post, error) {
	query := `
	SELECT 
		id,
		title,
		content,
		pubTime,
		link
	FROM news.posts
	WHERE id = $1;`

	var post *Post

	row := s.db.QueryRow(ctx, query, id)
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
	if err != nil {
		return nil, err
	}

	return post, nil
}
