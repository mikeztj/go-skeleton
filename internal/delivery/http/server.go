package http

import (
	"net/http"

	"github.com/vilbert/go-skeleton/pkg/grace"
)

// SkeletonHandler ...
type SkeletonHandler interface {
	SkeletonHandler(w http.ResponseWriter, r *http.Request)
}

// Server ...
type Server struct {
	server   *http.Server
	Skeleton SkeletonHandler
}

// Serve is serving HTTP gracefully on port x ...
func (s *Server) Serve(port string) error {
	return grace.Serve(port, s.Handler())
}
