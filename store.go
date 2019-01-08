package main

import (
	"database/sql"
)

// Store has 2 methods, to get and to post
type Store interface {
	CreateEmployee(employee *Employee) error
	GetEmployees() ([]*Employee, error)
}

// The `dbStore` struct will implement the `Store` interface
// It also takes the sql DB connection object, which represents
// the database connection.
type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateEmployee(employee *Employee) error {
	// 'Employee' is a simple struct which has "name" and "team" attributes
	_, err := store.db.Query("INSERT INTO employees(name, team) VALUES ($1,$2)", employee.Name, employee.Team)
	return err
}

func (store *dbStore) GetEmployees() ([]*Employee, error) {

	rows, err := store.db.Query("SELECT name, team from employees")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create an empty array of employees
	employees := []*Employee{}
	for rows.Next() {
		// For each row returned by the table, create a pointer to a employee,
		employee := &Employee{}

		if err := rows.Scan(&employee.Name, &employee.Team); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

// The store variable is a package level variable that will be available for
// use throughout our application code
var store Store

// InitStore method is to initialize the store. This will
// typically be done at the beginning of our application (in this case, when the server starts up)
// This can also be used to set up the store as a mock, which we will be observing
// later on
func InitStore(s Store) {
	store = s
}
