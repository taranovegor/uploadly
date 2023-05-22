package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

type StaticHandler struct {
	rootPath string
}

func NewStaticHandler(
	rootPath string,
) StaticHandler {
	return StaticHandler{
		rootPath: rootPath,
	}
}

func (handler StaticHandler) get(writer http.ResponseWriter, request *http.Request) {
	rctx := chi.RouteContext(request.Context())
	pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
	fs := http.StripPrefix(pathPrefix, http.FileServer(http.Dir(handler.rootPath)))
	fs.ServeHTTP(writer, request)
}
