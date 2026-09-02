package user

import "github.com/go-chi/chi/v5"

type Router struct {
	*Handler
}

func NewUserRouter(handler *Handler) *Router {
	return &Router{
		handler,
	}
}

func (route *Router) Register(r chi.Router) {

	r.Route("/users", func(r chi.Router) {
		r.Get("/", route.GetUsers)
		r.Post("/create", route.CreateUser)
	})
}
