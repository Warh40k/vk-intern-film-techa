package handler

import "net/http"

func (h *Handler) ListFilms(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get method"))
}

func (h *Handler) SearchFilm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get method"))

}

func (h *Handler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Post method"))

}

func (h *Handler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete method"))

}

func (h *Handler) PatchFilm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Patch method"))

}

func (h *Handler) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Put method"))

}
