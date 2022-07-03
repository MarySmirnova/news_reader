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

func (m *Memdb) NewsAmount(filter string) (int, error) {
	return 1, nil
}

func (m *Memdb) GetNewsPage(filter string, page int, ipemsPerPage int) ([]*Post, error) {
	var posts []*Post

	for i := 0; i < 2; i++ {
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

func (m *Memdb) GetNewsByID(id int) (*Post, error) {
	return &Post{
		ID:      1,
		Title:   "Title",
		Content: "Content",
		PubTime: time.Now().Unix(),
		Link:    "Link",
	}, nil
}
