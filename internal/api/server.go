package api

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/MarySmirnova/news_reader/internal/config"
	"github.com/MarySmirnova/news_reader/internal/database"

	log "github.com/sirupsen/logrus"
)

type ContextKey string

const ContextReqIDKey ContextKey = "request_id"

const itemsPerPage = 15

type storage interface {
	GetLastNews(n int) ([]*database.Post, error)
	NewsAmount(filter string) (int, error)
	GetNewsPage(filter string, page int, ipemsPerPage int) ([]*database.Post, error)
	GetNewsByID(id int) (*database.Post, error)
}

type API struct {
	db         storage
	httpServer *http.Server
}

//New creates a new instance API.
func New(cfg config.API, db storage) *API {
	a := &API{
		db: db,
	}

	handler := mux.NewRouter()
	handler.Use(a.reqIDMiddleware)
	handler.Name("get_some_last_news").Path("/news/{n}").Methods(http.MethodGet).HandlerFunc(a.SomePostsHandler)
	handler.Name("get_all_news").Path("/news").Methods(http.MethodGet).HandlerFunc(a.AllPostsHandler)
	handler.Name("get_news_by_id").Path("/news/full/{id:[0-9]+}").Methods(http.MethodGet).HandlerFunc(a.PostHandler)

	handler.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))

	a.httpServer = &http.Server{
		Addr:         cfg.Listen,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		Handler:      handler,
	}

	return a
}

//GetHTTPServer returns the http.Server.
func (a *API) GetHTTPServer() *http.Server {
	return a.httpServer
}

func (a *API) reqIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqID int
		reqIDString := r.FormValue("request_id")

		if reqIDString == "" {
			reqID = a.generateReqID()
		}

		if reqIDString != "" {
			id, err := strconv.Atoi(reqIDString)
			if err != nil {
				a.writeResponseError(w, err, http.StatusBadRequest)
				return
			}
			reqID = id
		}

		ctx := context.WithValue(r.Context(), ContextReqIDKey, reqID)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			log.WithFields(log.Fields{
				"time":       time.Now().String(),
				"request_ip": strings.Split(r.RemoteAddr, ":")[0],
				"code":       w.Header().Get(),
				"request_id": r.Context().Value(ContextReqIDKey),
			}).Info("news reader response")
		}()

		next.ServeHTTP(w, r)
	})
}

func (a *API) generateReqID() int {
	max := 999999999999
	min := 100000

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func (a *API) getPageAndFilterParams(w http.ResponseWriter, r *http.Request) (int, string, error) {
	var page int
	filter := r.FormValue("filter")
	pageString := r.FormValue("page")
	if pageString == "" {
		page = 1
	}
	if pageString != "" {
		p, err := strconv.Atoi(pageString)
		if err != nil {
			return 0, "", err
		}
		page = p
	}

	return page, filter, nil
}

func (a *API) writeResponseError(w http.ResponseWriter, err error, code int) {
	log.WithError(err).Error("api error")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}
