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
    "log"
    _ "github.com/go-sql-driver/mysql"
)

type Stu struct {
	id   int
	Name string
}



func main() {
	r := mux.NewRouter()
	r.HandleFunc("/info", info).Methods("GET")
	r.HandleFunc("/categories", categories).Methods("GET")
	r.HandleFunc("/products", products).Methods("GET")
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

type Product struct {
    ID           int             `json:"id"`
    GID          int             `json:"gId"`
    GIDName      string          `json:"gIdName"`
    Name         string          `json:"name"`
    HebName      string          `json:"heb_name"`
    Price        float64         `json:"price"`
    Currency     int             `json:"currency"`
    PictureFolder string         `json:"picture_folder"`
    Color        string          `json:"color"`
    Category     int             `json:"category"`
    Sizes        json.RawMessage `json:"sizes"` // stored as JSON
    SizesIsrael  string          `json:"sizes_israel"`
    Description  string          `json:"description"`
    DescHeb      string          `json:"desc_heb"`
    About        string          `json:"about"`
    AboutHeb     string          `json:"about_heb"`
    CareHeb      string          `json:"care_heb"`
    Care         string          `json:"care"`
    Fabric       string          `json:"fabric"`
    FabricHeb    string          `json:"fabric_heb"`
}

func products(w http.ResponseWriter, r *http.Request){
	dsn := "u981786471_meir:mp496285MP@tcp(fr-int-web2000.main-hosting.eu:3306)/u981786471_bergs?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("DB connection failed:", err)
    }
    defer db.Close()

    // http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")

        id := r.URL.Query().Get("id")
        catid := r.URL.Query().Get("catid")

        var rows *sql.Rows
        var row *sql.Row
        var result any

        if id != "" {
            if _, err := strconv.Atoi(id); err == nil {
                row = db.QueryRow("SELECT * FROM products WHERE id = ? ORDER BY id ASC", id)
                var p Product
                err := row.Scan(
                    &p.ID, &p.GID, &p.GIDName, &p.Name, &p.HebName, &p.Price, &p.Currency,
                    &p.PictureFolder, &p.Color, &p.Category, &p.Sizes, &p.SizesIsrael,
                    &p.Description, &p.DescHeb, &p.About, &p.AboutHeb,
                    &p.CareHeb, &p.Care, &p.Fabric, &p.FabricHeb,
                )
                if err == sql.ErrNoRows {
                    json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
                    return
                } else if err != nil {
                    json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
                    return
                }
                result = p
            } else {
                json.NewEncoder(w).Encode(map[string]string{"error": "Invalid id format"})
                return
            }

        } else if catid != "" {
            var query string
            if _, err := strconv.Atoi(catid); err == nil {
                query = "SELECT * FROM products WHERE category = ? ORDER BY id ASC"
            } else {
                query = "SELECT * FROM products WHERE gIdName = ? ORDER BY id ASC"
            }

            rows, err = db.Query(query, catid)
            if err != nil {
                json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				log.Fatalf("Error opening DB: %v", err)
                return
            }
            defer rows.Close()

            var products []Product
            for rows.Next() {
                var p Product
                err := rows.Scan(
                    &p.ID, &p.GID, &p.GIDName, &p.Name, &p.HebName, &p.Price, &p.Currency,
                    &p.PictureFolder, &p.Color, &p.Category, &p.Sizes, &p.SizesIsrael,
                    &p.Description, &p.DescHeb, &p.About, &p.AboutHeb,
                    &p.CareHeb, &p.Care, &p.Fabric, &p.FabricHeb,
                )
                if err != nil {
                    json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
                    return
                }
                products = append(products, p)
            }
            result = products

        } else {
            json.NewEncoder(w).Encode(map[string]string{"error": "Either id or catid parameter is required"})
            return
        }

        // Output result
        json.NewEncoder(w).Encode(result)
    

    // fmt.Println("Server running on http://localhost:8080/products")
    // log.Fatal(http.ListenAndServe(":8080", nil))
    
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
