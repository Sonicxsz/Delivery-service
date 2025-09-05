package service_test

import (
	"arabic/internal/dto"
	"arabic/internal/model"
	"arabic/internal/service"
	"arabic/pkg/logger"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockICatalogRepository struct {
	mock.Mock
}

func (m *MockICatalogRepository) FindAll(ctx context.Context) ([]*model.Catalog, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Catalog), args.Error(1)
}
func (m *MockICatalogRepository) Delete(ctx context.Context, id uint) (bool, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(bool), args.Error(1)
}
func (m *MockICatalogRepository) Create(ctx context.Context, category *model.Catalog) (*model.Catalog, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(*model.Catalog), args.Error(1)
}
func (m *MockICatalogRepository) Update(ctx context.Context, query string, values []any) (bool, error) {
	args := m.Called(ctx, query, values)
	return args.Get(0).(bool), args.Error(1)
}
func (m *MockICatalogRepository) FindById(ctx context.Context, id uint) (*model.Catalog, bool, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Catalog), args.Get(1).(bool), args.Error(2)
}

func generateMockCatalogItems() func() *model.Catalog {
	count := 1
	return func() *model.Catalog {
		item := &model.Catalog{
			Id:              uint(count),
			Amount:          150,
			DiscountPercent: 0,
			Price:           150,
			Name:            "Salsa",
			Sku:             "236218361836821",
			CategoryId:      1,
			Description:     "Lorem ipsum ipsum Lorem",
		}
		count++
		return item
	}
}

func TestCatalogService_Create(t *testing.T) {
	logger.Init("Error", "./")
	generator := generateMockCatalogItems()
	mockData := generator()

	tests := []struct {
		name         string
		mockValue    *model.Catalog
		mockError    error
		expectErr    bool
		hasDuplicate bool
	}{
		{
			name:         "Success",
			mockValue:    mockData,
			mockError:    nil,
			expectErr:    false,
			hasDuplicate: false,
		},
		{
			name:         "Error with duplicate",
			mockValue:    mockData,
			mockError:    errors.New("has duplicate sku"),
			hasDuplicate: true,
			expectErr:    true,
		},
		{
			name:         "Error",
			mockValue:    nil,
			mockError:    errors.New("some error"),
			expectErr:    true,
			hasDuplicate: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockICatalogRepository{}
			srv := service.CatalogService{CatalogRepository: mockRepo}

			mockRepo.On("Create", mock.Anything, mock.Anything).Return(tc.mockValue, tc.mockError)

			id, err := srv.Create(context.Background(), &dto.CatalogCreateRequest{
				Name: mockData.Name,
			})

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, uint(0), id)
				if tc.hasDuplicate {
					assert.Contains(t, err.Error(), "exist")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, mockData.Id, id)
			}
		})
	}

}

func TestCatalogService_Update(t *testing.T) {
	logger.Init("Error", "./")
	name := "Hello"
	mockData := &dto.CatalogUpdateRequest{
		Id:   1,
		Name: &name,
	}

	tests := []struct {
		name      string
		mockValue bool
		mockError error
		expectErr bool
	}{
		{
			name:      "Success",
			mockValue: true,
			mockError: nil,
			expectErr: false,
		},
		{
			name:      "Not updated case",
			mockValue: false,
			mockError: nil,
			expectErr: true,
		},
		{
			name:      "Error case",
			mockValue: false,
			mockError: errors.New("some error"),
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockICatalogRepository{}

			srv := service.CatalogService{CatalogRepository: mockRepo}

			mockRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockValue, tc.mockError)

			err := srv.Update(context.Background(), mockData)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}

			mockRepo.AssertCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
			mockRepo.AssertNumberOfCalls(t, "Update", 1)
		})
	}

}

func TestCatalogService_GetAll(t *testing.T) {
	logger.Init("Error", "./")

	generator := generateMockCatalogItems()
	mockData := []*model.Catalog{generator(), generator(), generator()}

	tests := []struct {
		name       string
		mockReturn []*model.Catalog
		mockError  error
		expectErr  bool
	}{
		{
			name:       "success",
			mockReturn: mockData,
			mockError:  nil,
			expectErr:  false,
		},
		{
			name:       "repo error",
			mockReturn: nil,
			mockError:  errors.New("error case"),
			expectErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockICatalogRepository{}
			mockRepo.On("FindAll", mock.Anything).Return(tc.mockReturn, tc.mockError)

			srv := &service.CatalogService{CatalogRepository: mockRepo}
			result, err := srv.GetAll(context.Background())

			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, len(tc.mockReturn))
				assert.Equal(t, tc.mockReturn[0].Id, result[0].Id)
				assert.IsType(t, &dto.CatalogResponse{}, result[0])
			}

			mockRepo.AssertCalled(t, "FindAll", mock.Anything)
		})
	}
}

func TestCatalogService_Delete(t *testing.T) {
	logger.Init("Error", "./")
	generator := generateMockCatalogItems()
	mockData := generator()

	test := []struct {
		name       string
		mockReturn bool
		mockError  error
		expectErr  bool
	}{
		{
			name:       "success",
			mockReturn: true,
			mockError:  nil,
			expectErr:  false,
		},
		{
			name:       "repo error",
			mockReturn: false,
			mockError:  errors.New("error case"),
			expectErr:  true,
		},
		{
			name:       "repo error",
			mockReturn: true,
			mockError:  errors.New("error case"),
			expectErr:  true,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockICatalogRepository{}

			svr := service.CatalogService{
				CatalogRepository: mockRepo,
			}

			mockRepo.Mock.On("Delete", mock.Anything, mock.Anything).Return(tc.mockReturn, tc.mockError)

			err := svr.Delete(context.Background(), mockData.Id)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertCalled(t, "Delete", mock.Anything, mock.Anything)
			mockRepo.AssertNumberOfCalls(t, "Delete", 1)
		})
	}

}

func TestCatalogService_GetById(t *testing.T) {
	logger.Init("Error", "./")

	generator := generateMockCatalogItems()
	mockData := generator()

	test := []struct {
		name       string
		mockReturn *model.Catalog
		mockError  error
		mockOk     bool
		expectErr  bool
	}{
		{
			name:       "success",
			mockReturn: mockData,
			mockError:  nil,
			expectErr:  false,
			mockOk:     true,
		},
		{
			name:       "repo error",
			mockReturn: nil,
			mockError:  errors.New("error case"),
			expectErr:  true,
			mockOk:     false,
		},
		{
			name:       "repo not ok",
			mockReturn: nil,
			mockError:  nil,
			mockOk:     false,
			expectErr:  true,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockICatalogRepository{}
			srv := service.CatalogService{
				CatalogRepository: mockRepo,
			}
			mockRepo.On("FindById", mock.Anything, mock.Anything).Return(tc.mockReturn, tc.mockOk, tc.mockError)

			item, err := srv.GetById(context.Background(), mockData.Id)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, item)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, item.Id, mockData.Id)
				assert.IsType(t, &dto.CatalogResponse{}, item)
			}

			mockRepo.AssertCalled(t, "FindById", mock.Anything, mock.Anything)
			mockRepo.AssertNumberOfCalls(t, "FindById", 1)
		})
	}
}
