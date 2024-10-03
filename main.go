package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` // "pending" or "completed"
}

var tasks []Task
var nextID = 1


func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	json.NewDecoder(r.Body).Decode(&task)
	task.ID = nextID
	nextID++
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(task) // Return the created task as a response
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks) // Return all tasks
}


func getTaskByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"]) // Convert ID to integer
	for _, task := range tasks {
		if task.ID == id {
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}


func updateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, task := range tasks {
		if task.ID == id {
			var updatedTask Task
			json.NewDecoder(r.Body).Decode(&updatedTask)
			updatedTask.ID = id // Keep the same ID
			tasks[i] = updatedTask
			json.NewEncoder(w).Encode(updatedTask) // Return updated task
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...) // Remove task from slice
			json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted"})
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	// Initialize the list with one predefined task "Buy groceries"
	tasks = append(tasks, Task{
		ID:          nextID,
		Title:       "Buy groceries",
		Description: "Get milk, bread, and eggs",
		Status:      "pending",
	})
	nextID++ // Increment ID for next task

	// Set up the router
	r := mux.NewRouter()

	// Define routes for each CRUD operation
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks", getAllTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", getTaskByID).Methods("GET")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	// Start the server on port 8080
	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", r)
}
