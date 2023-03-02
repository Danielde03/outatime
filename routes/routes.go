package routes

import (
	"net/http"
	"outatime/handlers"
)

// Get routes Mux
func GetRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// connect to public folder
	fs1 := http.FileServer(http.Dir("./templates/public"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs1))

	// form routes
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)

	// page routes
	mux.HandleFunc("/", handlers.Root)

	return mux
}
