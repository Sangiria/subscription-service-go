package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"subscription-service-go/internal/models"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(sub *models.Subscription) error {
    args := m.Called(sub)
    return args.Error(0)
}

func (m *MockRepository) Update(id string, fields map[string]any) (*models.Subscription, error) { 
    args := m.Called(id, fields)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockRepository) Sum(sumReq models.SumSubscriptionPrice) (int, error) {
    args := m.Called(sumReq)
    return args.Int(0), args.Error(1) 
}

func (m *MockRepository) Get(id string) (*models.Subscription, error) { return nil, nil }
func (m *MockRepository) List(limit int, offest int) ([]models.Subscription, error) { return nil, nil }
func (m *MockRepository) Delete(id string) error { return nil }

func TestUpdateSubscription(t *testing.T) {
    for _, tt := range UpdateSubscriptionTests {
        t.Run(tt.name, func(t *testing.T) {
            e := echo.New()
            mockRepo := new(MockRepository)
            tt.setupMock(mockRepo)
            h := NewSubscriptionHandler(mockRepo)

            e.PATCH("/subscriptions/:id", h.UpdateSubscriptions)

            body, _ := json.Marshal(tt.input)
            
            targetURL := "/subscriptions/" + tt.paramID 
            
            req := httptest.NewRequest(http.MethodPatch, targetURL, bytes.NewReader(body))
            req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
            rec := httptest.NewRecorder()

            e.ServeHTTP(rec, req)

            assert.Equal(t, tt.expectedStatus, rec.Code)
            assert.Contains(t, rec.Body.String(), tt.expectedBody)
            mockRepo.AssertExpectations(t)
        })
    }
}

func TestSumSubscriptionPrice(t *testing.T) {
    for _, tt := range SumSubscriprionPriceTests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			mockRepo := new(MockRepository)
			tt.setupMock(mockRepo)
			h := NewSubscriptionHandler(mockRepo)
            e.GET("/subscriptions/sum", h.SumSubscriptionsPrice)

			q := make(url.Values)

			if tt.input.UserID != "" {
                q.Set("user_id", tt.input.UserID)
            }
			if tt.input.ServiceName != "" {
				q.Set("service_name", tt.input.ServiceName)
			}
			if tt.input.StartDate != "" {
				q.Set("start_date", tt.input.StartDate)
			}
			if tt.input.EndDate != "" {
				q.Set("end_date", tt.input.EndDate)
			}

			targetURL := fmt.Sprintf("/subscriptions/sum?%s", q.Encode())
			req := httptest.NewRequest(http.MethodGet, targetURL, nil)
			rec := httptest.NewRecorder()

            e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectedBody)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateSubscription(t *testing.T) {
    for _, tt := range CreateSubscriptionTests {
        t.Run(tt.name, func(t *testing.T) {
            e := echo.New()
            mockRepo := new(MockRepository)
            tt.setupMock(mockRepo)
            h := NewSubscriptionHandler(mockRepo)

            e.POST("/subscriptions", h.CreateSubscription)

            body, _ := json.Marshal(tt.input)

            req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewReader(body))
            req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
            rec := httptest.NewRecorder()

            e.ServeHTTP(rec, req)

            assert.Equal(t, tt.expectedStatus, rec.Code)
            assert.Contains(t, rec.Body.String(), tt.expectedBody)
            mockRepo.AssertExpectations(t)
        })
    }
}