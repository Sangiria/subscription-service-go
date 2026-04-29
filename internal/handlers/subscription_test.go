package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"subscription-service-go/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(sub *models.Subscription) error {
    args := m.Called(sub)
    return args.Error(0)
}

func (m *MockRepository) Get(id string) (*models.Subscription, error)
func (m *MockRepository) List(limit int, offest int) ([]models.Subscription, error)

func TestCreateSubscription(t *testing.T) {
    e := echo.New()

	testUUID := uuid.New().String()
    
    tests := []struct {
        name           string
        input          models.SubscriptionReq
        setupMock      func(m *MockRepository)
        expectedStatus int
        expectedBody   string
    }{
        {
            name: "success",
            input: models.SubscriptionReq{
                ServiceName: "Netflix", Price: 300, 
                UserId: testUUID, StartDate: "07-2023",
            },
            setupMock: func(m *MockRepository) {
                m.On("Create", mock.Anything).Return(nil)
            },
            expectedStatus: http.StatusOK,
            expectedBody:   "Netflix",
        },
        {
            name: "validation failed - negative price",
            input: models.SubscriptionReq{
                ServiceName: "Netflix", Price: -1, 
                UserId: testUUID, StartDate: "07-2023",
            },
            setupMock:      func(m *MockRepository) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   "Validation failed",
        },
		{
            name: "validation failed - empty name",
            input: models.SubscriptionReq{
                ServiceName: "", Price: 300, 
                UserId: testUUID, StartDate: "07-2023",
            },
            setupMock:      func(m *MockRepository) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   "Validation failed",
        },
		{
            name: "validation failed - invalid uuid",
            input: models.SubscriptionReq{
                ServiceName: "Netflix", Price: 300, 
                UserId: "string", StartDate: "07-2023",
            },
            setupMock:      func(m *MockRepository) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   "Validation failed",
        },
		{
            name: "validation failed - invalid date",
            input: models.SubscriptionReq{
                ServiceName: "Netflix", Price: 300, 
                UserId: testUUID, StartDate: "0791832023",
            },
            setupMock:      func(m *MockRepository) {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   "Validation failed",
        },
        {
            name: "conflict",
            input: models.SubscriptionReq{
                ServiceName: "Netflix", Price: 300, 
                UserId: testUUID, StartDate: "07-2023",
            },
            setupMock: func(m *MockRepository) {
                m.On("Create", mock.Anything).Return(gorm.ErrDuplicatedKey)
            },
            expectedStatus: http.StatusConflict,
            expectedBody:   "Record already exist",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(MockRepository)
            tt.setupMock(mockRepo)
            h := NewSubscriptionHandler(mockRepo)

            body, _ := json.Marshal(tt.input)
            req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
            req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
            rec := httptest.NewRecorder()
            c := e.NewContext(req, rec)

            h.CreateSubscription(c)

            assert.Equal(t, tt.expectedStatus, rec.Code)
            assert.Contains(t, rec.Body.String(), tt.expectedBody)
        })
    }
}