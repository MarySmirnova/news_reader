package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/MarySmirnova/news_reader/internal/config"
	"github.com/MarySmirnova/news_reader/internal/database"

	log "github.com/sirupsen/logrus"
)

type storage interface {
	GetLastNews(int) ([]*database.Post, error)
}

type API struct {
	db         storage
	httpServer *http.Server
}

func New(cfg config.API, db storage) *API {
	a := &API{
		db: db,
	}

	handler := mux.NewRouter()
	handler.Name("get_some_last_news").Path("/news/{n}").Methods(http.MethodGet).HandlerFunc(a.PostsHandler)

	a.httpServer = &http.Server{
		Addr:         cfg.Listen,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		Handler:      handler,
	}

	return a
}

func (a *API) GetHTTPServer() *http.Server {
	return a.httpServer
}

func (a *API) PostsHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(news)
}

func (a *API) writeResponseError(w http.ResponseWriter, err error, code int) {
	log.WithError(err).Error("api error")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}
