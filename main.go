package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.elastic.co/apm/module/apmgorilla/v2"

	postgres "go.elastic.co/apm/module/apmgormv2/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	var err error
	db, err = gorm.Open(postgres.Open("host=localhost user=goapm sslmode=disable password=secret dbname=goapm"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//db := db.WithContext(context.Background())

	// Migrate the schema
	db.Debug().AutoMigrate(&Product{})

	r := mux.NewRouter()
	r.HandleFunc("/hello/{name}", handler)
	apmgorilla.Instrument(r)
	log.Fatal(http.ListenAndServe(":8000", r))

}

func handler(w http.ResponseWriter, req *http.Request) {

	db := db.WithContext(req.Context())

	// Create
	db.Debug().Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.Debug().First(&product)

	vars := mux.Vars(req)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Hello %s", vars["name"]),
	})
}
