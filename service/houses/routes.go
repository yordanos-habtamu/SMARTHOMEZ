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
	router.HandleFunc("/houses/delete/{id}", auth.WithJwtAuth(h.handleDeleteHouse, h.userStore)).Methods(http.MethodDelete)
	router.HandleFunc("/houses/update", auth.WithJwtAuth(h.handleUpdateHouse, h.userStore)).Methods(http.MethodDelete)
	router.HandleFunc("/houses/id/{id}", h.handleSearchHouse).Methods(http.MethodGet)
	router.HandleFunc("/houses/category/{category}", h.handleHouseByCategory).Methods(http.MethodGet)
	router.HandleFunc("/houses/address/{address}", h.handleHouseByAddress).Methods(http.MethodGet)
	router.HandleFunc("/houses/price", h.handleHousesByPrice).Methods(http.MethodGet)
	router.HandleFunc("/houses/latest", h.handleLatestHouses).Methods(http.MethodGet)
	router.HandleFunc("/houses", h.handleFetchHouses).Methods(http.MethodGet)

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

	// Parse form fields
	numBathrooms, err := strconv.Atoi(r.FormValue("numBathrooms"))
	if err != nil {
		log.Printf("error parsing %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse the form %v", err))
		return
	}
	numBedrooms, err := strconv.Atoi(r.FormValue("numBedRooms"))
	if err != nil {
		log.Printf("error parsing %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse the form %v", err))
		return
	}
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		log.Printf("error parsing %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse the form %v", err))
		return
	}
	agentID, err := strconv.ParseUint(r.FormValue("agentId"), 10, 32)
	if err != nil {
		log.Printf("error parsing %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse the form %v", err))
		return
	}
	areaSqFt, err := strconv.Atoi(r.FormValue("areaSqFt"))
	if err != nil {
		log.Printf("error parsing %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse the form %v", err))
		return
	}
	isSold, err := strconv.ParseBool(r.FormValue("isSold"))
	if err != nil {
		log.Printf("error parsing %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse the form %v", err))
		return
	}

	// Handle image upload
	var imgUrl string
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("No image file provided or error: %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("image file is required"))
		return
	}
	defer file.Close()

	// Upload to Cloudinary
	imgUrl, err = utils.UploadHouseImage(file, header.Filename)
	if err != nil {
		log.Printf("Failed to upload image: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to upload image"))
		return
	}

	house = types.RegisterHousePayload{
		Address:      r.FormValue("address"),
		NumBathrooms: numBathrooms,
		NumBedrooms:  numBedrooms,
		Description:  r.FormValue("description"),
		Price:        price,
		Category:     r.FormValue("category"),
		ImgUrl:       imgUrl, // Use Cloudinary URL
		IsSold:       isSold,
		AgentID:      uint(agentID),
		AreaSqFt:     areaSqFt,
	}

	// Convert payload to House struct
	dbHouse := utils.MapPayloadToHouse(house)

	if err := utils.Validate.Struct(house); err != nil {
		log.Printf("Validation error: %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid house data"))
		return
	}
	newerr := h.store.CreateHouse(dbHouse)
	if newerr != nil {
		utils.WriteError(w, http.StatusBadRequest, newerr)
		return
	}

	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "House registered successfully", "imageUrl": imgUrl})
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

func (h *Handler) handleFetchHouses(w http.ResponseWriter, r *http.Request) {
	houses, err := h.store.GetAllHouses()
	if err != nil {
		log.Printf("something happend %v", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("erorr fetching houses %v", err))
	}
	if len(houses) == 0 {
		log.Print("No houses registered")
		utils.WriteError(w, http.StatusNoContent, fmt.Errorf("no houses found"))
	}
	utils.WriteJson(w, http.StatusOK, houses)
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
	if err := utils.ParseJson(r, payload); err != nil {
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

func (h *Handler) handleDeleteHouse(w http.ResponseWriter, r *http.Request) {
	var user = auth.GetUserfromContext(r.Context())
	if user.Role != "admin" {
		log.Printf("You are not authorized to take this action")
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("you are not authorized"))
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("error parsing")
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error parsing %v", err))
		return
	}
	error := h.store.DeleteHouse(uint(id))
	if error != nil {
		log.Printf("Something Wrong Happened %v", error)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happened %v", error))
		return
	}

	utils.WriteJson(w, http.StatusOK, "Successfuly Deleted")
}
func (h *Handler) handleHousesByPrice(w http.ResponseWriter, r *http.Request) {

	minPrice, _ := strconv.ParseFloat(r.FormValue("minPrice"), 64)
	maxPrice, _ := strconv.ParseFloat(r.FormValue("maxPrice"), 64)
	houses, err := h.store.GetHousesByPriceRange(minPrice, maxPrice)
	if err != nil {
		log.Printf("Something Wrong Happened")
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happened %v", err))
		return
	}
	utils.WriteJson(w, http.StatusOK, houses)
}

func (h *Handler) handleLatestHouses(w http.ResponseWriter, r *http.Request) {

	limit, err := strconv.ParseInt(r.FormValue("limit"), 10, 64)
	if err != nil {
		log.Printf("Something Wrong Happened")
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happened %v", err))
		return
	}
	houses, err := h.store.GetLatestHouses(int(limit))
	if err != nil {
		log.Printf("Something Wrong Happened %v", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something wrong happened %v", err))
		return
	}
	utils.WriteJson(w, http.StatusOK, houses)
}

func (h *Handler) handleHouseByCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	category := vars["category"]
	houses, err := h.store.GetHousesByCategory(category)
	if err != nil {
		log.Printf("Something wrong %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, fmt.Errorf("something happened %v", err))
		return
	}
	utils.WriteJson(w, http.StatusOK, houses)

}
func (h *Handler) handleHouseByAddress(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	address := vars["address"]
	house, err := h.store.GetHouseByAddress(address)
	if err != nil {
		log.Printf("Something wrong %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, fmt.Errorf("something happened %v", err))
		return
	}
	utils.WriteJson(w, http.StatusOK, house)

}
