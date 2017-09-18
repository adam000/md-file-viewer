package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.String(), "/")
	log.Printf("Routing request '/%s'", path)
	if path == "" {
		// Defer to directory handler for base.
		dirHandler(w, r, path)
		return
	}

	if fileInfo, err := os.Stat(path); os.IsNotExist(err) {
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if fileInfo.IsDir() {
			// Defer to directory handler.
			dirHandler(w, r, path)
			return
		}
	}

	path += ".md"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Serving up %s", path)

	markdown, err := ioutil.ReadFile(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	htmlResult := template.HTML(bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownCommon(markdown)))

	cssBytes, err := ioutil.ReadFile(filepath.Join(cfg.StyleDir, "file.css"))
	if err != nil {
		log.Println(err)
	}

	data := struct {
		Css      template.CSS
		Content  template.HTML
		FileName string
	}{
		Css:      template.CSS(cssBytes),
		Content:  htmlResult,
		FileName: path,
	}

	tpl.ExecuteTemplate(w, "page_file.html", data)
}
