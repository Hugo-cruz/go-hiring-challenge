package categories

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/mytheresa/go-hiring-challenge/models"
)

// MockCategoriesRepository is a mock implementation of CategoriesRepository
type MockCategoriesRepository struct {
	mock.Mock
}

func (m *MockCategoriesRepository) GetAll() ([]models.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCategoriesRepository) Create(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoriesRepository) FindByCode(code string) (*models.Category, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func TestCategoriesHandleList(t *testing.T) {
	t.Run("returns all categories", func(t *testing.T) {
		mockRepo := new(MockCategoriesRepository)

		categories := []models.Category{
			{
				ID:   1,
				Code: "clothing",
				Name: "Clothing",
			},
			{
				ID:   2,
				Code: "shoes",
				Name: "Shoes",
			},
			{
				ID:   3,
				Code: "accessories",
				Name: "Accessories",
			},
		}

		mockRepo.On("GetAll").Return(categories, nil)

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/categories", nil)

		handler.HandleList(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Contains(t, recorder.Body.String(), "clothing")
		assert.Contains(t, recorder.Body.String(), "Clothing")
		assert.Contains(t, recorder.Body.String(), "shoes")
		assert.Contains(t, recorder.Body.String(), "Shoes")
		assert.Contains(t, recorder.Body.String(), "accessories")
		assert.Contains(t, recorder.Body.String(), "Accessories")
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns empty categories list", func(t *testing.T) {
		mockRepo := new(MockCategoriesRepository)
		mockRepo.On("GetAll").Return([]models.Category{}, nil)

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/categories", nil)

		handler.HandleList(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Contains(t, recorder.Body.String(), `"categories":[]`)
		mockRepo.AssertExpectations(t)
	})
}

func TestCategoriesHandleCreate(t *testing.T) {
	t.Run("creates a new category", func(t *testing.T) {
		mockRepo := new(MockCategoriesRepository)

		mockRepo.On("Create", mock.MatchedBy(func(c *models.Category) bool {
			return c.Code == "electronics" && c.Name == "Electronics"
		})).Return(nil)

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()

		reqBody := CreateCategoryRequest{
			Code: "electronics",
			Name: "Electronics",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		request := httptest.NewRequest("POST", "/categories", bytes.NewReader(bodyBytes))

		handler.HandleCreate(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Contains(t, recorder.Body.String(), "electronics")
		assert.Contains(t, recorder.Body.String(), "Electronics")
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 400 when code is missing", func(t *testing.T) {
		mockRepo := new(MockCategoriesRepository)

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()

		reqBody := CreateCategoryRequest{
			Name: "Electronics",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		request := httptest.NewRequest("POST", "/categories", bytes.NewReader(bodyBytes))

		handler.HandleCreate(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Code and name are required")
	})

	t.Run("returns 400 when name is missing", func(t *testing.T) {
		mockRepo := new(MockCategoriesRepository)

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()

		reqBody := CreateCategoryRequest{
			Code: "electronics",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		request := httptest.NewRequest("POST", "/categories", bytes.NewReader(bodyBytes))

		handler.HandleCreate(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Code and name are required")
	})

	t.Run("returns 400 when request body is invalid", func(t *testing.T) {
		mockRepo := new(MockCategoriesRepository)

		handler := NewCategoriesHandler(mockRepo)
		recorder := httptest.NewRecorder()

		request := httptest.NewRequest("POST", "/categories", bytes.NewReader([]byte("invalid json")))

		handler.HandleCreate(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Invalid request body")
	})
}
