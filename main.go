package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Category struct {
	ID					int    `json:"id"`
	Name 				string `json:"name"`
	Description string `json:"description"`
}

var categoryList = []Category{
	{ID: 1, Name: "Makanan Ringan", Description: "Camilan dan makanan ringan lainnya"},
	{ID: 2, Name: "Minuman", Description: "Berbagai jenis minuman segar"},
}

const (
	contentTypeHeader = "Content-Type"
	jsonContentType   = "application/json"
)

func createCategory(w http.ResponseWriter, r *http.Request) {
	// baca data dari request body
	var categoryBaru Category
	err := json.NewDecoder(r.Body).Decode(&categoryBaru)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid request body",
		})
		return
	}

	// masukkan ke slice categoryList
	categoryBaru.ID = len(categoryList) + 1
	categoryList = append(categoryList, categoryBaru)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string {
		"message": "kategori berhasil ditambahkan",
	})
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// ambil id dari url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid kategori id",
		})
		return
	}

	// cari kategori berdasarkan id
	for _, category := range categoryList {
		if category.ID == id {
			json.NewEncoder(w).Encode(category)
			return
		}
	}
	
	// jika tidak ditemukan
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string {
		"error": "kategori not found",
	})
}

func updateCategoryByID(w http.ResponseWriter, r *http.Request) {
	// ambil id dari url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid kategori id",
		})
		return
	}

	// baca data dari request body
	var categoryUpdate Category
	err = json.NewDecoder(r.Body).Decode(&categoryUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid request body",
		})
		return
	}

	// cari dan update kategori
	for i	 := range categoryList {
		if categoryList[i].ID == id {
			categoryList[i].Name = categoryUpdate.Name
			categoryList[i].Description = categoryUpdate.Description
			json.NewEncoder(w).Encode(map[string]string {
				"message": "kategori berhasil diupdate",
				"data": fmt.Sprintf("%+v", categoryList[i]),
			})
			return
		}
	}
	
	// jika tidak ditemukan
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string {
		"error": "kategori not found",
	})
}

func deleteCategoryByID(w http.ResponseWriter, r *http.Request) {
	// ambil id dari url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid kategori id",
		})
		return
	}

	// cari dan hapus kategori
	for i	 := range categoryList {
		if categoryList[i].ID == id {
			categoryList = append(categoryList[:i], categoryList[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string {
				"message": "kategori berhasil dihapus",
			})
			return
		}
	}

	// jika tidak ditemukan
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string {
		"error": "kategori not found",
	})
}

type Config struct {
	PORT 		string `mapstructure:"PORT"`
	DB_CONN string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		PORT: 		viper.GetString("PORT"),
		DB_CONN: 	viper.GetString("DB_CONN"),
	}

	// setup database
	db, err := database.InitDB(config.DB_CONN)
	if err != nil {
		fmt.Println("failed to initialize database:", err)
	}
	defer db.Close()

	// Setup HTTP Handlers
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// setup routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)

	// ====== API KATEGORI ======

	// DELETE localhost:8080/api/categories/{id}
	// PUT localhost:8080/api/categories/{id}
	// GET localhost:8080/api/categories/{id}
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)

		switch r.Method {
			case http.MethodGet:
				getCategoryByID(w, r)
			case http.MethodPut:
				updateCategoryByID(w, r)
			case http.MethodDelete:
				deleteCategoryByID(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// POST localhost:8080/api/categories
	// GET localhost:8080/api/categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)

		switch r.Method {
			case http.MethodGet:
				json.NewEncoder(w).Encode(categoryList)
			case http.MethodPost:
				createCategory(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)
		json.NewEncoder(w).Encode(map[string]string {
			"status": "OK",
			"message": "Service is healthy",
		})
	})

	fmt.Println("server running on port " + config.PORT)
	err = http.ListenAndServe(":" + config.PORT, nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
