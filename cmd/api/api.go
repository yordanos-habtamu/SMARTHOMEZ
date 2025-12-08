package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yordanos-habtamu/realstate/service/houses"
	"github.com/yordanos-habtamu/realstate/service/referral"
	"github.com/yordanos-habtamu/realstate/service/user"
	"github.com/yordanos-habtamu/realstate/utils"
)

type ApiServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// Initialize Cloudinary
	if err := utils.InitCloudinary(); err != nil {
		log.Printf("Warning: Cloudinary not initialized: %v", err)
	}

	userStore := user.NewStore(s.db)
	houseStore := houses.NewStore(s.db)
	referralStore := referral.NewStore(s.db)

	houseHandler := houses.NewHandler(houseStore, userStore)
	houseHandler.RegisterRoutes(subrouter)

	userHandler := user.NewHandler(userStore, referralStore)
	userHandler.RegisterRoutes(subrouter)

	referralHandler := referral.NewHandler(referralStore, userStore)
	referralHandler.RegisterRoutes(subrouter)

	log.Printf("listening on %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}
