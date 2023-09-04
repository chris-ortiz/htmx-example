package main

import (
	"html/template"
	"htmx_sample/todoitem"
	"net/http"
	"strconv"
	"strings"
)

var page *template.Template

const tmpl = "index.tmpl"

var itemStore = todoitem.NewItemStore()

func init() {
	page = template.New(tmpl)
	page = template.Must(page.ParseFiles(tmpl))
}

func handleRootRequest(w http.ResponseWriter, r *http.Request) {
	items := itemStore.FindAll()

	if err := page.ExecuteTemplate(w, tmpl, items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleItemCreateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid http method", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err == nil {
		item := todoitem.TodoItem{Text: r.FormValue("input-todo")}
		itemStore.Add(&item)

		if err := page.ExecuteTemplate(w, "item", item); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleItemDeleteRequest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		itemStore.Delete(id)
	}
}

func main() {

	http.HandleFunc("/", handleRootRequest)
	http.HandleFunc("/item", handleItemCreateRequest)
	http.HandleFunc("/item/", handleItemDeleteRequest)
	http.ListenAndServe(":80", nil)
}
