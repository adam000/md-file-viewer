package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	requestBase := strings.TrimLeft(r.URL.String(), "/")
	diskPath := path.Join(cfg.RootDir, requestBase)
	log.Printf("Routing request '/%s' using data at '%s'", requestBase, diskPath)
	if requestBase == "" {
		// Defer to directory handler for base.
		dirHandler(w, r, requestBase, diskPath)
		return
	}

	if fileInfo, err := os.Stat(diskPath); os.IsNotExist(err) {
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if fileInfo.IsDir() {
		// Defer to directory handler.
		dirHandler(w, r, requestBase, diskPath)
		return
	}

	if _, err := os.Stat(diskPath + ".md"); os.IsNotExist(err) {
		otherFileHandler(diskPath, w, r)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	diskPath += ".md"
	log.Printf("Serving up %s", diskPath)

	markdown, err := ioutil.ReadFile(diskPath)
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
		FileName: requestBase,
	}

	tpl.ExecuteTemplate(w, "page_file.html", data)
}

func otherFileHandler(path string, w http.ResponseWriter, r *http.Request) {
	if !isImagePath(path) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filepath.Join(cfg.RootDir, path))
}

func isImagePath(path string) bool {
	isImage := false
	extension := strings.ToLower(filepath.Ext(path))

	for _, imageExt := range []string{".jpg", ".png", ".jpeg", ".gif"} {
		if extension == imageExt {
			isImage = true
			break
		}
	}

	return isImage
}
