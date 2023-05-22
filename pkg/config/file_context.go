package config

import (
	"errors"
	"fmt"
)

type FileContexts map[string]FileContext

type FileContext struct {
	MaxNumberOfFiles uint     `yaml:"max_number_of_files"`
	MaxContentLength uint     `yaml:"max_content_length"`
	MaxFileSize      uint     `yaml:"max_file_size"`
	AllowedMimeTypes []string `yaml:"allowed_mime_types"`
}

func (context FileContexts) GetContext(name string) (FileContext, error) {
	if val, available := context[name]; available {
		return val, nil
	}

	return FileContext{}, errors.New(fmt.Sprintf("Context `%s` not found", name))
}

func (context FileContext) IsMimeTypeAllowed(mimeType string) bool {
	if context.AllowedMimeTypes == nil {
		return true
	}

	for _, ext := range context.AllowedMimeTypes {
		if mimeType == ext {
			return true
		}
	}

	return false
}

func (context FileContext) IsContentLengthAcceptable(length int) bool {
	return context.MaxContentLength == 0 || uint(length) <= context.MaxContentLength
}

func (context FileContext) IsFileSizeAcceptable(length int) bool {
	return context.MaxFileSize == 0 || uint(length) <= context.MaxFileSize
}

func (context FileContext) IsNumberOfFilesAllowed(length int) bool {
	return context.MaxNumberOfFiles == 0 || uint(length) <= context.MaxNumberOfFiles
}
