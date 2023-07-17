package function

import (
	"net/http"

	"github.com/peksinsara/e-voting-RDBMS/database"
)

func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("Admin-Username")
		password := r.Header.Get("Admin-Password")

		db := database.GetDB()

		var isAdmin bool
		err := db.QueryRow("SELECT is_admin FROM User WHERE admin_username = ? AND admin_password = ?", username, password).Scan(&isAdmin)
		if err != nil || !isAdmin {
			http.Error(w, "Invalid admin credentials", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
