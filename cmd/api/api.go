package api

import (
	"database/sql"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/yordanos-habtamu/realstate/service/user"
	"github.com/yordanos-habtamu/realstate/service/houses"
)

type ApiServer struct {
	addr string
	db  *sql.DB
}
func NewApiServer(addr string, db *sql.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Run() error{
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter();
    userStore := user.NewStore(s.db)
    houseStore := houses.NewStore(s.db)
	houseHandler := houses.NewHandler(houseStore,userStore)
	houseHandler.RegisterRoutes(subrouter)
    userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)
	log.Printf("listening on %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}
