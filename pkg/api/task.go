package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NekichGO/Graduate-work/pkg/db"
)

const limit = 50
const format = "20060102"
const formatFindDate = "02.01.2006"

type TaskResponse struct {
	Task []*db.Task `json:"tasks"`
}

func checkTask(task *db.Task) error {
	if task.Title == "" {
		return errors.New("title is required")
	}

	n := time.Now()
	now := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
	if task.Date == "" {
		task.Date = n.Format(format)
	}
	t, err := time.Parse(format, task.Date)
	if err != nil {
		return errors.New("failed to parse task date")
	}
	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	if t.Before(now) {
		if len(task.Repeat) == 0 {
			task.Date = n.Format(format)
		} else {
			task.Date = next
		}
	}
	return nil
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Println("failed to decode body:", err)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	err = checkTask(&task)
	if err != nil {
		log.Println("failed to validate task:", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Println("failed to add task:", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id})

}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []*db.Task
	var err error
	var date time.Time

	search := r.URL.Query().Get("search")

	if search == "" {
		tasks, err = db.Tasks(limit)
		if err != nil {
			log.Println("failed to get tasks:", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	} else {
		date, err = time.Parse(formatFindDate, search)
		if err != nil {
			tasks, err = db.SearchTasks(search, limit)
			if err != nil {
				log.Println("failed to search tasks:", err)
				writeError(w, http.StatusInternalServerError, "internal server error")
				return
			}
		} else {
			tasks, err = db.SearchTasksByDate(date.Format(format), limit)
			if err != nil {
				log.Println("failed to search tasks by date:", err)
				writeError(w, http.StatusInternalServerError, "internal server error")
				return
			}
		}
	}
	if tasks == nil {
		resEmpty := make([]*db.Task, 0)
		writeJSON(w, http.StatusOK, TaskResponse{Task: resEmpty})
	} else {
		writeJSON(w, http.StatusOK, TaskResponse{Task: tasks})
	}

}

func getTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("failed to parse id:", err)
		writeError(w, http.StatusBadRequest, "invalid id task")
		return
	}
	task, err := db.TaskByID(idInt)
	if err != nil {
		log.Println("failed to get task:", err)
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func putTaskByID(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Println("failed to decode body:", err)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	err = checkTask(&task)
	if err != nil {
		log.Println("failed to validate task:", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		log.Println("failed to update task:", err)
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}

func postTaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("id is required for task done")
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("failed to parse id for task done:", err)
		writeError(w, http.StatusBadRequest, "invalid id task")
		return
	}

	task, err := db.TaskByID(idInt)
	if err != nil {
		log.Println("failed to get task by id:", err)
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(idInt)
		if err != nil {
			log.Println("failed to delete task:", err)
			if errors.Is(err, db.ErrNotFound) {
				writeError(w, http.StatusNotFound, "task not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{})
	} else {
		n := time.Now()
		now := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)

		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Println("failed to obtain next date:", err)
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		err = db.UpdateTaskDate(next, idInt)
		if err != nil {
			log.Println("failed to update task date:", err)
			if errors.Is(err, db.ErrNotFound) {
				writeError(w, http.StatusNotFound, "task not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{})
	}
}

func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		log.Println("id is required for delete")
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println("failed to parse id for delete:", err)
		writeError(w, http.StatusBadRequest, "invalid id task")
		return
	}

	err = db.DeleteTask(idInt)
	if err != nil {
		log.Println("failed to delete task:", err)
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}
