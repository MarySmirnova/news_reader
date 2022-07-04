package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

//PostsHandler waits for parameter n in the request path, returns the latest n news.
func (a *API) SomePostsHandler(w http.ResponseWriter, r *http.Request) {
	nn := mux.Vars(r)["n"]
	n, err := strconv.Atoi(nn)
	if err != nil {
		a.writeResponseError(w, err, http.StatusBadRequest)
		return
	}

	news, err := a.db.GetLastNews(n)
	if err != nil {
		a.writeResponseError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Code", strconv.Itoa(http.StatusOK))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(news)
}

//AllPostsHandler returns a page with news found by filter.
//Accepts "filter" and "page" parameters.
func (a *API) AllPostsHandler(w http.ResponseWriter, r *http.Request) {
	page, filter, err := a.getPageAndFilterParams(w, r)
	if err != nil {
		a.writeResponseError(w, err, http.StatusBadRequest)
	}

	itemsAmount, err := a.db.NewsAmount(filter)
	if err != nil {
		a.writeResponseError(w, err, http.StatusInternalServerError)
		return
	}

	posts, err := a.db.GetNewsPage(filter, page, itemsPerPage)
	if err != nil {
		a.writeResponseError(w, err, http.StatusInternalServerError)
		return
	}

	resp := ResponseNews{
		Posts: posts,
		Page: Page{
			TotalPages:   itemsAmount / itemsPerPage,
			NumberOfPage: page,
			ItemsPerPage: itemsPerPage,
		},
	}

	w.Header().Add("Code", strconv.Itoa(http.StatusOK))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

//PostHandler returns one piece of news by its id.
func (a *API) PostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		a.writeResponseError(w, err, http.StatusBadRequest)
		return
	}

	post, err := a.db.GetNewsByID(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			a.writeResponseError(w, err, http.StatusBadRequest)
			return
		}
		a.writeResponseError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Code", strconv.Itoa(http.StatusOK))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(post)
}
