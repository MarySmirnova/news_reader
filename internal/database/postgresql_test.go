package database

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/MarySmirnova/news_reader/internal/config"
	"github.com/caarlos0/env/v6"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var cfg config.Postgres

func testPGDB(t *testing.T) (*Store, func()) {
	godotenv.Load("../../.env")
	err := env.Parse(&cfg)
	assert.Nil(t, err)

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.TestDatabase)

	db, err := pgxpool.Connect(ctx, connString)
	assert.Nil(t, err)

	schemaName := "news"

	_, err = db.Exec(ctx, "CREATE SCHEMA "+schemaName)
	assert.Nil(t, err)

	createTableQuery := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s.posts (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		pubTime BIGINT NOT NULL CHECK (pubTime > 0),
		link TEXT NOT NULL UNIQUE);`, schemaName)

	_, err = db.Exec(ctx, createTableQuery)
	assert.Nil(t, err)

	return &Store{db: db}, func() {
		_, err := db.Exec(ctx, "DROP SCHEMA "+schemaName+" CASCADE")
		assert.Nil(t, err)
	}
}

func generateSomePosts(n int) []*Post {
	var posts []*Post
	for i := 0; i < n; i++ {
		post := Post{
			Title:   "Title " + strconv.Itoa(i),
			Content: "Content " + strconv.Itoa(i),
			PubTime: time.Now().Unix(),
			Link:    "Link " + strconv.Itoa(i),
		}
		posts = append(posts, &post)
	}

	return posts
}

func TestStore_WriteNews(t *testing.T) {
	db, cleanup := testPGDB(t)
	defer cleanup()

	var n = 10
	posts := generateSomePosts(n)

	err := db.WriteNews(posts)
	assert.Nil(t, err)
}

func TestStore_GetLastNews(t *testing.T) {
	db, cleanup := testPGDB(t)
	defer cleanup()

	posts := generateSomePosts(20)
	err := db.WriteNews(posts)
	assert.Nil(t, err)

	n := 10
	lastPosts, err := db.GetLastNews(n)
	assert.Nil(t, err)

	assert.Equal(t, n, len(lastPosts))
}
