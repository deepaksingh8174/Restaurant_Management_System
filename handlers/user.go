package handlers

import (
	"encoding/json"
	"example.com/database"
	"example.com/database/dbHelper"
	"example.com/model"
	"example.com/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"time"
)

//const SecretKey = "my_secret_key"

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var users model.Register
	decodeErr := json.NewDecoder(r.Body).Decode(&users)
	if decodeErr != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Status{Message: "error occurred while decoding"})
		return
	}
	exists, existErr := dbHelper.IsUserExists(users.Email)
	if existErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "failed to check user existence"})
		return
	}
	if exists {
		utils.RespondJSON(w, http.StatusConflict, utils.Status{Message: "user already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(users.Password)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred in hashing"})
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userId, userErr := dbHelper.CreateUser(tx, users.Name, users.Email, hashedPassword)
		if userErr != nil {
			utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "Unable to perform query"})
			return userErr
		}
		//creation of role
		roleErr := dbHelper.CreateRole(tx, userId, model.NormalUser)
		if roleErr != nil {
			utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "Unable to perform query"})
			return roleErr
		}
		return nil
	})
	if txErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "transaction cannot be committed"})
		return
	}
	utils.RespondJSON(w, http.StatusCreated, utils.Status{Message: "User Registered Successfully"})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user = model.Login{}
	decodeErr := json.NewDecoder(r.Body).Decode(&user)
	if decodeErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "enter appropriate value in postman"})
		return
	}
	// todo: give proper names
	userId, err := dbHelper.GetIdByPassword(user.Email, user.Password)
	if err != nil {

		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "failed to check database"})
		return
	}
	// todo: status bad request
	if userId == uuid.Nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Status{Message: "wrong credentials entered"})
		return
	}

	expirationTime := time.Now().Add(time.Hour * 24)

	claim := &model.Claims{
		Userid: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //what is the signifinance of unix function
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, tokenErr := token.SignedString([]byte(os.Getenv("SecretKey")))
	if tokenErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred in generation of token"})
		return
	}
	type tokenStruct struct {
		Token string `json:"token"`
	}
	t := tokenStruct{Token: tokenString}
	utils.RespondJSON(w, http.StatusOK, t)

}

func GetAllRestaurants(w http.ResponseWriter, r *http.Request) {
	restaurant, err := dbHelper.GetAllRestaurant()
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "Error occurred in fetching DB"})
		return
	}
	utils.RespondJSON(w, http.StatusOK, restaurant)

}

func GetDishes(w http.ResponseWriter, r *http.Request) {
	dishes, err := dbHelper.GetDish()
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred in fetching db"})
		return
	}
	encodeErr := json.NewEncoder(w).Encode(dishes)
	if encodeErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred in encoding"})
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func GetAddress(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*model.Claims)
	userId := claims.Userid
	address, err := dbHelper.GetAddress(userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred in fetching db"})
		return
	}
	utils.RespondJSON(w, http.StatusOK, address)

}

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var restaurant model.Restaurant
	decodeErr := json.NewDecoder(r.Body).Decode(&restaurant)
	if decodeErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred while decoding"})
		return
	}
	claims := r.Context().Value("claims").(*model.Claims)
	userId := claims.Userid
	creationErr := dbHelper.CreateRestaurant(userId, restaurant.Name, restaurant.Latitude, restaurant.Longitude, restaurant.Address)
	if creationErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "unable to perform query"})
		return
	}
	utils.RespondJSON(w, http.StatusCreated, utils.Status{Message: "restaurant table created"})
}

func CreateDishes(w http.ResponseWriter, r *http.Request) {
	var dish model.Dishes
	decodeErr := json.NewDecoder(r.Body).Decode(&dish)
	if decodeErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred while decoding"})
		return
	}
	creationErr := dbHelper.CreateDish(dish.Name, dish.Cost, dish.CreatedIn, dish.CreatedBy)
	if creationErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "unable to perform query"})
		return
	}
	utils.RespondJSON(w, http.StatusCreated, utils.Status{Message: "dishes table created"})
}

func GetRestaurant(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*model.Claims)
	userId := claims.Userid
	restaurant, err := dbHelper.GetRestaurant(userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "Error occurred in fetching DB"})
		return
	}
	utils.RespondJSON(w, http.StatusOK, restaurant)

	//todo: fix status error

}

func GetAllSubAdmin(w http.ResponseWriter, r *http.Request) {
	allSubAdmin, err := dbHelper.GetAllSubAdmin()
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred in fetching DB"})
		return
	}
	utils.RespondJSON(w, http.StatusOK, allSubAdmin)

}

func CreateSubAdmin(w http.ResponseWriter, r *http.Request) {
	var users model.Register
	decodeErr := json.NewDecoder(r.Body).Decode(&users)
	if decodeErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred while decoding"})
		return
	}
	exists, existErr := dbHelper.IsUserExists(users.Email)
	if existErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "failed to check user existence"})
		return
	}
	if exists {
		utils.RespondJSON(w, http.StatusBadRequest, utils.Status{Message: "user already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(users.Password)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred in hashing"})
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userId, userErr := dbHelper.CreateUser(tx, users.Name, users.Email, hashedPassword)
		if userErr != nil {
			utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "unable to perform query"})
			return userErr
		}
		roleErr := dbHelper.CreateRole(tx, userId, model.SubAdminUser)
		if roleErr != nil {
			utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "unable to perform query"})
			return roleErr
		}
		return nil
	})
	if txErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "transaction cannot be committed"})
		return
	}
	utils.RespondJSON(w, http.StatusCreated, utils.Status{Message: "SubAdmin Registered Successfully"})
}

func CreateAddress(w http.ResponseWriter, r *http.Request) {
	var address model.Address
	decodeErr := json.NewDecoder(r.Body).Decode(&address)
	if decodeErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred while decoding"})
		return
	}
	claims := r.Context().Value("claims").(*model.Claims)
	userId := claims.Userid
	creationErr := dbHelper.CreateAddress(userId, address.Latitude, address.Longitude)
	if creationErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "unable to perform query"})
		return
	}
	utils.RespondJSON(w, http.StatusCreated, utils.Status{Message: "address table created."})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*model.Claims)
	expirationTime := time.Now().Add(time.Millisecond * 1)
	claims.ExpiresAt = expirationTime.Unix()
}

func GetDishesByRestaurant(w http.ResponseWriter, r *http.Request) {
	//todo: give proper variable names
	restaurantId := r.URL.Query().Get("id")
	id, parseErr := uuid.Parse(restaurantId)
	if parseErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "Error when parsing into uuid"})
		return
	}
	dish, err := dbHelper.GetDishes(id)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred in fetching db"})
		return
	}
	utils.RespondJSON(w, http.StatusOK, dish)

}

func GetDistance(w http.ResponseWriter, r *http.Request) {
	var address model.Distance
	decodeErr := json.NewDecoder(r.Body).Decode(&address)
	if decodeErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred while decoding"})
		return
	}
	distance, distanceErr := dbHelper.CalculateDistance(address.AddressId, address.RestaurantId)
	if distanceErr != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "error occurred in calculating distance"})
		return
	}
	utils.RespondJSON(w, http.StatusOK, distance)

}
