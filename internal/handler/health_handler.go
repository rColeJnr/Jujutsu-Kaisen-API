package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

func Health(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			_ = writeJson(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
			return
		}

		_ = writeJson(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
