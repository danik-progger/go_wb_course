package connections

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	port    string
	mux     *chi.Mux
	handler http.Handler
}

func NewRouter(port string) *Router {
	mux := chi.NewRouter()
	return &Router{
		port:    port,
		mux:     mux,
		handler: mux,
	}
}

func (r *Router) AddMiddleware(m func(next http.Handler) http.Handler) {
	r.handler = m(r.mux)
}

func (r *Router) Get(endpoint string, f func(w http.ResponseWriter, r *http.Request)) {
	r.mux.Get(endpoint, f)
}
func (r *Router) Post(endpoint string, f func(w http.ResponseWriter, r *http.Request)) {
	r.mux.Post(endpoint, f)
}

func (ro *Router) GetUrlParam(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}

func (r *Router) Run() error {
	return http.ListenAndServe(r.port, r.handler)
}
