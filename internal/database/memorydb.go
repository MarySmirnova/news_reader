package database

import (
	"strconv"
	"time"
)

type Memdb struct{}

func NewMemoryDB() *Memdb {
	return &Memdb{}
}

func (m *Memdb) WriteNews(posts []*Post) error {
	return nil
}

func (m *Memdb) GetLastNews(n int) ([]*Post, error) {
	var posts []*Post

	for i := 0; i < n; i++ {
		post := Post{
			ID:      i,
			Title:   "Title " + strconv.Itoa(i),
			Content: "Content " + strconv.Itoa(i),
			PubTime: time.Now().Unix(),
			Link:    "Link " + strconv.Itoa(i),
		}
		posts = append(posts, &post)
	}

	return posts, nil
}
