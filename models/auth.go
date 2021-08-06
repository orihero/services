package models

import "time"

const (
	ROLE_CUSTOMER   = "customer"
	ROLE_SPECIALIST = "specialist"
	ROLE_ADMIN = "admin"
)

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Verification struct {
	Id    uint `gorm:"primaryKey;autoIncrement:true"`
	Email string
	Code  string
}

type User struct {
	Id        uint      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate time.Time `json:"birth_date"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
}
