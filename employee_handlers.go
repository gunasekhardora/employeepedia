package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Employee Structure
type Employee struct {
	Name string `json:"name"`
	Team string `json:"team"`
}

func getEmployeeHandler(w http.ResponseWriter, r *http.Request) {

	employees, err := store.GetEmployees()

	employeeListBytes, err := json.Marshal(employees)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(employeeListBytes)
}

func createEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	employee := Employee{}

	err := r.ParseForm()

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	employee.Name = r.Form.Get("name")
	employee.Team = r.Form.Get("team")

	err = store.CreateEmployee(&employee)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/assets/", http.StatusFound)
}
