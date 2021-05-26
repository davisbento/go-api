package handlers

import (
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/davisbento/go-api/core/articles"
	"github.com/davisbento/go-api/core/utils"
	"github.com/gorilla/mux"
)

func MakeArticlesHandler(r *mux.Router, n *negroni.Negroni, service articles.UseCase) {
	r.Handle("/v1/articles", n.With(
		negroni.Wrap(getAllArticles(service)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/v1/article/{id}", n.With(
		negroni.Wrap(getArticle(service)),
	)).Methods("GET", "OPTIONS")

}

/*
Para testar:
curl ' http://localhost:4000/v1/articles
*/
func getAllArticles(service articles.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getAllArticlesJSON(w, service)
	})
}

func getAllArticlesJSON(w http.ResponseWriter, service articles.UseCase) {
	w.Header().Set("Content-Type", "application/json")
	all, err := service.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.FormatJSONError(err.Error()))
		return
	}
	//vamos converter o resultado em JSON e gerar a resposta
	err = json.NewEncoder(w).Encode(all)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.FormatJSONError("Erro convertendo em JSON"))
		return
	}
}

/*
Para testar:
curl http://localhost:4000/v1/article/1
*/
func getArticle(service articles.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//@TODO este código está duplicado em todos os handlers. Pergunta: como podemos melhorar isso?
		w.Header().Set("Content-Type", "application/json")

		//vamos pegar o ID da URL
		//na definição do protocolo http, os parâmetros são enviados no formato de texto
		//por isso precisamos converter em int64
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.FormatJSONError(err.Error()))
			return
		}
		b, err := service.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(utils.FormatJSONError(err.Error()))
			return
		}
		//vamos converter o resultado em JSON e gerar a resposta
		err = json.NewEncoder(w).Encode(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.FormatJSONError("Erro convertendo em JSON"))
			return
		}
	})
}
