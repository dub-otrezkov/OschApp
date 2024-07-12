package tasks

import (
	"html/template"
	"net/http"
)

type TasksApp struct{}

func New() *TasksApp {
	return &TasksApp{}
}

func (t TasksApp) Init() {
	http.HandleFunc("/tasks", t.tasks_list_page)
	http.HandleFunc("/tasks/{id}", t.task_page)
}

func (TasksApp) tasks_list_page(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./files/tasks/taskslist.html")

	if err != nil {
		panic(err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func (TasksApp) task_page(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./files/tasks/task.html")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, struct {
		Taskid string
	}{
		Taskid: r.PathValue("id"),
	})
	if err != nil {
		panic(err)
	}
}
