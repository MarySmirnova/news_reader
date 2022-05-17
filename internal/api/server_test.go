package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MarySmirnova/news_reader/internal/config"
	"github.com/MarySmirnova/news_reader/internal/database"
	"github.com/stretchr/testify/assert"
)

func testAPI(t *testing.T) *API {
	return New(config.API{
		Listen:       ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}, database.NewMemoryDB())
}

func execRequest(req *http.Request, s *http.Server) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	s.Handler.ServeHTTP(resp, req)

	return resp
}

func TestAPI_PostsHandler_InvalidParameter(t *testing.T) {
	api := testAPI(t)

	req, _ := http.NewRequest(http.MethodGet, "/news/{n}", nil)
	resp := execRequest(req, api.httpServer)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestAPI_PostsHandler_GoodWay(t *testing.T) {
	api := testAPI(t)
	n := 10

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/news/%d", n), nil)
	resp := execRequest(req, api.httpServer)
	assert.Equal(t, http.StatusOK, resp.Code)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var posts []database.Post
	err = json.Unmarshal(body, &posts)
	assert.Nil(t, err)

	assert.Equal(t, n, len(posts))
}
