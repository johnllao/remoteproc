package main

import (
	"net/http"
	"os"
)

type CSSHandler struct {
	Path string
}

func NewCSSHandler(p string) *CSSHandler {
	return &CSSHandler{
		Path: p,
	}
}

func (h *CSSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var b []byte
	b, err = os.ReadFile(h.Path)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/css")
	w.Write(b)
}
