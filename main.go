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
		tpl.Execute(w, struct {
			Css      template.CSS
			Content  template.HTML
			FileName string
		}{
			Css:      template.CSS(css),
			Content:  htmlResult,
			FileName: path,
		})
	})

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

const wrapper = `
<html>
<head>
  <style type="text/css">
    {{.Css}}
  </style>
</head>
<body>
	<div id="main">
		<h1 class="file-name">{{.FileName}}</h1>
		<div id="content">
			{{.Content}}
		</div>
	</div>
</body>
</html>
`

const css = `
#main {
	width: 800px;
	margin: 0 auto;
	font-family: sans-serif;
}

.file-name {
	margin-bottom: 30px;
}

#content {
	border: 1px #DDD solid;
}

code {
	font-family: Inconsolata, monospace;
	padding: 0 3px;
}

pre {
	background: #DDD;
	border: 1px #888 solid;
	margin: 5px;
	padding: 5px;
}
`
