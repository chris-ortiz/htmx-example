package main

import (
	"html/template"
	"htmx_sample/todoitem"
	"net/http"
)

var page *template.Template

var itemStore = todoitem.NewItemStore()

func init() {
	page = template.New("main")
	page = template.Must(page.ParseFiles("index.tmpl"))
}

func handleRootRequest(w http.ResponseWriter, r *http.Request) {
	items := itemStore.FindAll()

	if err := page.ExecuteTemplate(w, "main", items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleItemRequest(w http.ResponseWriter, r *http.Request) {
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

func main() {

	http.HandleFunc("/", handleRootRequest)
	http.HandleFunc("/item", handleItemRequest)
	http.ListenAndServe(":80", nil)
}
