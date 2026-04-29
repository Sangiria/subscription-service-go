package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func (m *MockRepository) Get(id string) (*models.Subscription, error) { return nil, nil }
func (m *MockRepository) List(limit int, offest int) ([]models.Subscription, error) { return nil, nil }
func (m *MockRepository) Delete(id string) error { return nil }
func (m *MockRepository) Update(id string, fields map[string]any) (*models.Subscription, error) { return nil, nil }

func TestCreateSubscription(t *testing.T) {
    e := echo.New()

    for _, tt := range CreateSubscriptionTests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(MockRepository)
            tt.setupMock(mockRepo)
            h := NewSubscriptionHandler(mockRepo)

            body, _ := json.Marshal(tt.input)
            req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
            req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
            rec := httptest.NewRecorder()
            c := e.NewContext(req, rec)

            err := h.CreateSubscription(c)
            if err != nil {
                e.HTTPErrorHandler(c, err)
            }

            assert.Equal(t, tt.expectedStatus, rec.Code)
            assert.Contains(t, rec.Body.String(), tt.expectedBody)

            mockRepo.AssertExpectations(t)
        })
    }
}