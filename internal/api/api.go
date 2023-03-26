package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lekht/news-gateway/config"
)

type API struct {
	r                *mux.Router
	NewsAddress      string
	CommentsAddress  string
	FormatterAddress string
}

func New(cfg *config.API) *API {
	api := &API{
		r:                mux.NewRouter(),
		NewsAddress:      cfg.NewsAddr,
		CommentsAddress:  cfg.CommentsAddr,
		FormatterAddress: cfg.FormatterAddr,
	}

	api.endpoints()

	return api
}

func (a *API) Router() *mux.Router {
	return a.r
}

func (a *API) endpoints() {
	a.r.Use(a.accessMiddleware, a.requestIdMiddlware, a.logRequestMiddlware)
	a.r.Name("news_list").Methods(http.MethodGet).Path("/news").HandlerFunc(a.newsListHandler)
	a.r.Name("full_new").Methods(http.MethodGet).Path("/news/{id}").HandlerFunc(a.fullNewsHandler)
	a.r.Name("add_comment").Methods(http.MethodPost).Path("/news/comment").HandlerFunc(a.addCommentHandler)
}
