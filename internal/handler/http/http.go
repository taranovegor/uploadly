package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/taranovegor/uploadly/internal/config"
	"github.com/taranovegor/uploadly/pkg/dto"
	"log"
	"net/http"
)

func Init(
	router *chi.Mux,
	fileHandler FileHandler,
	staticHandler StaticHandler,
) {
	router.Group(func(group chi.Router) {
		group.Use(func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				handler.ServeHTTP(w, r)
			})
		})

		group.Post("/{context:\\w{1,255}}", fileHandler.save)
		group.Route("/{id:[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}}", func(route chi.Router) {
			route.Get("/", fileHandler.get)
			route.Delete("/", fileHandler.delete)
		})
	})

	router.Get(fmt.Sprintf("/%s/*", config.StaticHttpPath), staticHandler.get)
}

func response(writer http.ResponseWriter, code int, body interface{}) {
	marshal, err := json.Marshal(body)
	if err != nil {
		errno(writer, http.StatusInternalServerError, err)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	if _, err := writer.Write(marshal); err != nil {
		log.Println(err.Error())
	}
}

func errno(writer http.ResponseWriter, code int, err error) {
	message := http.StatusText(code)
	if config.IsDebug() && nil != err || len(err.Error()) == 0 {
		message = err.Error()
	}

	response(writer, code, dto.NewError(code, message))
}
