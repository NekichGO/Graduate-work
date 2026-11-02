package server

import (
	"net/http"

	"github.com/NekichGO/Graduate-work/pkg/api"
)

func Run(port, dir string) error {
	mux := http.NewServeMux()

	api.Init(mux)
	mux.Handle("GET /", http.FileServer(http.Dir(dir)))
	return http.ListenAndServe(":"+port, mux)
}
