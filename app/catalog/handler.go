package catalog

import (
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
)

// ProductsRepository defines the interface for accessing product data
type ProductsRepository interface {
	GetProductsByFilter(offset, limit int, categoryID *uint, priceLessThan *decimal.Decimal) ([]models.Product, int64, error)
	GetProductByCode(code string) (*models.Product, error)
}

type CatalogHandler struct {
	repo ProductsRepository
}

func NewCatalogHandler(r ProductsRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

type ProductResponse struct {
	Code     string                `json:"code"`
	Price    float64               `json:"price"`
	Category *CategoryResponse     `json:"category,omitempty"`
	Variants []VariantResponse     `json:"variants,omitempty"`
}

type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type VariantResponse struct {
	Name  string   `json:"name"`
	SKU   string   `json:"sku"`
	Price *float64 `json:"price,omitempty"`
}

type CatalogListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64            `json:"total"`
}

type ProductDetailResponse struct {
	Code     string           `json:"code"`
	Price    float64          `json:"price"`
	Category CategoryResponse `json:"category"`
	Variants []VariantResponse `json:"variants"`
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Parse optional filters
	var categoryID *uint
	if categoryCode := r.URL.Query().Get("category"); categoryCode != "" {
		// We'll validate the category code in the repository query
		// For now, we pass it as a query parameter filter
		categoryIDVal, err := parseUintFromString(categoryCode)
		if err == nil {
			categoryID = &categoryIDVal
		}
	}

	var priceLessThan *decimal.Decimal
	if priceStr := r.URL.Query().Get("price_less_than"); priceStr != "" {
		price, err := decimal.NewFromString(priceStr)
		if err == nil {
			priceLessThan = &price
		}
	}

	products, total, err := h.repo.GetProductsByFilter(offset, limit, categoryID, priceLessThan)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map products to response
	productResponses := make([]ProductResponse, len(products))
	for i, p := range products {
		productResponses[i] = mapProductToResponse(p, false)
	}

	response := CatalogListResponse{
		Products: productResponses,
		Total:    total,
	}

	api.OKResponse(w, response)
}

func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	product, err := h.repo.GetProductByCode(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	// Map to detail response
	variantResponses := make([]VariantResponse, len(product.Variants))
	for i, v := range product.Variants {
		variantResponses[i] = mapVariantToResponse(v, &product.Price)
	}

	var categoryResp CategoryResponse
	if product.Category != nil {
		categoryResp = CategoryResponse{
			Code: product.Category.Code,
			Name: product.Category.Name,
		}
	}

	response := ProductDetailResponse{
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Category: categoryResp,
		Variants: variantResponses,
	}

	api.OKResponse(w, response)
}

func mapProductToResponse(p models.Product, includeVariants bool) ProductResponse {
	resp := ProductResponse{
		Code:  p.Code,
		Price: p.Price.InexactFloat64(),
	}

	if p.Category != nil {
		resp.Category = &CategoryResponse{
			Code: p.Category.Code,
			Name: p.Category.Name,
		}
	}

	if includeVariants {
		resp.Variants = make([]VariantResponse, len(p.Variants))
		for i, v := range p.Variants {
			resp.Variants[i] = mapVariantToResponse(v, &p.Price)
		}
	}

	return resp
}

func mapVariantToResponse(v models.Variant, productPrice *decimal.Decimal) VariantResponse {
	resp := VariantResponse{
		Name: v.Name,
		SKU:  v.SKU,
	}

	// Use variant price if available, otherwise use product price
	if !v.Price.IsZero() {
		price := v.Price.InexactFloat64()
		resp.Price = &price
	} else if productPrice != nil {
		price := productPrice.InexactFloat64()
		resp.Price = &price
	}

	return resp
}

func parseUintFromString(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	return uint(val), err
}
