package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var tpl *template.Template

const wrapper = "<html><head></head><body>{{.}}</body></html>"

func main() {
	tpl, err := template.New("T").Parse(wrapper)
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	r := mux.NewRouter()

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimLeft(r.URL.String(), "/") + ".md"
		log.Print(path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		markdown, err := ioutil.ReadFile(path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		htmlResult := template.HTML(bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownCommon(markdown)))
		tpl.Execute(w, htmlResult)
	})

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
