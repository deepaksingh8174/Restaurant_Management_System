package dbHelper

import (
	"database/sql"
	"example.com/database"
	"example.com/model"
	"example.com/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func IsUserExists(email string) (bool, error) {
	SQL := `SELECT id from users WHERE email = TRIM(LOWER($1))`
	var id uuid.UUID
	err := database.Todo.Get(&id, SQL, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, nil
}

func CreateUser(tx *sqlx.Tx, name, email, password string) (uuid.UUID, error) {
	SQL := `INSERT INTO users(name,email,password) VALUES ($1,TRIM(LOWER($2)),$3) RETURNING id`
	var id uuid.UUID
	err := tx.QueryRowx(SQL, name, email, password).Scan(&id)
	if err != nil {
		return id, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, nil
	}
	return id, nil
}

func CreateRole(tx *sqlx.Tx, userid uuid.UUID, role string) error {
	SQL := `INSERT INTO user_role(user_id,role) VALUES ($1,$2)`
	_, err := tx.Exec(SQL, userid, role)
	return err
}

func GetIdByPassword(email, password string) (uuid.UUID, error) {
	SQL := `SELECT id, password FROM users where email = TRIM(LOWER($1))`
	var userid uuid.UUID
	var pass string
	err := database.Todo.QueryRowx(SQL, email).Scan(&userid, &pass)
	if err != nil {
		return userid, err
	}
	isMatch := utils.CheckPasswordHash(password, pass)
	if err != nil || !(isMatch) {
		return uuid.Nil, err
	}
	return userid, nil

}

func IsUserRole(userid uuid.UUID) (bool, error) {
	SQL := `SELECT role FROM user_role where user_id = $1`
	var role string
	err := database.Todo.QueryRowx(SQL, userid).Scan(&role)
	if err != nil {
		return false, err
	}
	if role == model.NormalUser {
		return true, nil
	}
	return false, nil
}

func IsSubAdminRole(userid uuid.UUID) (bool, error) {
	SQL := `SELECT role FROM user_role where user_id = $1`
	var role string
	err := database.Todo.QueryRowx(SQL, userid).Scan(&role)
	if err != nil {
		return false, err
	}
	if role == model.SubAdminUser || role == model.AdminUser {
		return true, nil
	}
	return false, nil
}

func IsAdminRole(userid uuid.UUID) (bool, error) {
	SQL := `SELECT role FROM user_role where user_id = $1`
	var role string
	err := database.Todo.QueryRowx(SQL, userid).Scan(&role)
	if err != nil {
		return false, err
	}
	if role == model.AdminUser {
		return true, nil
	}
	return false, nil
}

func CreateRestaurant(userId uuid.UUID, name string, latitude, longitude float64, address string) error {
	SQL := `INSERT INTO resturant(name,latitude,longitude,address,created_by) VALUES($1,$2,$3,$4,$5)`
	_, err := database.Todo.Exec(SQL, name, latitude, longitude, address, userId)
	return err
}

func GetRestaurant(userId uuid.UUID) ([]model.Restaurant, error) {
	var restaurants = make([]model.Restaurant, 0)
	SQL := `SELECT id,name,latitude,longitude,address,created_by FROM resturant where created_by = $1`
	err := database.Todo.Select(&restaurants, SQL, userId)
	return restaurants, err
}

func GetAllRestaurant() ([]model.Restaurant, error) {
	var restaurants = make([]model.Restaurant, 0)
	SQL := `SELECT id,name,latitude,longitude,address,created_by FROM resturant`
	err := database.Todo.Select(&restaurants, SQL)
	return restaurants, err
}

func CreateDish(name string, cost int, createdIn string, createdBy uuid.UUID) error {
	SQL := `INSERT INTO dishes(name,cost,created_in,created_by) VALUES ($1,$2,$3,$4)`
	_, err := database.Todo.Exec(SQL, name, cost, createdIn, createdBy)
	//if err != nil {
	//	return err
	//}
	return err
}

func GetDish() ([]model.Dishes, error) {
	var dishes = make([]model.Dishes, 0)
	SQL := `SELECT r.name,json_agg(json_build_array(d.name, d.cost,u.name))
        from users u
         join resturant r
              on u.id = r.created_by
         join dishes d
              on r.id = d.created_in
              group by  r.name
            order by  r.name`
	err := database.Todo.Select(&dishes, SQL)

	return dishes, err
}

func GetDishes(resturantId uuid.UUID) ([]model.Dishes, error) {
	SQL := `SELECT id,name,cost,created_in,created_by FROM dishes where created_in = $1`
	var dish = make([]model.Dishes, 0)
	err := database.Todo.Select(&dish, SQL, resturantId)

	return dish, err
}

func GetAllSubAdmin() ([]model.Register, error) {
	var allSubAdmin = make([]model.Register, 0)
	SQL := `SELECT u.id,u.name,u.email,u.password 
	FROM users u 
	join user_role ur 
	on u.id = ur.user_id
	where ur.role = 'subAdmin'`
	err := database.Todo.Select(&allSubAdmin, SQL)

	return allSubAdmin, err
}

func CreateAddress(userId uuid.UUID, latitude, longitude float64) error {
	SQL := `INSERT INTO address(user_id,latitude,longitude) VALUES ($1,$2,$3)`
	_, err := database.Todo.Exec(SQL, userId, latitude, longitude)
	return err
}

func GetAddress(userid uuid.UUID) ([]model.Address, error) {
	var address = make([]model.Address, 0)
	SQL := `SELECT id,latitude,longitude FROM address where user_id = $1`
	err := database.Todo.Select(&address, SQL, userid)

	return address, err
}

func FindRestaurantById(userid uuid.UUID) (bool, error) {
	SQL := `SELECT id FROM resturant where id = $1`
	var id uuid.UUID
	err := database.Todo.Get(&id, SQL, userid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return true, nil
}

func CalculateDistance(addressID, restaurantID uuid.UUID) (int64, error) {
	SQL := `SELECT (earth_distance(ll_to_earth(address.latitude, address.longitude),
                                  ll_to_earth(resturant.latitude, resturant.longitude)
                                  )/1000)::integer AS distance_km
			FROM address, resturant
			WHERE address.id = $1 AND resturant.id = $2`
	var distance int64
	err := database.Todo.QueryRowx(SQL, addressID, restaurantID).Scan(&distance)
	return distance, err
}
