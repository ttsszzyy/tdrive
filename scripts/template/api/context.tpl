package svc

import (
	{{.configImport}}
    "github.com/go-playground/validator/v10"
)

type ServiceContext struct {
	Config {{.config}}
    Validator *validator.Validate
	{{.middleware}}
}

func NewServiceContext(c {{.config}}) *ServiceContext {
	return &ServiceContext{
		Config: c,
        Validator: validator.New(),
		{{.middlewareAssignment}}
	}
}
