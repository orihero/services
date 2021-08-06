package controllers

import (
	"../database"
	"../models"
	"../utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		var err models.Error
		err = utils.SetError(err, "Error in reading payload.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	var dbuser models.User
	connection.Where("id = ?", user.Id).First(&dbuser)
	//Keep the password
	user.Password = dbuser.Password
	//update user details in database
	connection.Save(&user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	var users []models.User
	connection.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		var err models.Error
		err = utils.SetError(err, "Invalid USER_ID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)
	var dbuser models.User
	connection.Delete(&dbuser, id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted."))
}
