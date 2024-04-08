package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postman",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// getTasks возвращает все задачи, которые хранятся в мапе tasks.
func getTasks(w http.ResponseWriter, _ *http.Request) {

	resp, err := json.Marshal(tasks)
	if err != nil {
		log.Printf("Error json.Marshal: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Error json.Marshal: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// createTask принимает задачу в теле запроса и сохраняет её в мапе.
func createTasks(w http.ResponseWriter, r *http.Request) {
	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		log.Printf("Error ТNewDecoder().Decode(): %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// getTaskForId возвращает задачу с указанным в запросе пути ID.
func getTaskId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	TaskId, ok := tasks[id]
	if !ok {
		log.Print("Error: TaskId not found")
		http.Error(w, "getTaskForId: TaskId not found", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(TaskId)
	if err != nil {
		log.Printf("Error json.Marshal: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Error w.Write response: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// deleteTaskForId удаляет задачу из мапы по её ID.
func deleteTaskForId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	TaskId, ok := tasks[id]

	if !ok {
		log.Print("Error: TaskId not found")
		http.Error(w, "deleteTaskForId: TaskId not found", http.StatusBadRequest)
		return
	}

	delete(tasks, TaskId.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...
	r.Get("/tasks", getTasks)
	r.Post("/tasks", createTasks)
	r.Get("/tasks/{id}", getTaskId)
	r.Delete("/tasks/{id}", deleteTaskForId)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
