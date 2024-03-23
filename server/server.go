package server

import (
	"context"
	"example.com/handlers"
	"example.com/middleware"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"
)

type Server struct {
	*mux.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRoutes() *Server {
	router := mux.NewRouter()
	router.HandleFunc("/register", handlers.RegisterUser).Methods(http.MethodPost)
	router.HandleFunc("/login", handlers.LoginUser).Methods(http.MethodPost)
	router.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "My name is Deepak Singh") }).Methods("GET")
	userRouter := router.PathPrefix("/home").Subrouter()
	userRouter.Use(middleware.JWTMiddleWare)
	userRouter.HandleFunc("/all-restaurant", handlers.GetAllRestaurants).Methods(http.MethodGet)
	//userRouter.HandleFunc("/get-Dishes", handlers.GetDishes).Methods(http.MethodGet)
	userRouter.HandleFunc("/dishes", handlers.GetDishesByRestaurant).Methods(http.MethodGet)
	userRouter.HandleFunc("/address", handlers.CreateAddress).Methods(http.MethodPost)
	userRouter.HandleFunc("/address", handlers.GetAddress).Methods(http.MethodGet)
	userRouter.HandleFunc("/calculate-distance", handlers.GetDistance).Methods(http.MethodGet)
	userRouter.HandleFunc("/logout", handlers.Logout).Methods(http.MethodGet)

	subAdminRouter := userRouter.PathPrefix("/subAdmin").Subrouter()
	subAdminRouter.Use(middleware.SubAdminMiddleware)
	subAdminRouter.HandleFunc("/restaurant", handlers.CreateRestaurant).Methods(http.MethodPost)
	subAdminRouter.HandleFunc("/dish", handlers.CreateDishes).Methods(http.MethodPost)
	subAdminRouter.HandleFunc("/restaurant", handlers.GetRestaurant).Methods(http.MethodGet)

	adminRouter := userRouter.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AdminMiddleware)
	adminRouter.HandleFunc("/subAdmin", handlers.CreateSubAdmin).Methods(http.MethodPost)
	adminRouter.HandleFunc("/all-subAdmin", handlers.GetAllSubAdmin).Methods(http.MethodGet)

	return &Server{
		Router: router,
	}
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}

func (svc *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return svc.server.Shutdown(ctx)
}
