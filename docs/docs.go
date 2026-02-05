package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "description": "Kasir API documentation",
    "title": "Kasir API",
    "version": "1.0.0"
  },
  "host": "kasir-api-production-8ed0.up.railway.app",
  "basePath": "/",
  "schemes": ["https"],
  "paths": {
    "/health": {
      "get": {
        "summary": "Health check",
        "description": "Check if service is healthy",
        "produces": ["application/json"],
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    },
    "/api/products": {
      "get": {
        "summary": "List products",
        "description": "Retrieve all products with pagination",
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "query",
            "name": "page",
            "required": false,
            "type": "integer",
            "description": "Page number (starts at 1)"
          },
          {
            "in": "query",
            "name": "limit",
            "required": false,
            "type": "integer",
            "description": "Items per page (max 100)"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": { "$ref": "#/definitions/models.ProductListResponse" }
          }
        }
      },
      "post": {
        "summary": "Create product",
        "description": "Create a new product",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "product",
            "required": true,
            "schema": { "$ref": "#/definitions/models.Product" }
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": { "$ref": "#/definitions/models.Product" }
          },
          "400": { "description": "Bad Request" }
        }
      }
    },
    "/api/products/{id}": {
      "get": {
        "summary": "Get product by ID",
        "description": "Retrieve a product by its ID",
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": { "$ref": "#/definitions/models.Product" }
          },
          "404": { "description": "Not Found" }
        }
      },
      "put": {
        "summary": "Update product",
        "description": "Update a product by its ID",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          },
          {
            "in": "body",
            "name": "product",
            "required": true,
            "schema": { "$ref": "#/definitions/models.Product" }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": { "$ref": "#/definitions/models.Product" }
          },
          "400": { "description": "Bad Request" }
        }
      },
      "delete": {
        "summary": "Delete product",
        "description": "Delete a product by its ID",
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          }
        ],
        "responses": {
          "204": { "description": "No Content" }
        }
      }
    },
    "/api/categories": {
      "get": {
        "summary": "List categories",
        "description": "Retrieve all categories with pagination",
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "query",
            "name": "page",
            "required": false,
            "type": "integer",
            "description": "Page number (starts at 1)"
          },
          {
            "in": "query",
            "name": "limit",
            "required": false,
            "type": "integer",
            "description": "Items per page (max 100)"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": { "$ref": "#/definitions/models.CategoryListResponse" }
          }
        }
      },
      "post": {
        "summary": "Create category",
        "description": "Create a new category",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "category",
            "required": true,
            "schema": { "$ref": "#/definitions/models.Category" }
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": { "$ref": "#/definitions/models.Category" }
          },
          "400": { "description": "Bad Request" }
        }
      }
    },
    "/api/categories/{id}": {
      "get": {
        "summary": "Get category by ID",
        "description": "Retrieve a category by its ID. Use include=products to include products.",
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          },
          {
            "in": "query",
            "name": "include",
            "required": false,
            "type": "string",
            "enum": ["products"]
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": { "$ref": "#/definitions/models.Category" }
          },
          "404": { "description": "Not Found" }
        }
      },
      "put": {
        "summary": "Update category",
        "description": "Update a category by its ID",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          },
          {
            "in": "body",
            "name": "category",
            "required": true,
            "schema": { "$ref": "#/definitions/models.Category" }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": { "$ref": "#/definitions/models.Category" }
          },
          "400": { "description": "Bad Request" }
        }
      },
      "delete": {
        "summary": "Delete category",
        "description": "Delete a category by its ID",
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          }
        ],
        "responses": {
          "204": { "description": "No Content" }
        }
      }
    }
    ,
    "/api/checkout": {
      "post": {
        "summary": "Create transaction (checkout)",
        "description": "Creates a transaction and returns the simplified details.",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "checkout",
            "required": true,
            "schema": { "$ref": "#/definitions/models.CheckoutRequest" }
          }
        ],
        "responses": {
          "201": {
            "description": "Created",
            "schema": { "$ref": "#/definitions/models.TransactionResponse" }
          },
          "400": { "description": "Bad Request" }
        }
      }
    },
    "/api/transactions": {
      "get": {
        "summary": "List transactions",
        "description": "Returns transactions with total products, total amount, and transaction date. Supports month or date range filter.",
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "query",
            "name": "page",
            "required": false,
            "type": "integer",
            "description": "Page number (starts at 1)"
          },
          {
            "in": "query",
            "name": "limit",
            "required": false,
            "type": "integer",
            "description": "Items per page (max 100)"
          },
          {
            "in": "query",
            "name": "month",
            "required": false,
            "type": "string",
            "description": "Filter by month (YYYY-MM)",
            "example": "2026-01"
          },
          {
            "in": "query",
            "name": "start_date",
            "required": false,
            "type": "string",
            "description": "Filter start date (YYYY-MM-DD)",
            "example": "2026-01-01"
          },
          {
            "in": "query",
            "name": "end_date",
            "required": false,
            "type": "string",
            "description": "Filter end date (YYYY-MM-DD)",
            "example": "2026-01-31"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": { "$ref": "#/definitions/models.TransactionListResponse" }
          },
          "400": { "description": "Bad Request" },
          "500": { "description": "Internal Server Error" }
        }
      }
    },
    "/api/transactions/{id}": {
      "get": {
        "summary": "Get transaction detail",
        "description": "Returns transaction detail with simplified item fields.",
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "required": true,
            "type": "integer"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": { "$ref": "#/definitions/models.TransactionResponse" }
          },
          "400": { "description": "Bad Request" },
          "404": { "description": "Not Found" },
          "500": { "description": "Internal Server Error" }
        }
      }
    },
    "/api/report/hari-ini": {
      "get": {
        "summary": "Today's sales report",
        "description": "Returns total revenue, total transactions, and top-selling products for today (REPORT_TZ).",
        "produces": ["application/json"],
        "responses": {
          "200": {
            "description": "OK",
            "schema": { "$ref": "#/definitions/models.DailyReport" }
          },
          "500": { "description": "Internal Server Error" }
        }
      }
    },
    "/api/report": {
      "get": {
        "summary": "Report by date range",
        "description": "Returns report for a date range (inclusive).",
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "query",
            "name": "start_date",
            "required": true,
            "type": "string",
            "description": "Start date (YYYY-MM-DD)",
            "example": "2026-01-01"
          },
          {
            "in": "query",
            "name": "end_date",
            "required": true,
            "type": "string",
            "description": "End date (YYYY-MM-DD)",
            "example": "2026-02-01"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": { "$ref": "#/definitions/models.DailyReport" }
          },
          "400": { "description": "Bad Request" },
          "500": { "description": "Internal Server Error" }
        }
      }
    }
  },
  "definitions": {
    "models.PaginationMeta": {
      "type": "object",
      "properties": {
        "total": { "type": "integer" },
        "page": { "type": "integer" },
        "page_size": { "type": "integer" },
        "has_more": { "type": "boolean" }
      }
    },
    "models.ProductListResponse": {
      "type": "object",
      "properties": {
        "data": {
          "type": "array",
          "items": { "$ref": "#/definitions/models.Product" }
        },
        "meta": { "$ref": "#/definitions/models.PaginationMeta" }
      }
    },
    "models.CategoryListResponse": {
      "type": "object",
      "properties": {
        "data": {
          "type": "array",
          "items": { "$ref": "#/definitions/models.Category" }
        },
        "meta": { "$ref": "#/definitions/models.PaginationMeta" }
      }
    },
    "models.Product": {
      "type": "object",
      "properties": {
        "id": { "type": "integer" },
        "name": { "type": "string" },
        "price": { "type": "integer" },
        "stock": { "type": "integer" },
        "category_id": { "type": "integer" },
        "category": { "$ref": "#/definitions/models.Category" }
      }
    },
    "models.Category": {
      "type": "object",
      "properties": {
        "id": { "type": "integer" },
        "name": { "type": "string" },
        "description": { "type": "string" },
        "products": {
          "type": "array",
          "items": { "$ref": "#/definitions/models.Product" }
        }
      }
    },
    "models.CheckoutItem": {
      "type": "object",
      "properties": {
        "product_id": { "type": "integer", "example": 4 },
        "quantity": { "type": "integer", "example": 2 }
      }
    },
    "models.CheckoutRequest": {
      "type": "object",
      "properties": {
        "items": {
          "type": "array",
          "items": { "$ref": "#/definitions/models.CheckoutItem" },
          "example": [{ "product_id": 4, "quantity": 2 }]
        }
      }
    },
    "models.TransactionDetailResponse": {
      "type": "object",
      "properties": {
        "product_name": { "type": "string", "example": "Indomie Goreng" },
        "quantity": { "type": "integer", "example": 2 },
        "subtotal": { "type": "integer", "example": 12000 }
      }
    },
    "models.TransactionResponse": {
      "type": "object",
      "properties": {
        "id": { "type": "integer", "example": 3 },
        "total_amount": { "type": "integer", "example": 16000 },
        "transaction_date": { "type": "string", "example": "2026-02-05" },
        "details": {
          "type": "array",
          "items": { "$ref": "#/definitions/models.TransactionDetailResponse" },
          "example": [
            { "product_name": "Indomie Soto", "quantity": 2, "subtotal": 11000 },
            { "product_name": "Coca Cola", "quantity": 1, "subtotal": 5000 }
          ]
        }
      }
    },
    "models.TransactionListItem": {
      "type": "object",
      "properties": {
        "id": { "type": "integer", "example": 3 },
        "total_products": { "type": "integer", "example": 3 },
        "total_amount": { "type": "integer", "example": 16000 },
        "transaction_date": { "type": "string", "example": "2026-02-05" }
      }
    },
    "models.TransactionListResponse": {
      "type": "object",
      "properties": {
        "data": {
          "type": "array",
          "items": { "$ref": "#/definitions/models.TransactionListItem" },
          "example": [
            { "id": 3, "total_products": 3, "total_amount": 16000, "transaction_date": "2026-02-05" }
          ]
        },
        "meta": { "$ref": "#/definitions/models.PaginationMeta" },
        "total_revenue": { "type": "integer", "example": 45000 }
      }
    },
    "models.BestProduct": {
      "type": "object",
      "properties": {
        "nama": { "type": "string", "example": "Indomie Goreng" },
        "qty_terjual": { "type": "integer", "example": 12 }
      }
    },
    "models.DailyReport": {
      "type": "object",
      "properties": {
        "total_revenue": { "type": "integer", "example": 45000 },
        "total_transaksi": { "type": "integer", "example": 5 },
        "produk_terlaris": {
          "type": "array",
          "items": { "$ref": "#/definitions/models.BestProduct" },
          "example": [{ "nama": "Indomie Goreng", "qty_terjual": 12 }]
        }
      }
    }
  }
}`

type swaggerInfo struct{}

func (s swaggerInfo) ReadDoc() string {
	return docTemplate
}

func init() {
	swag.Register("swagger", swaggerInfo{})
}
