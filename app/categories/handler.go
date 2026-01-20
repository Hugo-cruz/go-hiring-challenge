package categories

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

// CategoriesRepository defines the interface for accessing category data
type CategoriesRepository interface {
	GetAll() ([]models.Category, error)
	Create(category *models.Category) error
	FindByCode(code string) (*models.Category, error)
}

type CategoriesHandler struct {
	repo CategoriesRepository
}

func NewCategoriesHandler(r CategoriesRepository) *CategoriesHandler {
	return &CategoriesHandler{
		repo: r,
	}
}

type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoriesListResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// HandleList returns all categories.
func (h *CategoriesHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.GetAll()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	categoryResponses := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		categoryResponses[i] = CategoryResponse{
			Code: c.Code,
			Name: c.Name,
		}
	}

	response := CategoriesListResponse{
		Categories: categoryResponses,
	}

	api.OKResponse(w, response)
}

// HandleCreate creates a new category.
func (h *CategoriesHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Code == "" || req.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Code and name are required")
		return
	}

	category := &models.Category{
		Code: req.Code,
		Name: req.Name,
	}

	if err := h.repo.Create(category); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := CategoryResponse{
		Code: category.Code,
		Name: category.Name,
	}

	api.OKResponse(w, response)
}
