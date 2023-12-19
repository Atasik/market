package service

import (
	"errors"
	"market/internal/model"
	mock_repository "market/internal/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProductService_Create(t *testing.T) {
	type mockBehaviour func(pr *mock_repository.MockProductRepo, ur *mock_repository.MockUserRepo, product model.Product)

	tests := []struct {
		name          string
		product       model.Product
		mockBehaviour mockBehaviour
		want          int
		wantErr       bool
	}{
		{
			name:    "DB OK",
			product: model.Product{},
			mockBehaviour: func(pr *mock_repository.MockProductRepo, ur *mock_repository.MockUserRepo, product model.Product) {
				pr.EXPECT().Create(product).Return(1, nil)
			},
			want: 1,
		},
		{
			name:    "DB Error",
			product: model.Product{},
			mockBehaviour: func(pr *mock_repository.MockProductRepo, ur *mock_repository.MockUserRepo, product model.Product) {
				pr.EXPECT().Create(product).Return(0, errors.New("something went wrong"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		c := gomock.NewController(t)
		defer c.Finish()

		repoProduct := mock_repository.NewMockProductRepo(c)
		repoUser := mock_repository.NewMockUserRepo(c)
		test.mockBehaviour(repoProduct, repoUser, test.product)

		personService := NewProductService(repoProduct, repoUser)

		got, err := personService.Create(test.product)
		if test.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		}
	}
}

func TestProductService_GetAll(t *testing.T) {
	type mockBehaviour func(pr *mock_repository.MockProductRepo, ur *mock_repository.MockUserRepo, q model.ProductQueryInput)

	tests := []struct {
		name          string
		q             model.ProductQueryInput
		mockBehaviour mockBehaviour
		want          []model.Product
		wantErr       bool
	}{
		{
			name: "DB OK",
			q:    model.ProductQueryInput{},
			mockBehaviour: func(pr *mock_repository.MockProductRepo, ur *mock_repository.MockUserRepo, q model.ProductQueryInput) {
				pr.EXPECT().GetAll(q).Return([]model.Product{
					{ID: 1, UserID: 1, Title: "test", Tags: []model.Tag{{Name: "tag"}}, Category: "test", Amount: 15},
				}, nil)
			},
			want: []model.Product{
				{ID: 1, UserID: 1, Title: "test", Tags: []model.Tag{{Name: "tag"}}, Category: "test", Amount: 15},
			},
		},
		{
			name: "DB Error",
			q:    model.ProductQueryInput{},
			mockBehaviour: func(pr *mock_repository.MockProductRepo, ur *mock_repository.MockUserRepo, q model.ProductQueryInput) {
				pr.EXPECT().GetAll(q).Return([]model.Product{}, errors.New("something went wrong"))
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		c := gomock.NewController(t)
		defer c.Finish()

		repoProduct := mock_repository.NewMockProductRepo(c)
		repoUser := mock_repository.NewMockUserRepo(c)
		test.mockBehaviour(repoProduct, repoUser, test.q)

		personService := NewProductService(repoProduct, repoUser)

		got, err := personService.GetAll(test.q)
		if test.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.want, got)
		}
	}
}

func TestProductService_Delete(t *testing.T) {
	type mockBehaviour func(pr *mock_repository.MockProductRepo, ur *mock_repository.MockUserRepo, userID, productID int)

	tests := []struct {
		name          string
		userID        int
		productID     int
		mockBehaviour mockBehaviour
		wantErr       bool
	}{
		{
			name:      "DB OK",
			userID:    1,
			productID: 1,
			mockBehaviour: func(pr *mock_repository.MockProductRepo, ur *mock_repository.MockUserRepo, userID, productID int) {
				ur.EXPECT().GetUserByID(1).Return(model.User{Role: model.ADMIN}, nil)
				pr.EXPECT().Delete(productID).Return(nil)
			},
		},
		{
			name:      "UserID Error",
			userID:    1,
			productID: 1,
			mockBehaviour: func(pr *mock_repository.MockProductRepo, ur *mock_repository.MockUserRepo, userID, productID int) {
				ur.EXPECT().GetUserByID(1).Return(model.User{}, errors.New("something went wrong"))
			},
			wantErr: true,
		},
		{
			name:      "Delete Error",
			userID:    1,
			productID: 1,
			mockBehaviour: func(pr *mock_repository.MockProductRepo, ur *mock_repository.MockUserRepo, userID, productID int) {
				ur.EXPECT().GetUserByID(1).Return(model.User{Role: model.ADMIN}, nil)
				pr.EXPECT().Delete(productID).Return(errors.New("something went wrong"))
			},
			wantErr: true,
		},
		{
			name:      "Wrong User ID error",
			userID:    1,
			productID: 1,
			mockBehaviour: func(pr *mock_repository.MockProductRepo, ur *mock_repository.MockUserRepo, userID, productID int) {
				ur.EXPECT().GetUserByID(1).Return(model.User{ID: 3}, nil)
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		c := gomock.NewController(t)
		defer c.Finish()

		repoProduct := mock_repository.NewMockProductRepo(c)
		repoUser := mock_repository.NewMockUserRepo(c)
		test.mockBehaviour(repoProduct, repoUser, test.userID, test.productID)

		personService := NewProductService(repoProduct, repoUser)

		err := personService.Delete(test.userID, test.productID)
		if test.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
