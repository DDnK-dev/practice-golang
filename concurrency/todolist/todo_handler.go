package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var rd *render.Render

var todoMap map[int]Todo
var lastID int = 0

func MakeWebHandler() http.Handler {
	todoMap = make(map[int]Todo)
	mux := mux.NewRouter()
	mux.Handle("/", http.FileServer(http.Dir("public")))
	mux.HandleFunc("/todos", GetTodoListHandler).Methods("GET")
	mux.HandleFunc("/todos", PostTodoHandler).Methods("POST")
	mux.HandleFunc("/todos/{id:[0-9]+}", RemoveTodoHandler).Methods("DELETE")
	mux.HandleFunc("/todos/{id:[0-9]+}", UpdateTodoHandler).Methods("PUT")
	return mux
}

func GetTodoListHandler(w http.ResponseWriter, r *http.Request) {
	list := make(Todos, 0)
	for _, todo := range todoMap {
		list = append(list, todo)
	}
	sort.Sort(list)
	_ = rd.JSON(w, http.StatusOK, list)
}

func PostTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	lastID++
	todo.ID = lastID
	todoMap[lastID] = todo
	_ = rd.JSON(w, http.StatusCreated, todo)
}

func RemoveTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if _, ok := todoMap[id]; ok {
		delete(todoMap, id)
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusNotFound, Success{false})
	}
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if todo, ok := todoMap[id]; ok {
		todo.Name = newTodo.Name
		todo.Completed = newTodo.Completed
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusBadRequest, Success{false})
	}
}
