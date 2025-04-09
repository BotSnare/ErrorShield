package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func errorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				renderErrorPage(w, "500.html", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func renderErrorPage(w http.ResponseWriter, templateName string, errorCode int) {
	w.WriteHeader(errorCode)
	tmpl, err := template.ParseFiles("public" + templateName)
	if err != nil {
		http.Error(w, "Error template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/forbidden" {
		renderErrorPage(w, "403.html", http.StatusForbidden)
		return
	}
	tmpl, err := template.ParseFiles("public/index.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	tmpl.Execute(w, nil)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":9090", errorHandler(mux))
}
