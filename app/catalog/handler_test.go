package catalog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/mytheresa/go-hiring-challenge/models"
)

// MockProductsRepository is a mock implementation of ProductsRepository
type MockProductsRepository struct {
	mock.Mock
}

func (m *MockProductsRepository) GetAllProducts() ([]models.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductsRepository) GetProductsByFilter(offset, limit int, categoryID *uint, priceLessThan *decimal.Decimal) ([]models.Product, int64, error) {
	args := m.Called(offset, limit, categoryID, priceLessThan)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductsRepository) GetProductByCode(code string) (*models.Product, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func TestCatalogHandleGet(t *testing.T) {
	t.Run("returns paginated products with default pagination", func(t *testing.T) {
		mockRepo := new(MockProductsRepository)

		products := []models.Product{
			{
				ID:    1,
				Code:  "PROD001",
				Price: decimal.NewFromFloat(10.99),
				Category: &models.Category{
					Code: "clothing",
					Name: "Clothing",
				},
			},
		}

		mockRepo.On("GetProductsByFilter", 0, 10, mock.MatchedBy(func(v *uint) bool { return v == nil }), mock.MatchedBy(func(d *decimal.Decimal) bool { return d == nil })).Return(products, int64(1), nil)

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Contains(t, recorder.Body.String(), "PROD001")
		assert.Contains(t, recorder.Body.String(), "10.99")
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns products with custom pagination parameters", func(t *testing.T) {
		mockRepo := new(MockProductsRepository)

		products := []models.Product{
			{
				ID:    2,
				Code:  "PROD002",
				Price: decimal.NewFromFloat(12.49),
			},
		}

		mockRepo.On("GetProductsByFilter", 1, 20, mock.MatchedBy(func(v *uint) bool { return v == nil }), mock.MatchedBy(func(d *decimal.Decimal) bool { return d == nil })).Return(products, int64(8), nil)

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog?offset=1&limit=20", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "PROD002")
		assert.Contains(t, recorder.Body.String(), `"total":8`)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns empty products list", func(t *testing.T) {
		mockRepo := new(MockProductsRepository)

		mockRepo.On("GetProductsByFilter", 0, 10, mock.MatchedBy(func(v *uint) bool { return v == nil }), mock.MatchedBy(func(d *decimal.Decimal) bool { return d == nil })).Return([]models.Product{}, int64(0), nil)

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"products":[]`)
		assert.Contains(t, recorder.Body.String(), `"total":0`)
		mockRepo.AssertExpectations(t)
	})
}

func TestCatalogHandleGetByCode(t *testing.T) {
	t.Run("returns product details with variants", func(t *testing.T) {
		mockRepo := new(MockProductsRepository)

		product := &models.Product{
			ID:    1,
			Code:  "PROD001",
			Price: decimal.NewFromFloat(10.99),
			Category: &models.Category{
				Code: "clothing",
				Name: "Clothing",
			},
			Variants: []models.Variant{
				{
					ID:        1,
					ProductID: 1,
					Name:      "Variant A",
					SKU:       "SKU001A",
					Price:     decimal.NewFromFloat(11.99),
				},
				{
					ID:        2,
					ProductID: 1,
					Name:      "Variant B",
					SKU:       "SKU001B",
					Price:     decimal.Zero,
				},
			},
		}

		mockRepo.On("GetProductByCode", "PROD001").Return(product, nil)

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog/PROD001", nil)
		request.SetPathValue("code", "PROD001")

		handler.HandleGetByCode(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "PROD001")
		assert.Contains(t, recorder.Body.String(), "Clothing")
		assert.Contains(t, recorder.Body.String(), "Variant A")
		assert.Contains(t, recorder.Body.String(), "Variant B")
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 404 when product not found", func(t *testing.T) {
		mockRepo := new(MockProductsRepository)
		mockRepo.On("GetProductByCode", "INVALID").Return(nil, assert.AnError)

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog/INVALID", nil)
		request.SetPathValue("code", "INVALID")

		handler.HandleGetByCode(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Product not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("inherits product price for variant without price", func(t *testing.T) {
		mockRepo := new(MockProductsRepository)

		product := &models.Product{
			ID:    1,
			Code:  "PROD001",
			Price: decimal.NewFromFloat(10.99),
			Category: &models.Category{
				Code: "clothing",
				Name: "Clothing",
			},
			Variants: []models.Variant{
				{
					ID:        1,
					ProductID: 1,
					Name:      "Variant Without Price",
					SKU:       "SKU001",
					Price:     decimal.Zero,
				},
			},
		}

		mockRepo.On("GetProductByCode", "PROD001").Return(product, nil)

		handler := NewCatalogHandler(mockRepo)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/catalog/PROD001", nil)
		request.SetPathValue("code", "PROD001")

		handler.HandleGetByCode(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		// The variant should inherit the product price (10.99)
		assert.Contains(t, recorder.Body.String(), "10.99")
		mockRepo.AssertExpectations(t)
	})
}
