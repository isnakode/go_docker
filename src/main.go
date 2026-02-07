package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"com.isnakode.hello/db"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

//go:generate docker run --rm -v .:/src -w /src sqlc/sqlc generate

func main() {
	_ = godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is required. Application cannot start without a database connection.")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatalf("Gagal membuat database pool: %v\n", err)
	}
	defer pool.Close()
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Database is not responding: %v\n", err)
	}
	queries := db.New(pool)
	app := http.NewServeMux()

	app.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		authors, err := queries.ListAuthors(ctx)

		msg := "Data berhasil diambil"
		if err != nil {
			log.Printf("error %v\n", err)
			msg = "error"
		}

		data := map[string]any{
			"message": msg,
			"data": map[string]any{
				"authors": authors,
			},
		}
		json.NewEncoder(w).Encode(data)
	})
	app.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := queries.CreateAuthor(ctx, db.CreateAuthorParams{
			Name: "isnaini",
			Bio:  pgtype.Text{String: "aja sendiri", Valid: true},
		})

		msg := "Data berhasil diambil"
		if err != nil {
			msg = "error"
		}

		data := map[string]any{
			"message": msg,
		}
		json.NewEncoder(w).Encode(data)
	})
	app.HandleFunc("GET /text", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("hi gemini"))
	})
	app.HandleFunc("GET /html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`
        <html>
            <body>
                <h1 style="color: blue;">Halo Dunia!</h1>
                <p>Ini dikirim dari Go server</p>
            </body>
        </html>`,
		))
	})
	app.HandleFunc("GET /json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		data := map[string]any{
			"status":  "success",
			"message": "Data berhasil diambil",
			"items":   []string{"Gopher", "Python", "Rust"},
			"count":   3,
		}
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", app)
}
