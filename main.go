package main

import (
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
	"database/sql"

	_"github.com/go-sql-driver/mysql"
)

var db *sql.DB;
var err error

func main()  {
	router := mux.NewRouter();
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1)/delivery")
	
	if err != nil {
		panic("Error on connect a database " + err.Error());
	} else {
		fmt.Println("Connected: ");
	}

	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{NIF}", getUserByNIF).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{NIF}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{NIF}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))

	defer db.Close();
}


type User struct {
	NIF string `json:"NIF"`
	Name string `json:"Name"`		
	Email string `json:"Email"`
	PasswordHash string `json:"Password"`
}

func getUsers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var users []User;

	result, err := db.Query("SELECT NIF, Name, Email FROM Users");
	
	if err != nil {
		panic(err.Error());
	}

	defer result.Close();

	for result.Next(){
		var user User
		err := result.Scan(&user.NIF, &user.Name, &user.Email)

		if err != nil {
			panic(err.Error());
		}
		
		users = append(users, user);
	}

	json.NewEncoder(writer).Encode(users);
}

func getUserByNIF(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request);

	var users []User;

	result, err := db.Query("SELECT NIF, Name, Email FROM Users WHERE NIF LIKE ?", "%" + params["NIF"] + "%");
	
	if err != nil {
		panic(err.Error());
	}

	defer result.Close();

	for result.Next(){
		var user User
		err := result.Scan(&user.NIF, &user.Name, &user.Email)

		if err != nil {
			panic(err.Error());
		}
		
		users = append(users, user);
	}

	json.NewEncoder(writer).Encode(users);
}

func createUser(w http.ResponseWriter, request *http.Request) {
	statement, err := db.Prepare("INSERT INTO Users (NIF, Name, Email, PasswordHash) VALUES(?, ?, ? ,?) ")

	if err != nil {
		panic(err.Error());
	}

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		panic(err.Error());
	}

	keyValue := make(map[string]string)
	json.Unmarshal(body, &keyValue);

	NIF := keyValue["NIF"]
	Name := keyValue["Name"]
	Email := keyValue["Email"]
	PasswordHash := keyValue["PasswordHash"]

	_, err = statement.Exec(NIF, Name, Email, PasswordHash)

	if err != nil {
		panic(err.Error());
	}

	fmt.Fprintf(w, "New user was created")
}

func updateUser(w http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	statement, err := db.Prepare("UPDATE Users SET Name = ?,  Email = ?, PasswordHash = ? WHERE NIF = ?")

	if err != nil {
		panic(err.Error());
	}

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		panic(err.Error());
	}

	keyValue := make(map[string]string)
	json.Unmarshal(body, &keyValue);

	NIF := params["NIF"]
	Name := keyValue["Name"]
	Email := keyValue["Email"]
	PasswordHash := keyValue["PasswordHash"]

	_, err = statement.Exec(Name, Email, PasswordHash, NIF)

	if err != nil {
		panic(err.Error());
	}

	fmt.Fprintf(w, "User " + NIF + " was edited")
}

func deleteUser(w http.ResponseWriter, request *http.Request) {

	params := mux.Vars(request);
	NIF := params["NIF"]
	statement, err := db.Prepare("DELETE FROM Users WHERE NIF = ?")
	statement.Exec(NIF);

	if err != nil {
		panic(err.Error());
	}

	fmt.Fprintf(w, "User " + NIF + " was deleted")
}
