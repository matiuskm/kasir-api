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
	"kasir-api/middlewares"
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
	PORT      string `mapstructure:"PORT"`
	DB_CONN   string `mapstructure:"DB_CONN"`
	REPORT_TZ string `mapstructure:"REPORT_TZ"`
	API_KEY	 	string `mapstructure:"API_KEY"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		PORT:      viper.GetString("PORT"),
		DB_CONN:   viper.GetString("DB_CONN"),
		REPORT_TZ: viper.GetString("REPORT_TZ"),
		API_KEY:   viper.GetString("API_KEY"),
	}
	if strings.TrimSpace(config.REPORT_TZ) == "" {
		config.REPORT_TZ = "Asia/Jakarta"
	}

	// setup database
	db, err := database.InitDB(config.DB_CONN)
	if err != nil {
		fmt.Println("failed to initialize database:", err)
	}
	defer db.Close()

	apiKeyMiddleware := middlewares.APIKey(config.API_KEY)

	// Setup HTTP Handlers
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo, config.REPORT_TZ)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// setup routes
	http.HandleFunc("/api/products", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(productHandler.HandleProducts))))
	http.HandleFunc("/api/products/", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(productHandler.HandleProductByID))))
	http.HandleFunc("/api/categories", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(categoryHandler.HandleCategories))))
	http.HandleFunc("/api/categories/", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(categoryHandler.HandleCategoryByID))))
	http.HandleFunc("/api/checkout", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(transactionHandler.HandleCheckout))))
	http.HandleFunc("/api/transactions", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(transactionHandler.HandleTransactions))))
	http.HandleFunc("/api/transactions/", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(transactionHandler.HandleTransactionByID))))
	http.HandleFunc("/api/report/hari-ini", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(transactionHandler.HandleReportHariIni))))
	http.HandleFunc("/api/report", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(transactionHandler.HandleReportRange))))
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Service is healthy",
		})
	})

	fmt.Println("server running on port " + config.PORT)
	err = http.ListenAndServe(":"+config.PORT, nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
