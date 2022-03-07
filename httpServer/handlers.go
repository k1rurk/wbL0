package httpServer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"wb_l0/helperJson"
)

func (s *Server) HandleStatic(w http.ResponseWriter, r *http.Request) {
	render(w, "static/index.html", nil)
}

func (s *Server) HandleSearch(w http.ResponseWriter, r *http.Request) {

	//
	id := r.URL.Query().Get("search_id")
	if id != "" {
		v, err := s.cache.Get(id)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "Заказ не найден")
			return
		}
		str, err := helperJson.PrettyStruct(v)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprint(w, str)
		return
	}
	render(w, "static/index.html", nil)

}

func render(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}
