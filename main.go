package main

import (
	"encoding/json"
	"fmt"
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
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// Обработчик для получения всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из слайса artists
	response, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(response)
}

// Обработчик для отправки задачи на сервер
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newTask.ID == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if newTask.Applications == nil || len(newTask.Applications) == 0 {
		newTask.Applications = []string{r.Header.Get("User-Agent")} // Добавление User-Agent, если applications пуст
	}

	if _, ok := tasks[newTask.ID]; ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")

	// Ищем задачу по ID
	task, ok := tasks[taskID]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Обработчик удаления задачи по ID
func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	if _, ok := tasks[taskID]; !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	delete(tasks, taskID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{taskID}", getTask)
	r.Delete("/tasks/{id}", deleteTaskByID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
