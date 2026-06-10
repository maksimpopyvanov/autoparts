package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"autoparts/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BrandHandler struct {
	db *pgxpool.Pool
}

func NewBrandHandler(db *pgxpool.Pool) *BrandHandler {
	return &BrandHandler{db: db}
}

func (h *BrandHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `SELECT id, name FROM brands ORDER BY name`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var brands []models.Brand
	for rows.Next() {
		var b models.Brand
		if err := rows.Scan(&b.ID, &b.Name); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		brands = append(brands, b)
	}
	if brands == nil {
		brands = []models.Brand{}
	}
	respondJSON(w, http.StatusOK, brands)
}

func (h *BrandHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "поле name обязательно")
		return
	}

	var b models.Brand
	err := h.db.QueryRow(r.Context(),
		`INSERT INTO brands (name) VALUES ($1) RETURNING id, name`, req.Name,
	).Scan(&b.ID, &b.Name)
	if err != nil {
		respondError(w, http.StatusConflict, "марка с таким именем уже существует")
		return
	}
	respondJSON(w, http.StatusCreated, b)
}

func (h *BrandHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var b models.Brand
	err = h.db.QueryRow(r.Context(),
		`UPDATE brands SET name=$1 WHERE id=$2 RETURNING id, name`, req.Name, id,
	).Scan(&b.ID, &b.Name)
	if err != nil {
		respondError(w, http.StatusNotFound, "марка не найдена")
		return
	}
	respondJSON(w, http.StatusOK, b)
}

func (h *BrandHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "неверный id")
		return
	}

	tag, err := h.db.Exec(r.Context(), `DELETE FROM brands WHERE id=$1`, id)
	if err != nil || tag.RowsAffected() == 0 {
		respondError(w, http.StatusNotFound, "марка не найдена")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
