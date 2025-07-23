package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Stu struct {
	id   int
	Name string
}



func main() {
	r := mux.NewRouter()
	r.HandleFunc("/info", info).Methods("GET")
	r.HandleFunc("/categories", categories).Methods("GET")
	r.HandleFunc("/add", add).Methods("POST")
	r.HandleFunc("/update/{id}", update).Methods("PUT")
	http.ListenAndServe(":8080", r)
}

func info(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("hello")
	// db, err := sql.Open("mysql", "root:mp496285MP@tcp(127.0.0.1:3306)/bergs")
	connStr := fmt.Sprintf("host=9qasp5v56q8ckkf5dc.leapcellpool.com port=6438 user=ufnsrbazgcetcbqwevru password=zmkiotezqmcqwpwsvrjnsxtmydznos dbname=cuarxlxvaahzbgdyqnep sslmode=require")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("select * from student")
	if err != nil {
		panic(err)
	}
	defer res.Close()

	for res.Next() {
		var stu Stu
		res.Scan(&stu.id, &stu.Name)
		str := "The product name is:" + stu.Name + " My id is:" + strconv.Itoa(stu.id)
		fmt.Fprintln(w, str)
	}
	fmt.Fprintln(w, "Thats all")
}

type Category struct {
	ID     int    `json:"id"`
	Name   string `json:"name"` 
	Idcol  int 		`json:"idcol"` 
	Hebrew string `json:"hebrew"` 
}
func categories(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("hello")
	// db, err := sql.Open("mysql", "root:mp496285MP@tcp(127.0.0.1:3306)/bergs")
	connStr := fmt.Sprintf("host=9qasp5v56q8ckkf5dc.leapcellpool.com port=6438 user=ufnsrbazgcetcbqwevru password=zmkiotezqmcqwpwsvrjnsxtmydznos dbname=cuarxlxvaahzbgdyqnep sslmode=require")
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM category;")
	if err != nil {
		panic(err)
	}
	defer res.Close()
	var categories []Category
	for res.Next() {
		var stu Category
		res.Scan(&stu.ID, &stu.Idcol, &stu.Name, &stu.Hebrew)
		// str := "The product name is:" + stu.Name + " Idcol is:" + strconv.Itoa(stu.Idcol)
		categories = append(categories, stu)
		// jsonData, err := json.Marshal(categories)
		// fmt.Fprintf(w, string(jsonData))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	jsonData, err := json.Marshal(categories)
	// fmt.Fprintln(w, jsonData)
	// Set proper content type and write JSON bytes
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)  // This writes the actual JSON
}

func add(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)
	var stu Stu
	json.Unmarshal(data, &stu)

	fmt.Fprintln(w, stu.Name)

	db, err := sql.Open("mysql", "root:mp496285MP@tcp(127.0.0.1:3306)/bergs")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("insert into student (name) values('" + stu.Name + "')")
	if err != nil {
		panic(err)
	}
	defer res.Close()

	fmt.Fprintln(w, stu.Name+" is my name and added to database")
}

func update(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)
	var stu Stu
	json.Unmarshal(data, &stu)
	id := mux.Vars(r)["id"]

	db, err := sql.Open("mysql", "root:mp496285MP@tcp(127.0.0.1:3306)/bergs")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("update student set name='" + stu.Name + "' where id=" + id)
	if err != nil {
		panic(err)
	}
	defer res.Close()

	fmt.Fprintln(w, "Data is updated")
}

func delete(w http.ResponseWriter, r *http.Request) {

}
