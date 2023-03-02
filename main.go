package main

import (
	"net/http"
	"outatime/middleware"
	"outatime/routes"
)

func main() {

	port := "3000"

	http.ListenAndServe(":"+port, middleware.CheckAuth(routes.GetRoutes()))

}
