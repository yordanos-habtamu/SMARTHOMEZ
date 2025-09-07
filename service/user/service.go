package user

import (
	"database/sql"
	"fmt"
	
	"github.com/yordanos-habtamu/realstate/types"
)

type Store struct {
  db *sql.DB
}

func NewStore (db *sql.DB) *Store {
  return &Store{db:db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
    u := new(types.User)
    row := s.db.QueryRow("SELECT id, first_name, last_name, sex, email, dob, password, created_at, role FROM users WHERE email = $1", email)

    err := row.Scan(
        &u.ID,
        &u.FirstName,
        &u.LastName,
        &u.Sex,
        &u.Email,
        &u.DoB,
        &u.Password,
        &u.CreatedAt,
        &u.Role,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }

     return u, nil
}


func (s *Store) GetAllUsers() ([]types.User, error) {
    rows, err := s.db.Query("SELECT id, first_name, last_name, sex, email, dob, password, created_at, role FROM users")
    if err != nil {
        return nil, fmt.Errorf("error fetching users: %v", err)
    }
    defer rows.Close()

    users := []types.User{}

    for rows.Next() {
        u := types.User{}
        err := rows.Scan(
            &u.ID,
            &u.FirstName,
            &u.LastName,
            &u.Sex,
            &u.Email,
            &u.DoB,
            &u.Password,
            &u.CreatedAt,
            &u.Role,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning user: %v", err)
        }
        users = append(users, u)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %v", err)
    }

    if len(users) == 0 {
        return nil, fmt.Errorf("no users found")
    }

    return users, nil
}

func (s *Store) GetUserById(id int) (*types.User,error){
  rows := s.db.QueryRow("SELECT * FROM users WHERE id = $1",id)
  u := new(types.User)
  err := rows.Scan(
	    &u.ID,
        &u.FirstName,
        &u.LastName,
        &u.Sex,
        &u.Email,
        &u.DoB,
        &u.Password,
        &u.CreatedAt,
        &u.Role,
  );
  if err!= nil {
	return nil,fmt.Errorf("error fetching user: %v",err)
  }
  if u.ID == 0{
    return nil, fmt.Errorf("user not found")
  }
  return u,nil
}

func (s *Store) CreateUser(user types.User) error {
  _,err := s.db.Exec("INSERT INTO users (first_name,last_name,email,password,dob,contact,sex,role) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",user.FirstName,user.LastName,user.Email,user.Password,user.DoB,user.Sex,user.Role)
  if err != nil {
    return err
  }
  return nil
}