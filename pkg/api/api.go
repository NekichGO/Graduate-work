package api

import (
	"net/http"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/nextdate", getNextDate)
	mux.HandleFunc("POST /api/task", postTask)
	mux.HandleFunc("GET /api/tasks", getAllTasks)
	mux.HandleFunc("GET /api/task", getTaskByID)
	mux.HandleFunc("PUT /api/task", putTaskByID)
	mux.HandleFunc("POST /api/task/done", postTaskDone)
	mux.HandleFunc("DELETE /api/task", deleteTaskByID)
}
