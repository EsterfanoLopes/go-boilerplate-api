// Package api holds api routes setup and tests
package api

import (
	comment "go-boilerplate/api/comment/v1"

	"go-boilerplate/api/healthcheck"

	"go-boilerplate/common"
	"go-boilerplate/common/response"
	"net/http"
	"net/http/pprof"
	"strings"
	"sync/atomic"
	"time"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/handlers"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

// @title Go Boilerplate
// @version 1.0
// @description Go Boilerplate API
// @contact.name Esterfano Lopes
// @contact.url https://github.com/EsterfanoLopes
// @contact.email esterfano.lopes@gmail.com

var sentryHandler = sentryhttp.New(sentryhttp.Options{
	Repanic: true,
})

var corsHandler = cors.New(cors.Options{
	AllowedOrigins: []string{
		"*.vivareal.com.br",
		"*.zapimoveis.com.br",
		"*.grupozap.com",
	},
	AllowedHeaders: []string{
		"*",
	},
	AllowedMethods: []string{
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
	},
	MaxAge: 2592000, // 1 month
})

var apiReady = int32(0)

func apiIsReady() {
	atomic.StoreInt32(&apiReady, 1)
}

func isAPIReady() bool {
	return atomic.LoadInt32(&apiReady) == 1
}

// Setup configures api routes
func Setup() {
	if isAPIReady() {
		return
	}

	r := mux.NewRouter(mux.WithServiceName("go-boilerplate-mux"), mux.WithIgnoreRequest(func(r *http.Request) bool {
		return strings.HasPrefix(r.URL.Path, "/healthcheck")
	}))

	r.Handle("/healthcheck/status", handler{
		handler: healthcheck.SimpleHandler,
	}.build()).Methods(http.MethodGet)

	r.Handle("/healthcheck", handler{
		handler: healthcheck.CompleteHandler,
	}.build()).Methods(http.MethodGet)

	setupCommentRoutes(r)

	srv := &http.Server{
		ReadTimeout:  time.Duration(common.Config.GetInt64("httpServerReadTimeoutSeconds")) * time.Second,
		WriteTimeout: time.Duration(common.Config.GetInt64("httpServerWriteTimeoutSeconds")) * time.Second,
		Addr:         ":9000",
		Handler:      handlers.CompressHandler(r),
	}
	common.Logger.Info("server is ready at http://localhost:9000")

	apiIsReady()

	common.Logger.Fatal(srv.ListenAndServe())
}

func setupSwagger(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/swagger/index.html", http.StatusMovedPermanently)
	})

	r.HandleFunc("/docs/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger/swagger.json")
	}).Methods(http.MethodGet)

	r.PathPrefix("/docs/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger/swagger.json"),
	)).Methods(http.MethodGet)
}

func setupDebugRoutes(r *mux.Router) {
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)
}

func setupCommentRoutes(r *mux.Router) {
	r.Handle("/v1/comment", handler{
		handler: comment.CommentPostHandler,
	}.build()).Methods(http.MethodPost)

	r.Handle("/v1/comment/{id:[0-9]+}", handler{
		handler: comment.CommentPutHandler,
	}.build()).Methods(http.MethodPut)

	r.Handle("/v1/comment/{id:[0-9]+}", handler{
		handler: comment.CommentGetHandler,
	}.build()).Methods(http.MethodGet)

	r.Handle("/v1/comment", handler{
		handler: comment.CommentsGetHandler,
	}.build()).Methods(http.MethodGet)

	r.Handle("/v1/comment/{id:[0-9]+}", handler{
		handler: comment.CommentDeleteHandler,
	}.build()).Methods(http.MethodDelete)
}

type handler struct {
	cors    bool
	handler http.HandlerFunc
}

func (o handler) build() http.Handler {
	h := errorHandler(o.handler)
	if o.cors {
		h = corsHandler.Handler(h)
	}
	return h
}

func errorHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer handleUnexpectedError(w, r)
		sentryHandler.Handle(h).ServeHTTP(w, r)
	})
}

func handleUnexpectedError(w http.ResponseWriter, r *http.Request) {
	err := recover()
	if err != nil {
		casted, ok := err.(error)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.WriteServerError(w, casted, "unexpected error")
	}
}
