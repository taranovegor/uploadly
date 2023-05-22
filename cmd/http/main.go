package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/taranovegor/uploadly/internal/builder"
	"github.com/taranovegor/uploadly/internal/config"
	"github.com/taranovegor/uploadly/internal/domain"
	"github.com/taranovegor/uploadly/internal/handler/http"
	"github.com/taranovegor/uploadly/internal/interactor"
	"github.com/taranovegor/uploadly/internal/repository"
	"github.com/taranovegor/uploadly/internal/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	netHttp "net/http"
)

func main() {
	fmt.Println(fmt.Sprintf("Uploadly! Version: %s", config.Version))

	cfg := config.Init()
	orm, err := gorm.Open(
		sqlite.Open(cfg.Database.Dsn),
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}

	if err := orm.AutoMigrate(
		&domain.File{},
	); err != nil {
		panic(err)
	}

	fileStorage := storage.NewLocal(cfg.Storage.StaticPath)
	fileBuilder := builder.NewFile(fileStorage)
	fileRepository := repository.NewFileRepository(orm)
	fileInteractor := interactor.NewFileInteractor(cfg.FileContext, fileRepository, fileBuilder, fileStorage)
	fileHandler := http.NewFileHandler(cfg.FileContext, fileInteractor)
	staticHandler := http.NewStaticHandler(cfg.Storage.StaticPath)

	router := chi.NewRouter()
	http.Init(
		router,
		fileHandler,
		staticHandler,
	)
	netHttp.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port), router)
}
