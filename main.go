package main

// @title Kasir API
// @version 1.0.0
// @description API documentation for Kasir.
// @host kasir-api-production-8ed0.up.railway.app
// @BasePath /

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	_ "kasir-api/docs"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger/v2" 
)

const (
	contentTypeHeader = "Content-Type"
	jsonContentType   = "application/json"
)

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
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// setup routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)
	http.Handle("/swagger/", httpSwagger.WrapHandler)

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
