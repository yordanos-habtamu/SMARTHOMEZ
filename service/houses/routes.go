package houses

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yordanos-habtamu/realstate/service/auth"
	"github.com/yordanos-habtamu/realstate/types"
	"github.com/yordanos-habtamu/realstate/utils"
)

type contextKey string

const myuser contextKey = "user"

type Handler struct {
	store     types.HouseStore
	userStore types.UserStore
}

func NewHandler(store types.HouseStore, user types.UserStore) *Handler {
	return &Handler{store: store, userStore: user}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/houses/register", auth.WithJwtAuth(h.handleRegisterhouse, h.userStore)).Methods(http.MethodPost)

}

func (h *Handler) handleRegisterhouse(w http.ResponseWriter, r *http.Request) {
	var user = auth.GetUserfromContext(r.Context())
	var house types.RegisterHousePayload
	if user.Role != "admin" {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("Agent and Client can not add houses"))
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Failed to parse the form %v", err))
		return
	}
	numBathrooms, _ := strconv.Atoi(r.FormValue("numBathrooms"))
	numBedrooms, _ := strconv.Atoi(r.FormValue("numBedRooms"))
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	agentID, _ := strconv.ParseUint(r.FormValue("agentId"), 10, 32)
	areaSqFt, _ := strconv.Atoi(r.FormValue("areaSqFt"))
	isSold, _ := strconv.ParseBool(r.FormValue("isSold"))
	house = types.RegisterHousePayload{
		Address:      r.FormValue("address"),
		NumBathrooms: numBathrooms,
		NumBedrooms:  numBedrooms,
		Description:  r.FormValue("description"),
		Price:        price,
		Category:     r.FormValue("category"),
		ImgUrl:       r.FormValue("imgUrl"),
		IsSold:       isSold,
		AgentID:      uint(agentID),
		AreaSqFt:     areaSqFt,
	}
	if err := utils.Validate.Struct(house); err != nil {
		log.Printf("Validation error: %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid house data"))
		return
	}
	err := h.store.CreateHouse(house)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
}
func (h *Handler) handleSearchHouse(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]

    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        log.Printf("error parsing id: %v", err)
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id: %v", err))
        return
    }

    house, err := h.store.GetHouseById(uint(id))
    if err != nil {
        log.Printf("error getting house: %v", err)
        utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("could not fetch house: %v", err))
        return
    }    
    utils.WriteJson(w, http.StatusOK, house)
}

func(h *Handler) handleFetchHouses(w http.ResponseWriter,r *http.Request){
	houses,err := h.store.GetAllHouses()
	if err!= nil{
        log.Printf("something happend %v",err)
		utils.WriteError(w,http.StatusInternalServerError,fmt.Errorf("Erorr fetching houses %v",err))
	}
	if len(houses)==0{
      log.Print("No houses registered")
		utils.WriteError(w,http.StatusNoContent,fmt.Errorf("No houses found"))
	}
	utils.WriteJson(w,http.StatusOK,houses)
}

func (h *Handler) handleUpdateHouse(w http.ResponseWriter, r *http.Request) {
    user := auth.GetUserfromContext(r.Context())
    if user.Role != "admin" {
        log.Printf("You are not authorized to take this action")
        utils.WriteError(w, http.StatusForbidden, fmt.Errorf("forbidden"))
        return
    }

    idStr := r.FormValue("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid id: %v", err))
        return
    }

    var payload types.RegisterHousePayload
    if err :=utils.ParseJson(r,payload); err != nil {
        utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
        return
    }

    house, err := h.store.UpdateHouse(uint(id), payload)
    if err != nil {
        log.Printf("error updating the house: %v", err)
        utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("could not update house: %v", err))
        return
    }

    utils.WriteJson(w, http.StatusOK, house)
}


func(h *Handler) handleDeleteHouse(w http.ResponseWriter,r *http.Request){
	var user = auth.GetUserfromContext(r.Context())
	if user.Role != "admin"{
		log.Printf("You are not authorized to take this action")
		utils.WriteError(w,http.StatusUnauthorized,fmt.Errorf("You are not authorized"))
		return
	}
    vars := mux.Vars(r)
	idStr := vars["id"]
	id,err:= strconv.ParseUint(idStr, 10, 32)
    if err!=nil{
		log.Printf("Error parsing")
		utils.WriteError(w,http.StatusInternalServerError,fmt.Errorf("Error parsing %v",err))
		return
	}
	error:= h.store.DeleteHouse(uint(id))
	if err!=nil{
		log.Printf("Something Wrong Happened %v",error)
		utils.WriteError(w,http.StatusInternalServerError,fmt.Errorf("Something wrong happened %v",error))
		return
	}
	
	utils.WriteJson(w,http.StatusOK,"Successfuly Deleted")
}
func(h *Handler) handleHousesByPrice(w http.ResponseWriter,r *http.Request){
  
	minPrice,_:=strconv.ParseFloat(r.FormValue("minPrice"),64)
    maxPrice,_:=strconv.ParseFloat(r.FormValue("maxPrice"),64)
	houses,err := h.store.GetHousesByPriceRange(minPrice,maxPrice)
	if err!=nil{
		log.Printf("Something Wrong Happened")
		utils.WriteError(w,http.StatusInternalServerError,fmt.Errorf("Something wrong happened %v",err))
		return
	}
	utils.WriteJson(w,http.StatusOK,houses)
}

func (h *Handler) handleLatestHouses(w http.ResponseWriter, r *http.Request){
  
	limit,err := strconv.ParseInt(r.FormValue("limit"),10,64)
     if err!=nil{
		log.Printf("Something Wrong Happened")
		utils.WriteError(w,http.StatusInternalServerError,fmt.Errorf("Something wrong happened %v",err))
		return
	 }
	 houses,err:= h.store.GetLatestHouses(int(limit))
	 utils.WriteJson(w,http.StatusOK,houses)
}

func(h *Handler) handleHouseByCategory(w http.ResponseWriter,r *http.Request){

}