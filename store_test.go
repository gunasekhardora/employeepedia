package main

import (
	"database/sql"
	"testing"

	// The "testify/suite" package is used to make the test suite
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
	/*
		The suite is defined as a struct, with the store and db as its
		attributes. Any variables that are to be shared between tests in a
		suite should be stored as attributes of the suite instance
	*/
	store *dbStore
	db    *sql.DB
}

func (s *StoreSuite) SetupSuite() {
	/*
		The database connection is opened in the setup, and
		stored as an instance variable,
		as is the higher level `store`, that wraps the `db`
	*/
	connString := "dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	s.store = &dbStore{db: db}
}

func (s *StoreSuite) SetupTest() {
	/*
		We delete all entries from the table before each test runs, to ensure a
		consistent state before our tests run. In more complex applications, this
		is sometimes achieved in the form of migrations
	*/
	_, err := s.db.Query("DELETE FROM employees")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *StoreSuite) TearDownSuite() {
	// Close the connection after all tests in the suite finish
	s.db.Close()
}

// This is the actual "test" as seen by Go, which runs the tests defined below
func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}

func (s *StoreSuite) TestCreateEmployee() {
	// Create a employee through the store `CreateEmployee` method
	s.store.CreateEmployee(&Employee{
		Name: "test name",
		Team: "test team",
	})

	// Query the database for the entry we just created
	res, err := s.db.Query(`SELECT COUNT(*) FROM employees WHERE name='test name' AND team='test team'`)
	if err != nil {
		s.T().Fatal(err)
	}

	// Get the count result
	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	// Assert that there must be one entry with the properties of the employee that we just inserted (since the database was empty before this)
	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}
}

func (s *StoreSuite) TestGetEmployee() {
	// Insert a sample employee into the `employees` table
	_, err := s.db.Query(`INSERT INTO employees (name, team) VALUES('name','team')`)
	if err != nil {
		s.T().Fatal(err)
	}

	// Get the list of employees through the stores `GetEmployees` method
	employees, err := s.store.GetEmployees()
	if err != nil {
		s.T().Fatal(err)
	}

	// Assert that the count of employees received must be 1
	nEmployees := len(employees)
	if nEmployees != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", nEmployees)
	}

	// Assert that the details of the employee is the same as the one we inserted
	expectedEmployee := Employee{"name", "team"}
	if *employees[0] != expectedEmployee {
		s.T().Errorf("incorrect details, expected %v, got %v", expectedEmployee, *employees[0])
	}
}
