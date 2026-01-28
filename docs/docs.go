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
