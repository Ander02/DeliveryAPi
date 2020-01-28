package main

import (
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
	"database/sql"

	_"github.com/denisenkom/go-mssqldb"
)

var db *sql.DB;
var err error

func main()  {
	router := mux.NewRouter();
	db, err = sql.Open("mssql", "Data Source=.\\SQLEXPRESS;Initial Catalog=Delivery;Integrated Security=True;Persist Security Info=False;Pooling=False;MultipleActiveResultSets=False;Connect Timeout=15;Encrypt=False;TrustServerCertificate=False")
	
	if err != nil {
		panic("Error on connect a database " + err.Error());
	} else {
		fmt.Println("Connected: ");
	}

	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))

	defer db.Close();
}


type User struct {
	NIF string `json:"NIF"`
	Name string `json:"Name"`		
	Email string `json:"Email"`
	Password string `json:"Password"`
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

func createUser(w http.ResponseWriter, r *http.Request) {
	statement, err := db.Prepare("INSER INTO Users (NIF, Name, Email, Password) VALUES(?, ?, ? ,?) ")

	if err != nil {
		panic(err.Error());
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error());
	}

	keyValue := make(map[string]string)
	json.Unmarshal(body, &keyValue);

	NIF := keyValue["NIF"]
	Name := keyValue["Name"]
	Email := keyValue["Email"]
	Password := keyValue["Password"]

	_, err = statement.Exec(NIF, Name, Email, Password)

	if err != nil {
		panic(err.Error());
	}

	fmt.Fprintf(w, "New user was created")
}