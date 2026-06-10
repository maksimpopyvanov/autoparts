package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"autoparts/internal/db"
	"autoparts/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed frontend/dist
var staticFiles embed.FS

func main() {
	ctx := context.Background()

	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}
	defer pool.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	categories := handlers.NewCategoryHandler(pool)
	brands := handlers.NewBrandHandler(pool)
	parts := handlers.NewPartHandler(pool)
	stock := handlers.NewStockHandler(pool)
	income := handlers.NewIncomeHandler(pool)
	outcome := handlers.NewOutcomeHandler(pool)

	r.Route("/api", func(r chi.Router) {
		r.Get("/categories", categories.List)
		r.Post("/categories", categories.Create)
		r.Put("/categories/{id}", categories.Update)
		r.Delete("/categories/{id}", categories.Delete)

		r.Get("/brands", brands.List)
		r.Post("/brands", brands.Create)
		r.Put("/brands/{id}", brands.Update)
		r.Delete("/brands/{id}", brands.Delete)

		r.Get("/parts", parts.List)
		r.Post("/parts", parts.Create)
		r.Get("/parts/{id}", parts.Get)
		r.Put("/parts/{id}", parts.Update)
		r.Delete("/parts/{id}", parts.Delete)

		r.Get("/stock", stock.List)

		r.Get("/income", income.List)
		r.Post("/income", income.Create)

		r.Get("/outcome", outcome.List)
		r.Post("/outcome", outcome.Create)
	})

	// Раздаём собранный фронт
	distFS, err := fs.Sub(staticFiles, "frontend/dist")
	if err != nil {
		log.Fatalf("не удалось подключить статику: %v", err)
	}
	fileServer := http.FileServer(http.FS(distFS))
	r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		_, err := distFS.Open(req.URL.Path[1:])
		if err != nil {
			// SPA fallback — отдаём index.html для любого неизвестного пути
			req.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, req)
	})

	port := getenv("PORT", "8080")
	log.Printf("сервер запущен на :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
