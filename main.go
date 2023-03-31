package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// The main idea of this project is to provide an REST API that can manage Human-resource which will be stored in a file
// Example of resource file -> Kevin, David, Micheal
// Those employee's names will be put inside a slice

// TODO : Implement an HTTP server

func readFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening the resoure")
	}
	// This will close the file once the execution inside the function near the end
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Cannot read the file info")
	}

	fileBuffer := make([]byte, stat.Size())
	_, err = file.Read(fileBuffer)
	if err != nil {
		fmt.Println("Error reading the file content")
	}
	return string(fileBuffer)
}

func writeFile(path string, content string) {
	fileContentBuffer := []byte(content)
	os.WriteFile(path, fileContentBuffer, 0644)
}

func getEmployees(res http.ResponseWriter, req *http.Request) {
	// employeeList := strings.Split(string(fileBuffer), ",")
	employeeList := readFile("resource.txt")
	fmt.Fprintf(res, employeeList)
}

func searchEmployee(res http.ResponseWriter, req *http.Request) {
	employeeList := strings.Split(readFile("resource.txt"), ", ")
	args := mux.Vars(req)
	employeeName := args["name"]

	for _, employee := range employeeList {
		if string(employee) == employeeName {
			fmt.Fprintf(res, string(employee))
			return
		}
	}
}

func updateEmployee(res http.ResponseWriter, req *http.Request) {
	employeeList := strings.Split(readFile("resource.txt"), ", ")
	args := mux.Vars(req)
	employeeName := args["name"]
	newName := string(req.Form.Get("newName"))

	for i, employee := range employeeList {
		if string(employee) == employeeName {
			// Update employee name
			employeeList[i] = newName
			return
		}
	}

	var empStr string
	for _, employee := range employeeList {
		empStr += employee
	}

	writeFile(empStr, "resource.txt")
}

func deleteEmployee(res http.ResponseWriter, req *http.Request) {
	// Get EmployeeList from "resource.txt" file
	employeeList := strings.Split(readFile("resource.txt"), ", ")

	// Get employee name from the request parameters
	args := mux.Vars(req)
	employeeName := args["name"]

	var newEmpList []string
	for _, employee := range employeeList {
		if string(employee) != employeeName {
			newEmpList = append(newEmpList, employee)
		}
	}

	var empStr string
	for _, employee := range newEmpList {
		empStr += employee + ", "
	}

	writeFile("resource.txt", empStr)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/employees", getEmployees).Methods("GET")
	r.HandleFunc("/employees/{name}", searchEmployee).Methods("GET")
	r.HandleFunc("/employees/{name}", updateEmployee).Methods("PUT")
	r.HandleFunc("/employees/{name}", deleteEmployee).Methods("DELETE")
	fmt.Println("Server is listening on port : 5050")

	if err := http.ListenAndServe(":5050", r); err != nil {
		fmt.Println("Server error")
	}
}
