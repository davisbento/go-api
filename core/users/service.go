package users

import (
	"database/sql"

	"github.com/davisbento/go-api/core/utils"
)

type UserCreated struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UseCase interface {
	GetAll() ([]*User, error)
	Get(Id int64) (*User, error)
	Store(u *User) (UserCreated, error)
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
	return nil, nil
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

func NewService(db *sql.DB) *Service {
	return &Service{
		DB: db,
	}
}
