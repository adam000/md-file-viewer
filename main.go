package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var tpl *template.Template
var cfg configuration

func main() {
	// Load config
	c, err := loadConfiguration()
	cfg = c

	if err != nil {
		fmt.Println("Error in loading config:", err)
		fmt.Println("Default config will be used.")
	}

	// Validate config
	if err := c.Validate(); err != nil {
		fmt.Println("Error validating config:", err)
		return
	}

	// Parse templates
	tpl = template.New("T")
	templateGlob := fmt.Sprintf("%s/*.html", c.TemplateDir)
	log.Printf("Parsing templates in '%s'", templateGlob)
	template.Must(tpl.ParseGlob(templateGlob))

	// Create router
	r := mux.NewRouter()

	r.PathPrefix("/").HandlerFunc(fileHandler)
	http.Handle("/", r)

	port := ":6060"
	log.Printf("Server running at localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
