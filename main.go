package main

import (
	"context" // For context.Context, used in MongoDB operations
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	// "net/url" // Added for URL encoding
	// "os"
	"strconv"
	// "strings" // Added for string manipulation
	// "time"    // Added for context timeout

	// "github.com/joho/godotenv"
	"github.com/gorilla/mux" // For routing

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver (pgx)
	_ "github.com/lib/pq"              // PostgreSQL driver (pq) - often used with pgx
)

// Stu struct for SQL student table
type Stu struct {
	ID   int    `json:"id"`   // Added json tag for consistency, note: field names should be exported
	Name string `json:"name"` // Exported field for JSON marshalling and scanning
}

// Category struct for SQL category table
type Category struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Idcol  int    `json:"idcol"`
	Hebrew string `json:"hebrew"`
}

// Product struct for both SQL products and MongoDB products collection
// Added `bson` tags for MongoDB mapping
type Product struct {
	ID            int             `json:"id" bson:"id"`
	GID           int             `json:"gId" bson:"gId"`
	GIDName       string          `json:"gIdName" bson:"gIdName"`
	Name          string          `json:"name" bson:"name"`
	HebName       string          `json:"heb_name" bson:"heb_name"`
	Price         float64         `json:"price" bson:"price"`
	Currency      int             `json:"currency" bson:"currency"`
	PictureFolder string          `json:"picture_folder" bson:"picture_folder"`
	Color         string          `json:"color" bson:"color"`
	Category      int             `json:"category" bson:"category"`
	Sizes         json.RawMessage `json:"sizes" bson:"sizes"` // stored as JSON
	SizesIsrael   string          `json:"sizes_israel" bson:"sizes_israel"`
	Description   string          `json:"description" bson:"description"`
	DescHeb       string          `json:"desc_heb" bson:"desc_heb"`
	About         string          `json:"about" bson:"about"`
	AboutHeb      string          `json:"about_heb" bson:"about_heb"`
	CareHeb       string          `json:"care_heb" bson:"care_heb"`
	Care          string          `json:"care" bson:"care"`
	Fabric        string          `json:"fabric" bson:"fabric"`
	FabricHeb     string          `json:"fabric_heb" bson:"fabric_heb"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/info", info).Methods("GET")
	r.HandleFunc("/categories", categories).Methods("GET")
	r.HandleFunc("/products", products).Methods("GET")
	r.HandleFunc("/products_mongo", products_mongo).Methods("GET")
	r.HandleFunc("/getsingle", getsingle).Methods("GET")
	r.HandleFunc("/add", add).Methods("POST")
	r.HandleFunc("/update/{id}", update).Methods("PUT")
	r.HandleFunc("/delete/{id}", deleteHandler).Methods("DELETE") // Renamed to avoid conflict with built-in
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// info handler for SQL student data
func info(w http.ResponseWriter, r *http.Request) {
	connStr := fmt.Sprintf("host=9qasp5v56q8ckkf5dc.leapcellpool.com port=6438 user=ufnsrbazgcetcbqwevru password=zmkiotezqmcqwpwsvrjnsxtmydznos dbname=cuarxlxvaahzbgdyqnep sslmode=require")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening DB connection: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	res, err := db.Query("select id, name from student") // Specify columns to avoid issues with `id` field
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying student table: %v", err), http.StatusInternalServerError)
		return
	}
	defer res.Close()

	for res.Next() {
		var stu Stu
		// Ensure `id` is exported in Stu struct (Stu.ID) for Scan to work
		if err := res.Scan(&stu.ID, &stu.Name); err != nil {
			log.Printf("Error scanning student row: %v", err)
			continue
		}
		str := "The product name is:" + stu.Name + " My id is:" + strconv.Itoa(stu.ID)
		fmt.Fprintln(w, str)
	}
	fmt.Fprintln(w, "Thats all")
}

// categories handler for SQL category data
func categories(w http.ResponseWriter, r *http.Request) {
	connStr := fmt.Sprintf("host=9qasp5v56q8ckkf5dc.leapcellpool.com port=6438 user=ufnsrbazgcetcbqwevru password=zmkiotezqmcqwpwsvrjnsxtmydznos dbname=cuarxlxvaahzbgdyqnep sslmode=require")
	db, err := sql.Open("pgx", connStr) // Using "pgx" driver
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening DB connection: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	res, err := db.Query("SELECT ID, Idcol, Name, Hebrew FROM category ORDER BY idcol ASC;")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying categories: %v", err), http.StatusInternalServerError)
		return
	}
	defer res.Close()

	var categories []Category
	for res.Next() {
		var cat Category
		if err := res.Scan(&cat.ID, &cat.Idcol, &cat.Name, &cat.Hebrew); err != nil {
			log.Printf("Error scanning category row: %v", err)
			continue
		}
		categories = append(categories, cat)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding categories to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

// products handler for SQL product data
func products(w http.ResponseWriter, r *http.Request) {
	dsn := "u981786471_meir:mp496285MP@tcp(fr-int-web2000.main-hosting.eu:3306)/u981786471_bergs?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		http.Error(w, fmt.Sprintf("DB connection failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	id := r.URL.Query().Get("id")
	catid := r.URL.Query().Get("catid")

	var result any

	if id != "" {
		if _, err := strconv.Atoi(id); err == nil {
			row := db.QueryRow("SELECT ID, GID, GIDName, Name, HebName, Price, Currency, PictureFolder, Color, Category, Sizes, SizesIsrael, Description, DescHeb, About, AboutHeb, CareHeb, Care, Fabric, FabricHeb FROM products WHERE id = ?", id)
			var p Product
			err := row.Scan(
				&p.ID, &p.GID, &p.GIDName, &p.Name, &p.HebName, &p.Price, &p.Currency,
				&p.PictureFolder, &p.Color, &p.Category, &p.Sizes, &p.SizesIsrael,
				&p.Description, &p.DescHeb, &p.About, &p.AboutHeb,
				&p.CareHeb, &p.Care, &p.Fabric, &p.FabricHeb,
			)
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
				return
			} else if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning product row: %v", err), http.StatusInternalServerError)
				return
			}
			result = p
		} else {
			http.Error(w, `{"error": "Invalid id format"}`, http.StatusBadRequest)
			return
		}

	} else if catid != "" {
		var query string
		if _, err := strconv.Atoi(catid); err == nil {
			query = "SELECT ID, GID, GIDName, Name, HebName, Price, Currency, PictureFolder, Color, Category, Sizes, SizesIsrael, Description, DescHeb, About, AboutHeb, CareHeb, Care, Fabric, FabricHeb FROM products WHERE category = ? ORDER BY id ASC"
		} else {
			query = "SELECT ID, GID, GIDName, Name, HebName, Price, Currency, PictureFolder, Color, Category, Sizes, SizesIsrael, Description, DescHeb, About, AboutHeb, CareHeb, Care, Fabric, FabricHeb FROM products WHERE gIdName = ? ORDER BY id ASC"
		}

		rows, err := db.Query(query, catid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying products: %v", err), http.StatusInternalServerError)
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
				http.Error(w, fmt.Sprintf("Error scanning product row: %v", err), http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}
		result = products

	} else {
		http.Error(w, `{"error": "Either id or catid parameter is required"}`, http.StatusBadRequest)
		return
	}

	// Output result
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding result to JSON: %v", err), http.StatusInternalServerError)
	}
}

// products_mongo handler for MongoDB product data
func products_mongo(w http.ResponseWriter, r *http.Request) {
	// Load environment variables from .env file
	MONGODB_URI := "mongodb+srv://meir:mp-496285MP@bergs.9zb9ptn.mongodb.net/?retryWrites=true&w=majority&appName=Bergs"

	uri := MONGODB_URI
	docs := "www.mongodb.com/docs/drivers/go/current/"

	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " + docs +
			"usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// id := r.URL.Query().Get("id")
	catid := r.URL.Query().Get("catid")

	coll := client.Database("Products").Collection("products")
	title := "rgo6209"
	var result bson.M
	i, err := strconv.Atoi(catid)
	cursor, err := coll.Find(context.TODO(), bson.D{{Key: "category", Value: i}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
		if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding categories to JSON: %v", err), http.StatusInternalServerError)
		return
	}
	}

	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", title)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding categories to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

// getsingle handler for SQL product data by gIdName
func getsingle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

	if r.Method == http.MethodOptions {
		return // Handle CORS preflight request
	}

	catid := r.URL.Query().Get("catid")
	if catid == "" {
		http.Error(w, `{"error":"Missing catid parameter"}`, http.StatusBadRequest)
		return
	}

	dsn := "u981786471_meir:mp496285MP@tcp(fr-int-web2000.main-hosting.eu:3306)/u981786471_bergs?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Error opening DB: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT ID, GID, GIDName, Name, HebName, Price, Currency, PictureFolder, Color, Category, Sizes, SizesIsrael, Description, DescHeb, About, AboutHeb, CareHeb, Care, Fabric, FabricHeb FROM products WHERE gIdName = ? ORDER BY id ASC", catid)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Query error: %v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(
			&p.ID, &p.GID, &p.GIDName, &p.Name, &p.HebName,
			&p.Price, &p.Currency, &p.PictureFolder, &p.Color, &p.Category,
			&p.Sizes, &p.SizesIsrael, &p.Description, &p.DescHeb, &p.About,
			&p.AboutHeb, &p.CareHeb, &p.Care, &p.Fabric, &p.FabricHeb,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Row scan error: %v"}`, err), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding products to JSON: %v", err), http.StatusInternalServerError)
	}
}

// add handler for SQL student data
func add(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	var stu Stu
	if err := json.Unmarshal(data, &stu); err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:mp496285MP@tcp(127.0.0.1:3306)/bergs")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening DB connection: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Use prepared statements to prevent SQL injection
	stmt, err := db.Prepare("INSERT INTO student (name) VALUES(?)")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error preparing statement: %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(stu.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting data: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, stu.Name+" is my name and added to database")
}

// update handler for SQL student data
func update(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	var stu Stu
	if err := json.Unmarshal(data, &stu); err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:mp496285MP@tcp(127.0.0.1:3306)/bergs")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening DB connection: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Use prepared statements to prevent SQL injection
	stmt, err := db.Prepare("UPDATE student SET name = ? WHERE id = ?")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error preparing statement: %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(stu.Name, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating data: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Data is updated")
}

// deleteHandler for SQL student data
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:mp496285MP@tcp(127.0.0.1:3306)/bergs")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening DB connection: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM student WHERE id = ?")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error preparing statement: %v", err), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting data: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No student found with the given ID", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "Data is deleted")
}
