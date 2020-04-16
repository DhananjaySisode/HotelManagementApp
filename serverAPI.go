package main

import (
	"HotelManagement/DataAccess"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllMenu(w http.ResponseWriter, r *http.Request) {
	var menuList, err = DataAccess.GetAllMenu()
	if err != nil {
		http.Error(w, "No Data Found!", http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(menuList)
	}
}

func GetMenuByCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["custId"]
	val, _ := strconv.Atoi(id)
	var menuList, _ = DataAccess.GetMenuByCustomer(val)
	if len(menuList) == 0 {
		http.Error(w, "No menu added for you", http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(menuList)
	}
}

func LoginCustomer(w http.ResponseWriter, r *http.Request) {
	var login DataAccess.Login
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&login)
	if err != nil {
		fmt.Println(err)
	}
	custDetails, _ := DataAccess.LoginCustomer(login)

	if custDetails == (DataAccess.CustomerMaster{}) {
		http.Error(w, "No such user exist", http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(custDetails)
	}
}

func RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	var cust DataAccess.CustomerMaster
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cust)
	if err != nil {
		fmt.Println(err)
	}
	custDetails, _ := DataAccess.RegisterCustomer(cust)

	if custDetails == (DataAccess.CustomerMaster{}) {
		http.Error(w, "Unable to add new customer", http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(custDetails)
	}
}

// func AddNewMenu(w http.ResponseWriter, r *http.Request) {

// }

func main() {

	var rout = mux.NewRouter()
	rout.HandleFunc("/menu", GetAllMenu).Methods("GET")
	rout.HandleFunc("/menu/{custId}", GetMenuByCustomer).Methods("GET")
	rout.HandleFunc("/login", LoginCustomer).Methods("POST")
	rout.HandleFunc("/register", RegisterCustomer).Methods("POST")
	//rout.HandleFunc("/addMenu", AddNewMenu).Methods("POST")

	http.ListenAndServe(":8080", rout)
}
