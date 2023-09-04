package handler

import (
	"bytes"
	"errors"
	"market/internal/model"
	"market/internal/service"
	mock_service "market/internal/service/mocks"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehaviour func(ru *mock_service.MockUser, rc *mock_service.MockCart, user model.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            model.User
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username": "testname", "password": "testpassword"}`,
			inputUser: model.User{
				Role:     model.USER,
				Username: "testname",
				Password: "testpassword",
			},
			mockBehaviour: func(ru *mock_service.MockUser, rc *mock_service.MockCart, user model.User) {
				ru.EXPECT().CreateUser(user).Return(1, nil)
				rc.EXPECT().Create(1).Return(1, nil)
				ru.EXPECT().GenerateToken(user.Username, user.Password).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:      "Wrong Input",
			inputBody: `{"username": "testname"}`,
			inputUser: model.User{
				Username: "testname",
			},
			mockBehaviour:        func(r *mock_service.MockUser, rc *mock_service.MockCart, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"username": "testname", "password": "testpassword"}`,
			inputUser: model.User{
				Role:     model.USER,
				Username: "testname",
				Password: "testpassword",
			},
			mockBehaviour: func(r *mock_service.MockUser, rc *mock_service.MockCart, user model.User) {
				r.EXPECT().CreateUser(user).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoUser := mock_service.NewMockUser(c)
			repoCart := mock_service.NewMockCart(c)
			test.mockBehaviour(repoUser, repoCart, test.inputUser)

			services := &service.Service{User: repoUser, Cart: repoCart}

			validate := validator.New()
			model.RegisterCustomValidations(validate) //nolint:errcheck
			logger := zap.NewNop().Sugar()
			h := &Handler{
				Services:  services,
				Logger:    logger,
				Validator: validate,
			}

			r := mux.NewRouter()
			r.HandleFunc("/api/register", h.signUp).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/register",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", appJSON)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehaviour func(r *mock_service.MockUser, inp signInInput)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            signInInput
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username": "testname", "password": "testpassword"}`,
			inputUser: signInInput{
				Username: "testname",
				Password: "testpassword",
			},
			mockBehaviour: func(r *mock_service.MockUser, inp signInInput) {
				r.EXPECT().GenerateToken(inp.Username, inp.Password).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:      "Wrong Input",
			inputBody: `{"username": "testname"}`,
			inputUser: signInInput{
				Username: "testname",
			},
			mockBehaviour:        func(r *mock_service.MockUser, inp signInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"username": "testname", "password": "testpassword"}`,
			inputUser: signInInput{
				Username: "testname",
				Password: "testpassword",
			},
			mockBehaviour: func(r *mock_service.MockUser, inp signInInput) {
				r.EXPECT().GenerateToken(inp.Username, inp.Password).Return("", errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repoUser := mock_service.NewMockUser(c)
			test.mockBehaviour(repoUser, test.inputUser)

			services := &service.Service{User: repoUser}

			validate := validator.New()
			model.RegisterCustomValidations(validate) //nolint:errcheck
			logger := zap.NewNop().Sugar()
			h := &Handler{
				Services:  services,
				Logger:    logger,
				Validator: validate,
			}

			r := mux.NewRouter()
			r.HandleFunc("/api/login", h.signIn).Methods("POST")

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/login",
				bytes.NewBufferString(test.inputBody))
			req.Header.Set("Content-Type", appJSON)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
