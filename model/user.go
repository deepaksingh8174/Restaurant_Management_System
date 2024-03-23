package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	AdminUser    = "admin"
	SubAdminUser = "subAdmin"
	NormalUser   = "user"
)

type Register struct {
	Id       uuid.UUID
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type Login struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type Claims struct {
	Userid uuid.UUID `json:"userid"`
	jwt.StandardClaims
}

type Restaurant struct {
	Id        uuid.UUID
	Name      string    `json:"name" db:"name"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	Address   string    `json:"address" db:"address"`
	CreatedBy uuid.UUID `db:"created_by"`
}

type Dishes struct {
	Id        uuid.UUID
	Name      string    `json:"name" db:"name"`
	Cost      int       `json:"cost" db:"cost"`
	CreatedIn string    `json:"created_in" db:"created_in"`
	CreatedBy uuid.UUID `json:"created_by" db:"created_by"`
}

type Address struct {
	Id        uuid.UUID
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
}

type Distance struct {
	RestaurantId uuid.UUID `json:"rid"`
	AddressId    uuid.UUID `json:"aid"`
}

type LatAndLong struct {
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
}

//type Restaurants struct {
//	Id uuid.UUID
//	Name      string  `json:"name" db:"name"`
//	Latitude  float64 `json:"latitude" db:"latitude"`
//	Longitude float64 `json:"longitude" db:"longitude"`
//	Address   string  `json:"address" db:"address"`
//	CreatedBy uuid.UUID
//}
