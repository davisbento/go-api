package users

import (
	"database/sql"
)

type UseCase interface {
	GetAll() ([]*User, error)
	Get(Id int64) (*User, error)
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

func NewService(db *sql.DB) *Service {
	return &Service{
		DB: db,
	}
}
