package handler

import (
	"github.com/mp1947/ya-url-shortener/internal/service"
)

const (
	contentTypePlain  = "text/plain; charset=utf-8"
	contentTypeJSON   = "application/json; charset=utf-8"
	requestBindingErr = "invalid request: error parsing request params"
	requestBodyGetErr = "error getting request body"
)

// HandlerService provides HTTP handler methods and holds a reference to the core service logic.
// Service is the main business logic layer that HandlerService delegates requests to.
type HandlerService struct {
	Service service.Service
}
