package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
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

func dirHandler(w http.ResponseWriter, r *http.Request, dirPath, diskPath string) {
	thisDirPath := path.Base(dirPath)
	if thisDirPath == "." {
		thisDirPath = "Root"
	}
	data := struct {
		Css       template.CSS
		Directory fileOrDir
	}{
		Css: ".fa { width: 20px; }",
		Directory: fileOrDir{
			Name:     thisDirPath,
			FullPath: dirPath,
			Parent:   "/" + path.Dir(dirPath),
			IsDir:    true,
		},
	}

	files, err := ioutil.ReadDir(diskPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error reading directory: %s", err)
		return
	}
	for _, file := range files {
		fileName := file.Name()
		if file.IsDir() {
			dir := fileOrDir{
				Name:     fileName,
				FullPath: "/" + path.Join(dirPath, fileName),
				IsDir:    true,
			}
			data.Directory.FilesUnder = append(data.Directory.FilesUnder, dir)
			// TODO: recurse 1-2 layers?
		} else if strings.HasSuffix(fileName, ".md") {
			fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			dir := fileOrDir{
				Name:     fileNameWithoutExt,
				FullPath: "/" + path.Join(dirPath, fileNameWithoutExt),
				IsDir:    false,
			}
			data.Directory.FilesUnder = append(data.Directory.FilesUnder, dir)
		} else if isImagePath(filepath.Join(dirPath, file.Name())) {
			dir := fileOrDir{
				Name:     file.Name(),
				FullPath: "/" + path.Join(dirPath, file.Name()),
				IsDir:    false,
			}
			data.Directory.FilesUnder = append(data.Directory.FilesUnder, dir)
		} else {
			// TODO different kind of handler for other files?
			log.Printf("Skipping non-md/image file %s", fileName)
		}
	}

	tpl.ExecuteTemplate(w, "page_dir.html", data)
}
