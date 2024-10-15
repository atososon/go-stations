package router

import (
    "database/sql"
    "github.com/TechBowl-japan/go-stations/handler"
    "github.com/TechBowl-japan/go-stations/service"
    "net/http"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
    // Create a new TODOService
    todoService := service.NewTODOService(todoDB)

    // register routes
    mux := http.NewServeMux()
    mux.Handle("/healthz", handler.NewHealthzHandler())
    mux.Handle("/todos", handler.NewTODOHandler(todoService))
    return mux
}
