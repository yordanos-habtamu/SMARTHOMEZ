package houses

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/mux"
	"github.com/yordanos-habtamu/realstate/service/auth"
	"github.com/yordanos-habtamu/realstate/types"
	"github.com/yordanos-habtamu/realstate/utils"
)

type contextKey string
const myuser contextKey = "user"

type Handler struct {
	store types.HouseStore
	userStore types.UserStore
}

func NewHandler (store types.HouseStore,user types.UserStore ) *Handler{
	return &Handler{store:store,userStore:user}
}

func (h *Handler) RegisterRoutes (router *mux.Router){
	router.HandleFunc("/houses/register",auth.WithJwtAuth(h.handleRegisterhouse,h.userStore)).Methods(http.MethodPost)
}

func (h *Handler ) handleRegisterhouse (w http.ResponseWriter, r *http.Request){
	var user = auth.GetUserfromContext(r.Context())
	var house types.RegisterHousePayload
	if user.Role != "admin"{
	    utils.WriteError(w,http.StatusForbidden,fmt.Errorf("Agent and Client can not add houses"))
		return
 	} 
	 if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w,http.StatusBadRequest,fmt.Errorf("Failed to parse the form %v",err))
		return
	}
	numBathrooms, _ := strconv.Atoi(r.FormValue("numBathrooms"))
   numBedrooms, _ := strconv.Atoi(r.FormValue("numBedRooms"))
   price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
  agentID, _ := strconv.ParseUint(r.FormValue("agentId"), 10, 32)
  areaSqFt, _ := strconv.Atoi(r.FormValue("areaSqFt"))
  isSold, _ := strconv.ParseBool(r.FormValue("isSold"))
	house = types.RegisterHousePayload{
      Address : r.FormValue("address"),
	NumBathrooms:  numBathrooms,
	NumBedrooms:numBedrooms,
	Description: r.FormValue("description"), 
    Price :price,
	Category: r.FormValue("category"),
	ImgUrl :r.FormValue("imgUrl"),
	IsSold:isSold,
	AgentID:uint(agentID),
	AreaSqFt:areaSqFt,
	}
	if err:= utils.Validate.Struct(house); err!=nil{
		log.Printf("Validation error: %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid house data"))
		return
	}
	 err := h.store.CreateHouse(house)
	if err!=nil{
		utils.WriteError(w,http.StatusBadRequest,err)
		return
	}
}
func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request){
	var userID uint = auth.GetUserIdfromContext(r.Context())
	var cart types.CartCheckoutPayload
	if err := utils.ParseJson(r,&cart); err !=nil {
	utils.WriteError(w,http.StatusBadRequest,err)
	return
	}
	if err := utils.Validate.Struct(cart); err != nil {
		log.Printf("Validation error: %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	productIDs, err:= getCartItemsIDs(cart.Items)
	if err != nil {
		utils.WriteError(w,http.StatusBadRequest,err)
	}
   ps,err := h.productStore.GetProductByIds(productIDs)
   if err != nil {
	utils.WriteError(w,http.StatusInternalServerError,err)
   }
   orderId, totalPrice,err := h.CreateOrder(ps,cart.Items,int(userID))
   if err != nil {
	utils.WriteError(w,http.StatusInternalServerError,err)
   }
   utils.WriteJson(w,http.StatusOK,map[string] any{
	   "order_id":orderId,
	   "total_price":totalPrice,
   })
}  