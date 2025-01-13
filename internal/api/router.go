package api

import (
	"net/http"
)

type Rest struct {
	server *http.Server
}

func New() *Rest {
	return &Rest{}
}

func (r *Rest) Start() {

}
