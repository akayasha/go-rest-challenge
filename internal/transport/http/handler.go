package http

import (
	"encoding/json"
	"errors"
	"go-rest-challenge/internal/domain"
	"go-rest-challenge/internal/repository"
	"io"
	"net/http"
	"strconv"

	"go-rest-challenge/internal/usecase"
)

type Handler struct {
	usecase *usecase.BookUsecase
}

func NewHandler(u *usecase.BookUsecase) *Handler {
	return &Handler{usecase: u}
}

// Level 1
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true}`))
}

// Level 2
func (h *Handler) Echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid body"}`))
		return
	}

	// If empty body → return {}
	if len(body) == 0 {
		w.Write([]byte(`{}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// Level 5
//func (h *Handler) Token(w http.ResponseWriter, r *http.Request) {
//	writeJSON(w, 200, map[string]string{"token": "supertoken"})
//}

// Level 3
func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Year   int    `json:"year"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, 400, "invalid body")
		return
	}

	book, err := h.usecase.Create(input.Title, input.Author, input.Year)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}

	writeJSON(w, 201, book)
}
func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	books := h.usecase.GetAll(author, page, limit)

	// FORCE an empty slice if it's nil
	// This turns "null" into "[]" in the JSON output
	if books == nil {
		books = []domain.Book{}
	}

	writeJSON(w, 200, books)
}

func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	book, err := h.usecase.GetByID(id)
	if err != nil {
		writeError(w, 404, "not found")
		return
	}

	writeJSON(w, 200, book)
}

// Level 4
func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}

	var input struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Year   int    `json:"year"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, 400, "invalid body")
		return
	}

	book, err := h.usecase.Update(id, input.Title, input.Author, input.Year)
	if err != nil {

		// 🔥 THIS IS THE IMPORTANT PART
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, 404, "not found")
			return
		}

		writeError(w, 400, err.Error())
		return
	}

	writeJSON(w, 200, book)
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Token(w http.ResponseWriter, r *http.Request) {
	var req authRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid json")
		return
	}

	// simple hardcoded auth (for challenge)
	if req.Username != "admin" || req.Password != "password" {
		writeError(w, 401, "invalid credentials")
		return
	}

	resp := authResponse{
		Token: "supertoken",
	}

	writeJSON(w, 200, resp)
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := h.usecase.Delete(id); err != nil {
		writeError(w, 404, "not found")
		return
	}

	writeJSON(w, 200, map[string]string{"message": "deleted"})
}
