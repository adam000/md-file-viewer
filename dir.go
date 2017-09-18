package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type fileOrDir struct {
	Name       string
	FullPath   string
	Parent     string
	IsDir      bool
	FilesUnder dirListing
}

type dirListing []fileOrDir

func dirHandler(w http.ResponseWriter, r *http.Request, path string) {
	path = filepath.Join(cfg.RootDir, path)

	data := struct {
		Css       template.CSS
		Directory fileOrDir
	}{
		Directory: fileOrDir{
			Name:     filepath.Base(path),
			FullPath: path,
			Parent:   "/" + filepath.Dir(path),
			IsDir:    true,
		},
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			dir := fileOrDir{
				Name:     file.Name(),
				FullPath: "/" + filepath.Join(path, file.Name()),
				IsDir:    true,
			}
			data.Directory.FilesUnder = append(data.Directory.FilesUnder, dir)
			// TODO: recurse 1-2 layers?
		} else if strings.HasSuffix(file.Name(), ".md") {
			// TODO maybe get rid of the million file.Name() calls?
			dir := fileOrDir{
				Name:     strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())),
				FullPath: "/" + strings.TrimSuffix(filepath.Join(path, file.Name()), filepath.Ext(file.Name())),
				IsDir:    false,
			}
			data.Directory.FilesUnder = append(data.Directory.FilesUnder, dir)
		} else {
			// TODO different kind of handler for other files?
			log.Printf("Skipping non-md file %s", file.Name())
		}
	}

	tpl.ExecuteTemplate(w, "page_dir.html", data)

}
