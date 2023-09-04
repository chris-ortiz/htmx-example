package todoitem

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type TodoItem struct {
	Id   int64
	Text string
}
type ItemStore interface {
	Add(item *TodoItem)
	FindAll() []TodoItem
	Delete(id int)
}

type ItemDB struct {
	db *sql.DB
}

func (itemDB ItemDB) Delete(id int) {
	_, _ = itemDB.db.Exec("DELETE FROM item WHERE id = ?", id)
}

func (itemDB ItemDB) Add(item *TodoItem) {
	result, err := itemDB.db.Exec("INSERT INTO item (description) VALUES (?)", item.Text)
	if err != nil {
		log.Fatal(err)
	}

	item.Id, _ = result.LastInsertId()
}

func (itemDB ItemDB) FindAll() []TodoItem {
	rows, err := itemDB.db.Query("SELECT * from item")

	if err != nil {
		log.Fatal(err)
	}

	var items []TodoItem

	for rows.Next() {
		var item TodoItem
		_ = rows.Scan(&item.Id, &item.Text)
		items = append(items, item)
	}

	return items
}

func NewItemStore() ItemStore {
	db := createDB()
	return &ItemDB{db: db}
}

func createDB() *sql.DB {
	db, err := sql.Open("sqlite3", "item.sqlite")

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS item (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   description text NOT NULL);`)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
