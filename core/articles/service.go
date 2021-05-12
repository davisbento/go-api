package articles

import (
	"database/sql"
)

type UseCase interface {
	GetAll() ([]*Article, error)
	Get(Id int64) (*Article, error)
}

//a struct Service agora tem uma conexão com o banco de dados dentro dela
type Service struct {
	DB *sql.DB
}

//vamos implementar as funções na próxima etapa
func (s *Service) GetAll() ([]*Article, error) {
	//result é um slice de ponteiros do tipo Beer
	var result []*Article

	//vamos sempre usar a conexão que está dentro do Service
	rows, err := s.DB.Query("select id, title, description, content from articles")
	//se existe erro a função deve retorná-lo e ele vai ser tratado
	//por quem chamou o pacote. Essa é uma boa prática em Go
	if err != nil {
		return nil, err
	}
	//a função defer garante que o comando rows.Close vai ser executado na saída da função
	//desta forma não precisamos nos preocupar em fechar a conexão
	defer rows.Close()
	for rows.Next() {
		var a Article
		err = rows.Scan(&a.Id, &a.Title, &a.Description, &a.Content)
		if err != nil {
			return nil, err
		}
		//o comando append adiciona novos itens a um slice, sempre no final
		result = append(result, &a)
	}
	return result, nil
}

func (s *Service) Get(id int64) (*Article, error) {
	return nil, nil
}

//esta função retorna um ponteiro em memória para uma estrutura
//a função agora recebe uma conexão com o banco de dados
func NewService(db *sql.DB) *Service {
	return &Service{
		DB: db,
	}
}
