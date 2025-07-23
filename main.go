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
)

type Stu struct {
	id   int
	Name string
}

var host = "9qasp5v56q8ckkf5dc.leapcellpool.com user=ufnsrbazgcetcbqwevru password=zmkiotezqmcqwpwsvrjnsxtmydznos dbname=cuarxlxvaahzbgdyqnep port=6438 sslmode=require"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/info", info).Methods("GET")
	r.HandleFunc("/add", add).Methods("POST")
	r.HandleFunc("/update/{id}", update).Methods("PUT")
	http.ListenAndServe(":8080", r)
}

func info(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("hello")
	db, err := sql.Open("postgresql", host)
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
