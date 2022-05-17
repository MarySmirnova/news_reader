package rss

import (
	"testing"

	"github.com/MarySmirnova/news_reader/internal/config"
	"github.com/MarySmirnova/news_reader/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestNewsParser_readAllRSS(t *testing.T) {
	db := database.NewMemoryDB()
	links := []string{"https://habr.com/ru/rss/hub/go/all/?fl=ru", "https://habr.com/ru/rss/best/daily/?fl=ru"}

	p := NewNewsParser(config.RSS{
		Links:         links,
		RequestPeriod: 1,
	}, db)

	posts := p.readAllRSS()

	assert.True(t, len(posts) > 0)
}
