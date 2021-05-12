package handlers

import (
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/davisbento/go-api/core/users"
	"github.com/gorilla/mux"
)

func MakeUsersHandlers(r *mux.Router, n *negroni.Negroni, service users.UseCase) {
	r.Handle("/v1/users", n.With(
		negroni.Wrap(getAllUsers(service)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/v1/user/{id}", n.With(
		negroni.Wrap(getUser(service)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/v1/users", n.With(
		negroni.Wrap(storeUser(service)),
	)).Methods("POST", "OPTIONS")
}

func getAllUsers(service users.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getAllUsersJSON(w, service)
	})
}

func getAllUsersJSON(w http.ResponseWriter, service users.UseCase) {
	w.Header().Set("Content-Type", "application/json")
	all, err := service.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(formatJSONError(err.Error()))
		return
	}
	err = json.NewEncoder(w).Encode(all)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(formatJSONError("Erro convertendo em JSON"))
		return
	}
}

func getUser(service users.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(formatJSONError(err.Error()))
			return
		}
		u, err := service.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(formatJSONError(err.Error()))
			return
		}
		err = json.NewEncoder(w).Encode(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(formatJSONError("Erro convertendo em JSON"))
			return
		}
	})
}

/*
Para testar:
curl "POST" "http://localhost:4000/v1/users" \
     -H 'Accept: application/json' \
     -H 'Content-Type: application/json' \
     -d $'{
  "name": "Davi",
  "password": "123456"
}'
*/
func storeUser(service users.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//@TODO este código está duplicado em todos os handlers. Pergunta: como podemos melhorar isso?
		w.Header().Set("Content-Type", "application/json")

		//vamos pegar os dados enviados pelo usuário via body
		var u users.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(formatJSONError(err.Error()))
			return
		}
		//@TODO precisamos validar os dados antes de salvar na base de dados. Pergunta: Como fazer isso?
		created, err := service.Store(&u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(formatJSONError(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)
	})
}
