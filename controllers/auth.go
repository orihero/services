package controllers

import (
	"../database"
	"../models"
	"../utils"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	var authDetails models.Authentication

	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		var err models.Error
		err = utils.SetError(err, "Error in reading payload.")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var authUser models.User
	connection.Where("email = ?", authDetails.Email).First(&authUser)
	if authUser.Email == "" {
		var err models.Error
		err = utils.SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}

	check := utils.CheckPasswordHash(authDetails.Password, authUser.Password)

	if !check {
		var err models.Error
		err = utils.SetError(err, "Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}

	if !authUser.IsActive {
		var err models.Error
		err = utils.SetError(err, "Sorry your account has not been activated yet!")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := utils.GenerateJWT(authUser.Email, authUser.Role)
	if err != nil {
		var err models.Error
		err = utils.SetError(err, "Failed to generate token")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}

	token := models.Token{User: authUser, TokenString: validToken}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func Register(w http.ResponseWriter, r *http.Request) {
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
	connection.Where("email = ?", user.Email).First(&dbuser)

	//check email is alredy registered or not
	if dbuser.Email != "" {
		var err models.Error
		err = utils.SetError(err, "Email already in use")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}

	user.Password, err = utils.GenerateHashPassword(user.Password)
	if err != nil {
		log.Fatalln("Error in password hashing.")
	}
	//Make sure that user is  not active
	user.IsActive = false
	//insert user details in database
	connection.Create(&user)
	id, _ := uuid.NewV1()
	verification := models.Verification{
		Email: user.Email,
		Code:  id.String(),
	}
	connection.Create(&verification)
	utils.SendEmailVerification(&user, fmt.Sprintf("http://localhost:8080/verifications/%s", verification.Code))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Account was successfully created check your email for activation link!"))
}

func Verify(w http.ResponseWriter, r *http.Request) {
	//TODO check if email is already  activated
	params := mux.Vars(r)
	code := params["code"]
	if code == "" {
		var err models.Error
		err = utils.SetError(err, "Invalid link")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)
	verification := models.Verification{Code: code}
	connection.Find(&verification).First(&verification)
	if verification.Email == "" {
		var err models.Error
		err = utils.SetError(err, "Invalid link")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	user := models.User{Email: verification.Email}
	connection.Find(&user).First(&user)
	if user.Password == "" {
		var err models.Error
		err = utils.SetError(err, "Invalid link")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	user.IsActive = true
	connection.Save(&user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Account was successfully activated!"))
}
