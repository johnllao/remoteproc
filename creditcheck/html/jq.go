package main

import (
	"net/http"
	"os"
)

type JQHandler struct {
	Path string
}

func NewJQHandler(p string) *JQHandler {
	return &JQHandler{
		Path: p,
	}
}

func (h *JQHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var b []byte
	b, err = os.ReadFile(h.Path)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/javascript")
	w.Write(b)
}
