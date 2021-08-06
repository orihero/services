package utils

import (
	"../env"
	"../models"
	"crypto/tls"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
	"time"
)

//--------------HELPER FUNCTIONS---------------------

//set error message in Error struct
func SetError(err models.Error, message string) models.Error {
	err.IsError = true
	err.Message = message
	return err
}

//take password as input and generate new hash password from it
func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//compare plain password with hash password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//Generate JWT token
func GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte(env.SecretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("Something went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func SendEmailVerification(user *models.User, link string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "oriherodemos@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", "aka.orihero@gmail.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "Email verification to Services")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", fmt.Sprintf("Hi %s %s\nThis is email for you to confirm your email. Please visit link below  to confirm your email:\n\n%s", user.FirstName, user.LastName, link))

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "oriherodemos@gmail.com", "xitb6ldim")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
