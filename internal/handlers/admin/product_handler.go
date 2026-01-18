package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"hosting-backend/internal/models"
	"github.com/go-chi/chi/v5"
)

// CreateProductHandler cria um novo produto.
// Rota: POST /admin/products
func CreateProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p models.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}

		// Validação básica
		if p.Name == "" || p.Price <= 0 {
			http.Error(w, "Nome e preço são obrigatórios e o preço deve ser positivo", http.StatusBadRequest)
			return
		}

		lastInsertID, err := models.CreateProduct(db, &p)
		if err != nil {
			http.Error(w, "Erro ao criar produto", http.StatusInternalServerError)
			return
		}

		p.ID = int(lastInsertID)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

// GetProductsHandler lista todos os produtos.
// Rota: GET /admin/products
func GetProductsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := models.GetAllProducts(db)
		if err != nil {
			http.Error(w, "Erro ao buscar produtos", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

// GetProductHandler busca um único produto pelo ID.
// Rota: GET /admin/products/{id}
func GetProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		product, err := models.GetProductByID(db, id)
		if err != nil {
			http.Error(w, "Erro ao buscar produto", http.StatusInternalServerError)
			return
		}
		if product == nil {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

// UpdateProductHandler atualiza um produto existente.
// Rota: PUT /admin/products/{id}
func UpdateProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		var p models.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Corpo da requisição inválido", http.StatusBadRequest)
			return
		}
		p.ID = id // Garante que o ID da URL seja usado

		if err := models.UpdateProduct(db, &p); err != nil {
			http.Error(w, "Erro ao atualizar produto", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

// DeleteProductHandler remove um produto.
// Rota: DELETE /admin/products/{id}
func DeleteProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		if err := models.DeleteProduct(db, id); err != nil {
			http.Error(w, "Erro ao deletar produto", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
