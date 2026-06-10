package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"autoparts/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PartHandler struct {
	db *pgxpool.Pool
}

func NewPartHandler(db *pgxpool.Pool) *PartHandler {
	return &PartHandler{db: db}
}

func (h *PartHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT p.id, p.category_id, p.name, p.article, p.description, p.created_at, c.name
		FROM parts p
		LEFT JOIN categories c ON c.id = p.category_id
		ORDER BY p.name
	`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var parts []models.PartWithCategory
	for rows.Next() {
		var p models.PartWithCategory
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Article, &p.Description, &p.CreatedAt, &p.CategoryName); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		p.Brands = []models.Brand{}
		parts = append(parts, p)
	}
	rows.Close()

	if parts == nil {
		parts = []models.PartWithCategory{}
		respondJSON(w, http.StatusOK, parts)
		return
	}

	brandRows, err := h.db.Query(r.Context(), `
		SELECT pb.part_id, b.id, b.name
		FROM part_brands pb
		JOIN brands b ON b.id = pb.brand_id
		ORDER BY b.name
	`)
	if err != nil {
		respondJSON(w, http.StatusOK, parts)
		return
	}
	defer brandRows.Close()

	idx := make(map[int]int, len(parts))
	for i, p := range parts {
		idx[p.ID] = i
	}
	for brandRows.Next() {
		var partID int
		var b models.Brand
		if err := brandRows.Scan(&partID, &b.ID, &b.Name); err != nil {
			continue
		}
		if i, ok := idx[partID]; ok {
			parts[i].Brands = append(parts[i].Brands, b)
		}
	}

	respondJSON(w, http.StatusOK, parts)
}

func (h *PartHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "неверный id")
		return
	}

	var p models.PartWithCategory
	err = h.db.QueryRow(r.Context(), `
		SELECT p.id, p.category_id, p.name, p.article, p.description, p.created_at, c.name
		FROM parts p
		LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1
	`, id).Scan(&p.ID, &p.CategoryID, &p.Name, &p.Article, &p.Description, &p.CreatedAt, &p.CategoryName)
	if err != nil {
		respondError(w, http.StatusNotFound, "запчасть не найдена")
		return
	}

	p.Brands = []models.Brand{}
	brandRows, err := h.db.Query(r.Context(), `
		SELECT b.id, b.name FROM part_brands pb
		JOIN brands b ON b.id = pb.brand_id
		WHERE pb.part_id = $1 ORDER BY b.name
	`, id)
	if err == nil {
		defer brandRows.Close()
		for brandRows.Next() {
			var b models.Brand
			if err := brandRows.Scan(&b.ID, &b.Name); err == nil {
				p.Brands = append(p.Brands, b)
			}
		}
	}

	respondJSON(w, http.StatusOK, p)
}

func (h *PartHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CategoryID  *int    `json:"category_id"`
		Name        string  `json:"name"`
		Article     *string `json:"article"`
		Description *string `json:"description"`
		BrandIDs    []int   `json:"brand_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "поле name обязательно")
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(r.Context())

	var p models.Part
	err = tx.QueryRow(r.Context(), `
		INSERT INTO parts (category_id, name, article, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, category_id, name, article, description, created_at
	`, req.CategoryID, req.Name, req.Article, req.Description,
	).Scan(&p.ID, &p.CategoryID, &p.Name, &p.Article, &p.Description, &p.CreatedAt)
	if err != nil {
		respondError(w, http.StatusConflict, "запчасть с таким артикулом уже существует")
		return
	}

	if _, err := tx.Exec(r.Context(), `INSERT INTO stock (part_id, quantity) VALUES ($1, 0)`, p.ID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, brandID := range req.BrandIDs {
		if _, err := tx.Exec(r.Context(), `INSERT INTO part_brands (part_id, brand_id) VALUES ($1, $2)`, p.ID, brandID); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if err := tx.Commit(r.Context()); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, p)
}

func (h *PartHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "неверный id")
		return
	}

	var req struct {
		CategoryID  *int    `json:"category_id"`
		Name        string  `json:"name"`
		Article     *string `json:"article"`
		Description *string `json:"description"`
		BrandIDs    []int   `json:"brand_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		respondError(w, http.StatusBadRequest, "поле name обязательно")
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(r.Context())

	var p models.Part
	err = tx.QueryRow(r.Context(), `
		UPDATE parts SET category_id=$1, name=$2, article=$3, description=$4
		WHERE id=$5
		RETURNING id, category_id, name, article, description, created_at
	`, req.CategoryID, req.Name, req.Article, req.Description, id,
	).Scan(&p.ID, &p.CategoryID, &p.Name, &p.Article, &p.Description, &p.CreatedAt)
	if err != nil {
		respondError(w, http.StatusNotFound, "запчасть не найдена")
		return
	}

	if _, err := tx.Exec(r.Context(), `DELETE FROM part_brands WHERE part_id=$1`, id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, brandID := range req.BrandIDs {
		if _, err := tx.Exec(r.Context(), `INSERT INTO part_brands (part_id, brand_id) VALUES ($1, $2)`, p.ID, brandID); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if err := tx.Commit(r.Context()); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, p)
}

func (h *PartHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "неверный id")
		return
	}

	tag, err := h.db.Exec(r.Context(), `DELETE FROM parts WHERE id=$1`, id)
	if err != nil || tag.RowsAffected() == 0 {
		respondError(w, http.StatusNotFound, "запчасть не найдена")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
