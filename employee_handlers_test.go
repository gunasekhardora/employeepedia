package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestGetEmployeesHandler(t *testing.T) {
	// Initialize the mock store
	mockStore := InitMockStore()

	/* Define the data that we want to return when the mocks `GetEmployees` method is
	called
	Also, we expect it to be called only once
	*/
	mockStore.On("GetEmployees").Return([]*Employee{{"Dora", "Platforms"}}, nil).Once()

	req, err := http.NewRequest("GET", "", nil)

	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(getEmployeeHandler)

	// Now, when the handler is called, it should cal our mock store, instead of
	// the actual one
	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := Employee{"Dora", "Platforms"}
	b := []Employee{}
	err = json.NewDecoder(recorder.Body).Decode(&b)

	if err != nil {
		t.Fatal(err)
	}

	actual := b[0]

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}

	// the expectations that we defined in the `On` method are asserted here
	mockStore.AssertExpectations(t)
}

func TestCreateEmployeesHandler(t *testing.T) {

	mockStore := InitMockStore()
	/*
	 Similarly, we define our expectations for th `CreateEmployee` method.
	 We expect the first argument to the method to be the employee struct
	 defined below, and tell the mock to return a `nil` error
	*/
	mockStore.On("CreateEmployee", &Employee{"Dora", "Platforms"}).Return(nil)

	form := newCreateEmployeeForm()
	req, err := http.NewRequest("POST", "", bytes.NewBufferString(form.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(createEmployeeHandler)

	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	mockStore.AssertExpectations(t)
}

func newCreateEmployeeForm() *url.Values {
	form := url.Values{}
	form.Set("name", "Dora")
	form.Set("team", "Platforms")
	return &form
}
