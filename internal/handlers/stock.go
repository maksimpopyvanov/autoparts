package handlers

import (
	"encoding/json"
	"net/http"

	"autoparts/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StockHandler struct {
	db *pgxpool.Pool
}

func NewStockHandler(db *pgxpool.Pool) *StockHandler {
	return &StockHandler{db: db}
}

func (h *StockHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT s.part_id, p.name, p.article, c.name, s.quantity
		FROM stock s
		JOIN parts p ON p.id = s.part_id
		LEFT JOIN categories c ON c.id = p.category_id
		ORDER BY p.name
	`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var items []models.Stock
	for rows.Next() {
		var s models.Stock
		if err := rows.Scan(&s.PartID, &s.PartName, &s.Article, &s.CategoryName, &s.Quantity); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		items = append(items, s)
	}
	if items == nil {
		items = []models.Stock{}
	}
	respondJSON(w, http.StatusOK, items)
}

type IncomeHandler struct {
	db *pgxpool.Pool
}

func NewIncomeHandler(db *pgxpool.Pool) *IncomeHandler {
	return &IncomeHandler{db: db}
}

func (h *IncomeHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT i.id, i.part_id, p.name, i.quantity, i.date::text, i.comment
		FROM income i
		JOIN parts p ON p.id = i.part_id
		ORDER BY i.date DESC, i.id DESC
	`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var items []models.Income
	for rows.Next() {
		var inc models.Income
		if err := rows.Scan(&inc.ID, &inc.PartID, &inc.PartName, &inc.Quantity, &inc.Date, &inc.Comment); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		items = append(items, inc)
	}
	if items == nil {
		items = []models.Income{}
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *IncomeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PartID   int     `json:"part_id"`
		Quantity int     `json:"quantity"`
		Date     *string `json:"date"`
		Comment  *string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PartID == 0 || req.Quantity <= 0 {
		respondError(w, http.StatusBadRequest, "поля part_id и quantity (>0) обязательны")
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(r.Context())

	var inc models.Income
	err = tx.QueryRow(r.Context(), `
		INSERT INTO income (part_id, quantity, date, comment)
		VALUES ($1, $2, COALESCE($3::date, CURRENT_DATE), $4)
		RETURNING id, part_id, quantity, date::text, comment
	`, req.PartID, req.Quantity, req.Date, req.Comment,
	).Scan(&inc.ID, &inc.PartID, &inc.Quantity, &inc.Date, &inc.Comment)
	if err != nil {
		respondError(w, http.StatusBadRequest, "запчасть не найдена")
		return
	}

	if _, err = tx.Exec(r.Context(), `
		UPDATE stock SET quantity = quantity + $1 WHERE part_id = $2
	`, req.Quantity, req.PartID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// получаем имя запчасти для ответа
	h.db.QueryRow(r.Context(), `SELECT name FROM parts WHERE id=$1`, req.PartID).Scan(&inc.PartName)

	respondJSON(w, http.StatusCreated, inc)
}

type OutcomeHandler struct {
	db *pgxpool.Pool
}

func NewOutcomeHandler(db *pgxpool.Pool) *OutcomeHandler {
	return &OutcomeHandler{db: db}
}

func (h *OutcomeHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT o.id, o.part_id, p.name, o.quantity, o.date::text, o.comment
		FROM outcome o
		JOIN parts p ON p.id = o.part_id
		ORDER BY o.date DESC, o.id DESC
	`)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var items []models.Outcome
	for rows.Next() {
		var out models.Outcome
		if err := rows.Scan(&out.ID, &out.PartID, &out.PartName, &out.Quantity, &out.Date, &out.Comment); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		items = append(items, out)
	}
	if items == nil {
		items = []models.Outcome{}
	}
	respondJSON(w, http.StatusOK, items)
}

func (h *OutcomeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PartID   int     `json:"part_id"`
		Quantity int     `json:"quantity"`
		Date     *string `json:"date"`
		Comment  *string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PartID == 0 || req.Quantity <= 0 {
		respondError(w, http.StatusBadRequest, "поля part_id и quantity (>0) обязательны")
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback(r.Context())

	var current int
	err = tx.QueryRow(r.Context(), `SELECT quantity FROM stock WHERE part_id=$1`, req.PartID).Scan(&current)
	if err != nil {
		respondError(w, http.StatusBadRequest, "запчасть не найдена на складе")
		return
	}
	if current < req.Quantity {
		respondError(w, http.StatusConflict, "недостаточно товара на складе")
		return
	}

	var out models.Outcome
	err = tx.QueryRow(r.Context(), `
		INSERT INTO outcome (part_id, quantity, date, comment)
		VALUES ($1, $2, COALESCE($3::date, CURRENT_DATE), $4)
		RETURNING id, part_id, quantity, date::text, comment
	`, req.PartID, req.Quantity, req.Date, req.Comment,
	).Scan(&out.ID, &out.PartID, &out.Quantity, &out.Date, &out.Comment)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if _, err = tx.Exec(r.Context(), `
		UPDATE stock SET quantity = quantity - $1 WHERE part_id = $2
	`, req.Quantity, req.PartID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// получаем имя запчасти для ответа
	h.db.QueryRow(r.Context(), `SELECT name FROM parts WHERE id=$1`, req.PartID).Scan(&out.PartName)

	respondJSON(w, http.StatusCreated, out)
}
