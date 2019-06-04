package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var todo Todo

func TestIndexkHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/index", nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Index)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "Welcome to he!\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestTodoCreateHandler(t *testing.T) {

	var jsonStr = []byte(`{"name":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TodoCreate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(io.LimitReader(rr.Body, 1048576))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &todo)

	// Check the response body is what we expect.
	expected := "Buy cheese and bread for breakfast."
	if todo.Name != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			todo.Name, expected)
	}

}

func TestTodoIndexHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/todos", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TodoIndex)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(io.LimitReader(rr.Body, 1048576))

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &todos)

	// Check the response body is what we expect.
	fmt.Println(todos[0])
	expected := "Buy cheese and bread for breakfast."
	if todos[0].Name != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestTodoShowHandler(t *testing.T) {

	//req, err := http.NewRequest("GET", "/todos", nil)
	req, err := http.NewRequest("GET", "todos/", nil)
	fmt.Println(todo.Id.String())
	req = mux.SetURLVars(req, map[string]string{"todoId": todo.Id.Hex()})
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TodoShow)

	handler.ServeHTTP(rr, req)
	fmt.Print("rr", rr.Body.String())

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(io.LimitReader(rr.Body, 1048576))

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &todo)

	// Check the response body is what we expect.
	expected := "Buy cheese and bread for breakfast."
	if todo.Name != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestTodoUpdateHandler(t *testing.T) {

	//req, err := http.NewRequest("GET", "/todos", nil)
	var jsonStr = []byte(`{"name":"Buy call possible drugs"}`)
	req, err := http.NewRequest("PUT", "todos/", bytes.NewBuffer(jsonStr))
	fmt.Println(todo.Id.String())
	req = mux.SetURLVars(req, map[string]string{"todoId": todo.Id.Hex()})
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TodoUpdate)

	handler.ServeHTTP(rr, req)
	fmt.Print("rr", rr.Body.String())

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(io.LimitReader(rr.Body, 1048576))

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &todo)

	// Check the response body is what we expect.
	expected := "Buy call possible drugs"
	if todo.Name != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//TodoDestroy

func TestTodoDestroyHandler(t *testing.T) {

	var validAction JsonValidAction

	req, err := http.NewRequest("DELETE", "/todos", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{"todoId": todo.Id.Hex()})
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TodoDestroy)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body, err := ioutil.ReadAll(io.LimitReader(rr.Body, 1048576))

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &validAction)

	// Check the response body is what we expect.
	expected := "action completed"
	if validAction.Text != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
