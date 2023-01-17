package main

import (
	"net/http"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", mailIndexer)
	http.ListenAndServe(":3000", nil)
}

func mailIndexer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content Type", "text/html")
	folderTree := BuildTree("/Users/israelguerrero/Downloads/enron_mail_small")
	bodyDocument := GetBodyDocument(folderTree)
	//tpl.Execute(os.Stdout, bodyDocument)
	tpl.ExecuteTemplate(w, "mailIndexer.html", bodyDocument)
}
