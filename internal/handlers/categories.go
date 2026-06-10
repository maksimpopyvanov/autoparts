package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"autoparts/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryHandler struct {
	db *pgxpool.Pool
}

func NewCategoryHandler(db *pgxpool.Pool) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `SELECT id, name FROM categories ORDER BY name`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		categories = append(categories, c)
	}
	if categories == nil {
		categories = []models.Category{}
	}
	respondJSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "поле name обязательно")
		return
	}

	var c models.Category
	err := h.db.QueryRow(r.Context(),
		`INSERT INTO categories (name) VALUES ($1) RETURNING id, name`, req.Name,
	).Scan(&c.ID, &c.Name)
	if err != nil {
		respondError(w, http.StatusConflict, "категория с таким именем уже существует")
		return
	}
	respondJSON(w, http.StatusCreated, c)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "неверный id")
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "поле name обязательно")
		return
	}

	var c models.Category
	err = h.db.QueryRow(r.Context(),
		`UPDATE categories SET name=$1 WHERE id=$2 RETURNING id, name`, req.Name, id,
	).Scan(&c.ID, &c.Name)
	if err != nil {
		respondError(w, http.StatusNotFound, "категория не найдена")
		return
	}
	respondJSON(w, http.StatusOK, c)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "неверный id")
		return
	}

	tag, err := h.db.Exec(r.Context(), `DELETE FROM categories WHERE id=$1`, id)
	if err != nil || tag.RowsAffected() == 0 {
		respondError(w, http.StatusNotFound, "категория не найдена")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
