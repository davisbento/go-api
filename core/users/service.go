package users

import (
	"database/sql"
	"fmt"

	"github.com/davisbento/go-api/core/utils"
)

type UserCreated struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	AuthToken string `json:"authToken"`
}

type UseCase interface {
	GetAll() ([]*User, error)
	Get(Id int64) (*User, error)
	getByEmail(email string) (*User, error)
	Store(u *User) (UserCreated, error)
	Login(u *UserLogin) (LoginResponse, error)
}

type Service struct {
	DB *sql.DB
}

func (s *Service) GetAll() ([]*User, error) {
	var result []*User

	rows, err := s.DB.Query("select id, name from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a User
		err = rows.Scan(&a.Id, &a.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, &a)
	}
	return result, nil
}

func (s *Service) Get(id int64) (*User, error) {
	var u User

	stmt, err := s.DB.Prepare("select id, name, email from users where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&u.Id, &u.Name, &u.Email)
	if err != nil {
		return nil, err
	}
	//deve retornar a posição da memória de b
	return &u, nil
}

func (s *Service) Store(u *User) (UserCreated, error) {
	hashedPassword, err := utils.HashPassword(u.Password)

	user := UserCreated{}

	if err != nil {
		panic(err)
	}

	//iniciamos uma transação
	tx, err := s.DB.Begin()
	if err != nil {
		return user, err
	}

	sqlStatement := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)`

	_, err = s.DB.Exec(sqlStatement, u.Name, u.Email, hashedPassword)
	if err != nil {
		tx.Rollback()
		return user, err
	}

	tx.Commit()

	user.Name = u.Name
	user.Email = u.Email
	return user, nil
}

func (s *Service) getByEmail(email string) (*User, error) {
	var u User

	sqlStatement := `SELECT id, name, email, password FROM users WHERE email=$1`
	row := s.DB.QueryRow(sqlStatement, email)
	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *Service) Login(u *UserLogin) (LoginResponse, error) {
	response := LoginResponse{}

	user, err := s.getByEmail(u.Email)
	if err != nil {
		return response, err
	}

	isValid := utils.ComparePasswords(user.Password, u.Password)

	if !isValid {
		return response, fmt.Errorf("password-invalid")
	}

	response.AuthToken = "3123123jdajaja"

	return response, nil
}

func NewService(db *sql.DB) *Service {
	return &Service{
		DB: db,
	}
}
