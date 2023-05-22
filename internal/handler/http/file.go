package http

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/taranovegor/uploadly/internal/domain"
	"github.com/taranovegor/uploadly/pkg/config"
	"github.com/taranovegor/uploadly/pkg/dto"
	"mime/multipart"
	"net/http"
	"sync"
)

type FileHandler struct {
	contexts   config.FileContexts
	interactor domain.InterfaceFileInteractor
}

func NewFileHandler(
	contexts config.FileContexts,
	interactor domain.InterfaceFileInteractor,
) FileHandler {
	return FileHandler{
		contexts:   contexts,
		interactor: interactor,
	}
}

func (handler FileHandler) save(writer http.ResponseWriter, request *http.Request) {
	contextName := chi.URLParam(request, "context")
	context, err := handler.contexts.GetContext(contextName)
	if err != nil {
		errno(writer, http.StatusNotFound, err)

		return
	}

	if !context.IsContentLengthAcceptable(int(request.ContentLength)) {
		errno(
			writer,
			http.StatusRequestEntityTooLarge,
			errors.New(fmt.Sprintf("Content Length is `%dB` and over `%dB` limit", request.ContentLength, context.MaxContentLength)),
		)

		return
	}

	if err := request.ParseMultipartForm(int64(context.MaxContentLength)); err != nil {
		errno(writer, http.StatusInternalServerError, err)

		return
	}

	var filesLength int
	for _, fileHeaders := range request.MultipartForm.File {
		filesLength += len(fileHeaders)
	}
	if !context.IsNumberOfFilesAllowed(filesLength) {
		errno(
			writer,
			http.StatusRequestEntityTooLarge,
			errors.New(fmt.Sprintf("Attempt to upload `%d` files with a limit of `%d`", filesLength, context.MaxNumberOfFiles)),
		)

		return
	}

	responseMap := make(dto.FilesSaveMap)
	var wg sync.WaitGroup
	for field, fileHeaders := range request.MultipartForm.File {
		responseMap[field] = make(dto.FileSaveSlice, len(fileHeaders))
		for i, fileHeader := range fileHeaders {
			wg.Add(1)
			go func(field string, i int, fileHeader *multipart.FileHeader) {
				var fileResponse interface{}
				file, err := fileHeader.Open()
				if err == nil {
					fileResponse, err = handler.interactor.Save(contextName, fileHeader.Filename, file)
				}

				if err != nil {
					fileResponse = dto.NewError(0, err.Error())
				}

				responseMap[field][i] = fileResponse

				defer wg.Done()
			}(field, i, fileHeader)
		}
	}
	wg.Wait()

	response(writer, http.StatusOK, responseMap)
}

func (handler FileHandler) get(writer http.ResponseWriter, request *http.Request) {
	id, err := uuid.Parse(chi.URLParam(request, "id"))
	if err != nil {
		errno(writer, http.StatusBadRequest, err)

		return
	}

	file, err := handler.interactor.Get(id)
	if err != nil {
		errno(writer, http.StatusNotFound, err)

		return
	}

	response(writer, http.StatusOK, file)
}

func (handler FileHandler) delete(writer http.ResponseWriter, request *http.Request) {
	id, err := uuid.Parse(chi.URLParam(request, "id"))
	if err != nil {
		errno(writer, http.StatusBadRequest, err)

		return
	}

	if err := handler.interactor.Delete(id); err != nil {
		errno(writer, http.StatusNotFound, err)

		return
	}

	response(writer, http.StatusNoContent, "")
}
