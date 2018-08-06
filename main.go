package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type TaskItem struct {
	Id          int    `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type TaskServer struct {
	Tasks  []TaskItem
	LastId int
}

type DataResponse struct {
	Status string      `json:"status"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func NewTaskServer() *TaskServer {
	return &TaskServer{
		Tasks:  make([]TaskItem, 0),
		LastId: 10,
	}
}

func (t *TaskServer) HandleCreateTask(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(404)
		rw.Write([]byte("Not Found"))
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeJson(rw, 400, &DataResponse{Status: StatusError, Error: err.Error()})
		return
	}

	item := TaskItem{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		writeJson(rw, 400, &DataResponse{Status: StatusError, Error: err.Error()})
		return
	}
	if item.Id != 0 {
		writeJson(rw, 400, &DataResponse{Status: StatusError, Error: "A value for the ID can't be provided"})
		return
	}
	if item.Title == "" {
		writeJson(rw, 400, &DataResponse{Status: StatusError, Error: "A title must be provided"})
		return
	}
	if item.Description == "" {
		writeJson(rw, 400, &DataResponse{Status: StatusError, Error: "A description must be provided"})
		return
	}
	item.Id = t.LastId
	t.LastId++
	t.Tasks = append(t.Tasks, item)
	writeJson(rw, 201, &DataResponse{Status: StatusSuccess, Data: item})
}

func (t *TaskServer) HandleDeleteTask(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		rw.WriteHeader(404)
		rw.Write([]byte("Not Found"))
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeJson(rw, 400, &DataResponse{Status: StatusError, Error: err.Error()})
		return
	}

	item := TaskItem{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		writeJson(rw, 400, &DataResponse{Status: StatusError, Error: err.Error()})
		return
	}
	if item.Id == 0 {
		writeJson(rw, 400, &DataResponse{Status: StatusError, Error: "A task ID must be provided"})
		return
	}

	for i, task := range t.Tasks {
		if task.Id == item.Id {
			t.Tasks = append(t.Tasks[0:i], t.Tasks[i+1:]...)
			writeJson(rw, 200, &DataResponse{Status: StatusSuccess})
			return
		}
	}

	writeJson(rw, 404, &DataResponse{Status: StatusError, Error: fmt.Sprintf("The item with id %d does not exist", item.Id)})
}

func (t *TaskServer) HandleListTasks(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		rw.WriteHeader(404)
		rw.Write([]byte("Not Found"))
		return
	}
	writeJson(rw, 200, &DataResponse{Status: StatusSuccess, Data: t.Tasks})
}

func (t *TaskServer) HandleHealth(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		rw.WriteHeader(404)
		rw.Write([]byte("Not Found"))
		return
	}
	writeJson(rw, 200, &DataResponse{Status: "ok"})
}

func writeJson(rw http.ResponseWriter, code int, resp interface{}) {
	rw.Header().Set("Content-Type", "applicaiton/json; charset=UTF-8")
	rw.WriteHeader(code)
	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		panic(err)
	}
}

func main() {
	ts := NewTaskServer()
	http.HandleFunc("/", ts.HandleHealth)
	http.HandleFunc("/api/create_task", ts.HandleCreateTask)
	http.HandleFunc("/api/delete_task", ts.HandleDeleteTask)
	http.HandleFunc("/api/list_tasks", ts.HandleListTasks)
	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
