package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	//"go.mongodb.org/mongo-driver/bson/primitive"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to he!\n")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	//starting with one ...
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(getCollection("todo")); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var todoId string

	todoId = vars["todoId"]

	/** err != nil {
		panic(err)
	}**/

	todo := RepoFindTodo(todoId)
	log.Print("todo111vars", vars)
	if len(todo.Name) > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/todos

*/

func TodoDestroy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var todoId string
	var jsonValidAction JsonValidAction
	todoId = vars["todoId"]
	deletedCount := DestroyTodo(todoId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if deletedCount == 1 {
		jsonValidAction.Bool = true
		jsonValidAction.Text = "action completed"
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(jsonValidAction); err != nil {
			panic(err)
		}
	}
}

func TodoCreate(w http.ResponseWriter, r *http.Request) {
	//log.Printf("TODODO 111")
	var todo Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := RepoCreateTodo(todo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

func TodoUpdate(w http.ResponseWriter, r *http.Request) {

	var todo Todo
	var tmpTodo Todo
	vars := mux.Vars(r)
	var todoId string
	todoId = vars["todoId"]
	todo = RepoFindTodo(todoId)
	ids, err := primitive.ObjectIDFromHex(todoId)
	//todo.Id = ids
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &tmpTodo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	tmpTodo.Id = ids
	//loop over tmp to assign todo with attr

	var todoInterface map[string]interface{}
	inrec, _ := json.Marshal(todo)
	json.Unmarshal(inrec, &todoInterface)

	var tmpTodoInterface map[string]interface{}
	tmpInrec, _ := json.Marshal(tmpTodo)
	json.Unmarshal(tmpInrec, &tmpTodoInterface)
	//fmt.Println("Map: ", todoInterface, tmpTodoInterface)
	for field, val := range todoInterface {
		fmt.Println("KV Pair: ", field == "_id", val)
		//check if tmp value is empty if so do not assign
		if field != "" {
			todoInterface[field] = tmpTodoInterface[field]
		}
	}

	mapstructure.Decode(todoInterface, &todo)
	if err != nil {
		panic(err)
	}

	todo = UpdateOne(todo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		panic(err)
	}

}
