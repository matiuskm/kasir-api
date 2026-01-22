package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID		int     `json:"id"`
	Nama  string  `json:"nama"`
	Harga int 		`json:"harga"`
	Stok 	int 		`json:"stok"`
}

type Category struct {
	ID					int    `json:"id"`
	Name 				string `json:"name"`
	Description string `json:"description"`
}

var produkList = []Produk{
	{ID: 1, Nama: "Indomie Rebus", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Aqua 200ml", Harga: 2000, Stok: 25},
}

var categoryList = []Category{
	{ID: 1, Name: "Makanan Ringan", Description: "Camilan dan makanan ringan lainnya"},
	{ID: 2, Name: "Minuman", Description: "Berbagai jenis minuman segar"},
}

const (
	contentTypeHeader = "Content-Type"
	jsonContentType   = "application/json"
)

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	// ambil id dari url
		idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string {
				"error": "invalid produk id",
			})
			return
		}

		// cari produk berdasarkan id
			for _, produk := range produkList {
				if produk.ID == id {
					json.NewEncoder(w).Encode(produk)
					return
				}
			}
			
			// jika tidak ditemukan
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string {
				"error": "produk not found",
			})
}

func updateProdukByID(w http.ResponseWriter, r *http.Request) {
	// ambil id dari url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid produk id",
		})
		return
	}

	// baca data dari request body
	var produkUpdate Produk
	err = json.NewDecoder(r.Body).Decode(&produkUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid request body",
		})
		return
	}

	// cari dan update produk
	for i	 := range produkList {
		if produkList[i].ID == id {
			produkList[i].Nama = produkUpdate.Nama
			produkList[i].Harga = produkUpdate.Harga
			produkList[i].Stok = produkUpdate.Stok
			json.NewEncoder(w).Encode(map[string]string {
				"message": "produk berhasil diupdate",
				"data": fmt.Sprintf("%+v", produkList[i]),
			})
			return
		}
	}

	// jika tidak ditemukan
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string {
		"error": "produk not found",
	})
}

func hapusProdukByID(w http.ResponseWriter, r *http.Request) {
	// ambil id dari url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid produk id",
		})
		return
	}

	// cari dan hapus produk
	for i	 := range produkList {
		if produkList[i].ID == id {
			produkList = append(produkList[:i], produkList[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string {
				"message": "produk berhasil dihapus",
			})
			return
		}
	}

	// jika tidak ditemukan
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string {
		"error": "produk not found",
	})
}

func createProduk(w http.ResponseWriter, r *http.Request) {
	// baca data dari request body
	var produkBaru Produk
	err := json.NewDecoder(r.Body).Decode(&produkBaru)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid request body",
		})
		return
	}

	// masukkan ke slice produkList
	produkBaru.ID = len(produkList) + 1
	produkList = append(produkList, produkBaru)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string {
		"message": "produk berhasil ditambahkan",
	})
}

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
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
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
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
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
	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
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

func main() {

	// ====== API KATEGORI ======

	// DELETE localhost:8080/api/kategori/{id}
	// PUT localhost:8080/api/kategori/{id}
	// GET localhost:8080/api/kategori/{id}
	http.HandleFunc("/api/kategori/", func(w http.ResponseWriter, r *http.Request) {
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
				

	// POST localhost:8080/api/kategori
	// GET localhost:8080/api/kategori
	http.HandleFunc("/api/kategori", func(w http.ResponseWriter, r *http.Request) {
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

	// ====== API PRODUK ======

	// DELETE localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// GET localhost:8080/api/produk/{id}
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)

		switch r.Method {
			case http.MethodPut:
				updateProdukByID(w, r)
			case http.MethodGet:
				getProdukByID(w, r)
			case http.MethodDelete:
				hapusProdukByID(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// GET localhost:8080/api/produk
	// POST localhost:8080/api/produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)
		
		switch r.Method {
			case http.MethodGet:
				json.NewEncoder(w).Encode(produkList)
			case http.MethodPost:
				createProduk(w, r)
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

	fmt.Println("server running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
