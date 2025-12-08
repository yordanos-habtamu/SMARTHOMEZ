package houses

import (
	"database/sql"
	"fmt"

	"github.com/yordanos-habtamu/realstate/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetHouseByAddress(address string) (*types.House, error) {
	h := new(types.House)
	row := s.db.QueryRow("SELECT id,category,address,price,num_bedrooms,num_bathrooms,area_sq_ft,description,img_url,created_at,updated_at,is_sold,agent_id FROM house WHERE address = $1", address)

	err := row.Scan(
		&h.ID,
		&h.Category,
		&h.Address,
		&h.Price,
		&h.NumBedrooms,
		&h.NumBathrooms,
		&h.AreaSqFt,
		&h.Description,
		&h.ImgURL,
		&h.CreatedAt,
		&h.UpdatedAt,
		&h.IsSold,
		&h.AgentID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("house not found")
		}
		return nil, err
	}

	return h, nil
}

func (s *Store) GetHousesByCategory(category string) ([]types.House, error) {
	rows, err := s.db.Query("SELECT id,category,address,price,num_bedrooms,num_bathrooms,area_sq_ft,description,img_url,created_at,updated_at,is_sold,agent_id FROM house WHERE category = $1", category)
	if err != nil {
		return nil, fmt.Errorf("error fetching houses: %v", err)
	}
	defer rows.Close()

	houses := []types.House{}

	for rows.Next() {
		h := types.House{}
		err := rows.Scan(
			&h.ID,
			&h.Category,
			&h.Address,
			&h.Price,
			&h.NumBedrooms,
			&h.NumBathrooms,
			&h.AreaSqFt,
			&h.Description,
			&h.ImgURL,
			&h.CreatedAt,
			&h.UpdatedAt,
			&h.IsSold,
			&h.AgentID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning house: %v", err)
		}
		houses = append(houses, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	if len(houses) == 0 {
		return nil, fmt.Errorf("no houses found")
	}

	return houses, nil
}

func (s *Store) GetAllHouses() ([]types.House, error) {
	rows, err := s.db.Query("SELECT id,category,address,price,num_bedrooms,num_bathrooms,area_sq_ft,description,img_url,created_at,updated_at,is_sold,agent_id FROM house")
	if err != nil {
		return nil, fmt.Errorf("error fetching houses: %v", err)
	}
	defer rows.Close()

	houses := []types.House{}

	for rows.Next() {
		h := types.House{}
		err := rows.Scan(
			&h.ID,
			&h.Category,
			&h.Address,
			&h.Price,
			&h.NumBedrooms,
			&h.NumBathrooms,
			&h.AreaSqFt,
			&h.Description,
			&h.ImgURL,
			&h.CreatedAt,
			&h.UpdatedAt,
			&h.IsSold,
			&h.AgentID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning house: %v", err)
		}
		houses = append(houses, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	if len(houses) == 0 {
		return nil, fmt.Errorf("no houses found")
	}

	return houses, nil
}

func (s *Store) GetHouseById(id uint) (*types.House, error) {
	rows := s.db.QueryRow("SELECT id,category,address,price,num_bedrooms,num_bathrooms,area_sq_ft,description,img_url,created_at,updated_at,is_sold,agent_id FROM house WHERE id = $1", id)
	h := new(types.House)
	err := rows.Scan(
		&h.ID,
		&h.Category,
		&h.Address,
		&h.Price,
		&h.NumBedrooms,
		&h.NumBathrooms,
		&h.AreaSqFt,
		&h.Description,
		&h.ImgURL,
		&h.CreatedAt,
		&h.UpdatedAt,
		&h.IsSold,
		&h.AgentID,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching house: %v", err)
	}
	if h.ID == 0 {
		return nil, fmt.Errorf("house not found")
	}
	return h, nil
}

func (s *Store) CreateHouse(h types.House) error {
	_, err := s.db.Exec("INSERT INTO house (category,address,price,num_bedrooms,num_bathrooms,area_sq_ft,description,img_url,is_sold,agent_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)",
		h.Category,
		h.Address,
		h.Price,
		h.NumBedrooms,
		h.NumBathrooms,
		h.AreaSqFt,
		h.Description,
		h.ImgURL,
		h.IsSold,
		h.AgentID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateHouse(id uint, payload types.RegisterHousePayload) (*types.House, error) {
	i, err := s.GetHouseById(id)
	if err != nil {
		return nil, fmt.Errorf("Something happened %v", err)
	}
	if i == nil {
		return nil, fmt.Errorf("house is not found")
	}
	var updatedHouse types.House

	er := s.db.QueryRow(`
    UPDATE house
    SET address = $1, price = $2, num_bedrooms = $3, num_bathrooms = $4,
        area_sq_ft = $5, description = $6, img_url = $7, is_sold = $8,
        agent_id = $9, updated_at = NOW()
    WHERE id = $10
    RETURNING id, category, address, price, num_bedrooms, num_bathrooms,
              area_sq_ft, description, img_url, created_at, updated_at,
              is_sold, agent_id
	`, payload.Address, payload.Price, payload.NumBedrooms, payload.NumBathrooms, payload.AreaSqFt, payload.Description, payload.ImgUrl, payload.IsSold, payload.AgentID, id).Scan(
		&updatedHouse.ID,
		&updatedHouse.Category,
		&updatedHouse.Address,
		&updatedHouse.Price,
		&updatedHouse.NumBedrooms,
		&updatedHouse.NumBathrooms,
		&updatedHouse.AreaSqFt,
		&updatedHouse.Description,
		&updatedHouse.ImgURL,
		&updatedHouse.CreatedAt,
		&updatedHouse.UpdatedAt,
		&updatedHouse.IsSold,
		&updatedHouse.AgentID,
	)
	if er != nil {
		return nil, er
	}
	return &updatedHouse, nil
}

func (s *Store) DeleteHouse(id uint) error {
	_, err := s.db.Exec("DELETE FROM house WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting house: %v", err)
	}
	return nil
}

func (s *Store) GetHousesByPriceRange(minPrice, maxPrice float64) ([]types.House, error) {
	rows, err := s.db.Query("SELECT id,category,address,price,num_bedrooms,num_bathrooms,area_sq_ft,description,img_url,created_at,updated_at,is_sold,agent_id FROM house WHERE price >= $1 AND price <= $2", minPrice, maxPrice)
	if err != nil {
		return nil, fmt.Errorf("error fetching houses: %v", err)
	}
	defer rows.Close()

	houses := []types.House{}

	for rows.Next() {
		h := types.House{}
		err := rows.Scan(
			&h.ID,
			&h.Category,
			&h.Address,
			&h.Price,
			&h.NumBedrooms,
			&h.NumBathrooms,
			&h.AreaSqFt,
			&h.Description,
			&h.ImgURL,
			&h.CreatedAt,
			&h.UpdatedAt,
			&h.IsSold,
			&h.AgentID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning house: %v", err)
		}
		houses = append(houses, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	if len(houses) == 0 {
		return nil, fmt.Errorf("no houses found")
	}

	return houses, nil
}

func (s *Store) GetLatestHouses(limit int) ([]types.House, error) {
	query := "SELECT id,category,address,price,num_bedrooms,num_bathrooms,area_sq_ft,description,img_url,created_at,updated_at,is_sold,agent_id FROM house ORDER BY created_at DESC LIMIT $1"
	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching houses: %v", err)
	}
	defer rows.Close()

	houses := []types.House{}

	for rows.Next() {
		h := types.House{}
		err := rows.Scan(
			&h.ID,
			&h.Category,
			&h.Address,
			&h.Price,
			&h.NumBedrooms,
			&h.NumBathrooms,
			&h.AreaSqFt,
			&h.Description,
			&h.ImgURL,
			&h.CreatedAt,
			&h.UpdatedAt,
			&h.IsSold,
			&h.AgentID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning house: %v", err)
		}
		houses = append(houses, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	if len(houses) == 0 {
		return nil, fmt.Errorf("no houses found")
	}

	return houses, nil
}
