package api

import "github.com/MarySmirnova/news_reader/internal/database"

type ResponseNews struct {
	Page  Page
	Posts []*database.Post
}

type Page struct {
	TotalPages   int // общее количество страниц по запросу
	NumberOfPage int // номер страницы
	ItemsPerPage int // количество новостей на одной странице
}
