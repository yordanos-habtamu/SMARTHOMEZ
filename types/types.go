package types
import (
  "time"
)

type UserStore interface{
	GetUserByEmail(email string) (*User,error)
	CreateUser(User) error
	GetUserById(id int) (*User,error)
	GetAllUsers() ([]User,error)
}

type HouseStore interface {
	CreateHouse(RegisterHousePayload) (error)
 	GetHouseByAddress(address string) (*House, error)
	GetHouseById(id uint) (*House, error)
	GetAllHouses() ([]House, error)
	UpdateHouse(id uint, payload RegisterHousePayload) (*House, error)
	DeleteHouse(id uint) error
	GetHousesByCategory(category string) ([]House, error)
	GetHousesByPriceRange(minPrice, maxPrice float64) ([]House, error)
	GetLatestHouses(limit int) ([]House, error)
}




type User struct {
	ID  uint  `json:"Id"`
	FirstName string `json:"firstName"`
	LastName string   `json:"lastName"`
	Sex string `json:"sex"`
	Email string     `json:"email"`
	DoB time.Time   `json:"DoB"`
	Password string  `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	Role   string  `json:"role"`
}
type House struct {
	ID        uint      `json:"id"`
	Address   string    `json:"address"`
	Price     float64   `json:"price"`
	NumBathrooms int    `json:"numBathrooms"`
	NumBedrooms  int       `json:"numRooms"`
	AreaSqFt    int        `json:"area_sq_ft"`
	Category  string    `json:"category"`
	ImageURL  string    `json:"imageUrl"`
	IsSold  bool      `json:"isSold"`
	AgentID   uint      `json:"agentId"`
 	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Description string  `json:"description"`
	ImgURL string   `json:"imgUrl"`
}


type RegisterHousePayload struct {
	Address string `json:"address" validate:"required"`
	NumBathrooms int `json:"numBathrooms" validate:"required"`
	NumBedrooms int `json:"numBedRooms" validate:"required"`
	Description string `json:"description" validate:"required"`
    Price float64 `json:"price" validate:"required"`
	Category string `json:"catagory" validate:"required"`
	ImgUrl string `json:"imgUrl" validate:"required"`
	IsSold bool `json:"isSold"`
	AgentID uint `json:"agentId"`
	AreaSqFt int `json:"areaSqFt" validate:"required"`
	
}

type RegisterUserPayload struct {
	FirstName string  `json:"firstName" validate:"required"`
	LastName  string   `json:"lastName" validate:"required"`
	Email     string    `json:"email" validate:"required,email"` 
	Contact  string  `json:"contact" validate:"required"`
	DoB       string    `json:"DoB" validate:"required"`
	Sex       string      `json:"sex" validate:"required"`
	Password   string      `json:"password" validate:"required,min=6,max=12"`
	Role       string      `json:"role"`
}

type LoginUserPayload struct {
	Email string 	`json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

