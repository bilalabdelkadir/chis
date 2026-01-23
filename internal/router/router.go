package router

import (
	"net/http"

	"github.com/bilalabdelkadir/chis/pkg/apperror"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	mux chi.Router
}

func NewRouter() *Router {
	return &Router{
		mux: chi.NewRouter(),
	}
}

func (r *Router) Get(pattern string, h apperror.HandlerFunc) {
	r.mux.Get(pattern, apperror.Wrap(h))
}

func (r *Router) Post(pattern string, h apperror.HandlerFunc) {
	r.mux.Post(pattern, apperror.Wrap(h))
}

func (r *Router) Put(pattern string, h apperror.HandlerFunc) {
	r.mux.Put(pattern, apperror.Wrap(h))
}

func (r *Router) Patch(pattern string, h apperror.HandlerFunc) {
	r.mux.Patch(pattern, apperror.Wrap(h))
}

func (r *Router) Delete(pattern string, h apperror.HandlerFunc) {
	r.mux.Delete(pattern, apperror.Wrap(h))
}

func (r *Router) Options(pattern string, h apperror.HandlerFunc) {
	r.mux.Options(pattern, apperror.Wrap(h))
}

func (r *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	r.mux.Use(middlewares...)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) Route(pattern string, fn func(r *Router)) {
	r.mux.Route(pattern, func(cr chi.Router) {
		fn(&Router{mux: cr})
	})
}
