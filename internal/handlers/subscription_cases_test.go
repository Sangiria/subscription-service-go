package handlers

import (
	"encoding/json"
	"net/http"
	"subscription-service-go/internal/models"
	"subscription-service-go/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

const testUUID = "550e8400-e29b-41d4-a716-446655440000"

var SumSubscriprionPriceTests = []struct {
    name string
    input models.SumSubscriptionPriceParams
    setupMock func(m *MockRepository)
    expectedStatus int
    expectedBody string
}{
    {
        name: "success",
        input: models.SumSubscriptionPriceParams{
            UserID: testUUID,
            ServiceName: "Netflix",
			StartDate:   "01-2026",
			EndDate:     "12-2026",
        },
        setupMock: func(m *MockRepository) {
            m.On("Sum", models.SumSubscriptionPriceParams{
                UserID:      testUUID,
                ServiceName: "Netflix",
                StartDate:   "01-2026",
                EndDate:     "12-2026",
            }).Return(1200, nil)
        },
        expectedStatus: http.StatusOK,
		expectedBody: `{"total":1200}`,
    },
    {
        name: "invalid date format",
		input: models.SumSubscriptionPriceParams{
			UserID: testUUID,
			StartDate: "2026-01",
		},
		setupMock: func(m *MockRepository) {},
		expectedStatus: http.StatusBadRequest,
		expectedBody: "Invalid parameters",
    },
    {
        name: "user_id missing",
		input: models.SumSubscriptionPriceParams{
			ServiceName: "Spotify",
		},
		setupMock: func(m *MockRepository) {},
		expectedStatus: http.StatusBadRequest,
		expectedBody: "Invalid parameters",
    },
    {
        name: "internal server error",
        input: models.SumSubscriptionPriceParams{
            UserID: testUUID,
            ServiceName: "Netflix",
			StartDate: "01-2026",
			EndDate: "12-2026",
        },
        setupMock: func(m *MockRepository) {
            m.On("Sum", mock.Anything).Return(0, gorm.ErrInvalidDB)
        },
        expectedStatus: http.StatusInternalServerError,
        expectedBody: "Error calculating subscription sum price",
    },
}

var UpdateSubscriptionTests = []struct{
    name           	string
	paramID			string
    input          	models.SubscriptionUpdateReq
    setupMock      	func(m *MockRepository)
    expectedStatus 	int
    checkResponse   func(t *testing.T, body []byte)
}{
	{
        name: "success",
        paramID: testUUID,
        input: models.SubscriptionUpdateReq{ServiceName: new("Netflix Premium"), Price: new(500)},
        setupMock: func(m *MockRepository) {
            m.On("Update", testUUID, mock.MatchedBy(func(fields map[string]any) bool {
                return fields["service_name"] == "Netflix Premium" && fields["price"] == 500
            })).Return(&models.Subscription{ServiceName: "Netflix Premium", Price: 500}, nil)
        },
        expectedStatus: http.StatusOK,
        checkResponse: func(t *testing.T, body []byte) {
            var actual models.Subscription
            err := json.Unmarshal(body, &actual)
            assert.NoError(t, err)
            assert.Equal(t, "Netflix Premium", actual.ServiceName)
            assert.Equal(t, 500, actual.Price)
        },
    },
	{
        name: "invalid uuid format",
        paramID: "not-a-uuid",
        input: models.SubscriptionUpdateReq{ServiceName: new("New")},
        setupMock: func(m *MockRepository) {},
        expectedStatus: http.StatusBadRequest,
        checkResponse: func(t *testing.T, body []byte) {
            var actual map[string]string
            json.Unmarshal(body, &actual)
            assert.Equal(t, "Invalid parameters", actual["message"])
        },
    },
    {
        name: "negative price",
        paramID: testUUID,
        input: models.SubscriptionUpdateReq{ServiceName: new("Netflix Premium"), Price: new(-500)},
        setupMock: func(m *MockRepository) {},
        expectedStatus: http.StatusBadRequest,
        checkResponse: func(t *testing.T, body []byte) {
            var actual map[string]string
            json.Unmarshal(body, &actual)
            assert.Equal(t, "Invalid parameters", actual["message"])
        },
    },
    {
        name: "invalid date format",
        paramID: testUUID,
        input: models.SubscriptionUpdateReq{StartDate: new("2026-01")},
        setupMock: func(m *MockRepository) {},
        expectedStatus: http.StatusBadRequest,
        checkResponse: func(t *testing.T, body []byte) {
            var actual map[string]string
            json.Unmarshal(body, &actual)
            assert.Equal(t, "Invalid parameters", actual["message"])
        },
    },
	{
        name: "empty request body",
        paramID: testUUID,
        input: models.SubscriptionUpdateReq{},
        setupMock: func(m *MockRepository) {},
        expectedStatus: http.StatusBadRequest,
        checkResponse: func(t *testing.T, body []byte) {
            var actual map[string]string
            json.Unmarshal(body, &actual)
            assert.Equal(t, "Nothing to update", actual["message"])
        },
    },
	{
        name: "record not found",
        paramID: testUUID,
        input: models.SubscriptionUpdateReq{ServiceName: new("Ghost")},
        setupMock: func(m *MockRepository) {
            m.On("Update", testUUID, mock.Anything).Return(nil, gorm.ErrRecordNotFound)
        },
        expectedStatus: http.StatusNotFound,
        checkResponse: func(t *testing.T, body []byte) {
            var actual map[string]string
            json.Unmarshal(body, &actual)
            assert.Equal(t, "This subscription doesn't exist", actual["message"])
        },
    },
	{
        name: "internal server error",
        paramID: testUUID,
        input: models.SubscriptionUpdateReq{ServiceName: new("Error")},
        setupMock: func(m *MockRepository) {
            m.On("Update", testUUID, mock.Anything).Return(nil, gorm.ErrInvalidDB)
        },
        expectedStatus: http.StatusInternalServerError,
        checkResponse: func(t *testing.T, body []byte) {
            var actual map[string]string
            json.Unmarshal(body, &actual)
            assert.Equal(t, "Error updating subscription record", actual["message"])
        },
    },
}

var CreateSubscriptionTests = []struct {
	name           string
	input          models.SubscriptionCreateReq
	setupMock      func(m *MockRepository)
	expectedStatus int
	checkResponse  func(t *testing.T, body []byte)
}{
	{
		name: "success",
		input: models.SubscriptionCreateReq{
			ServiceName: "Netflix", Price: 300,
			UserId: testUUID, StartDate: "07-2023",
		},
		setupMock: func(m *MockRepository) {
			m.On("Create", mock.MatchedBy(func(sub *models.Subscription) bool {
				return sub.ServiceName == "Netflix" &&
					sub.Price == 300 &&
					sub.UserId == testUUID &&
					sub.StartDate.Equal(*utils.ParseToDate("07-2023"))
			})).Return(nil)
		},
		expectedStatus: http.StatusOK,
		checkResponse: func(t *testing.T, body []byte) {
			var response struct {
				Subscription models.Subscription `json:"subscription"`
			}
			err := json.Unmarshal(body, &response)
			assert.NoError(t, err)

			actual := response.Subscription
			assert.Equal(t, "Netflix", actual.ServiceName)
			assert.Equal(t, 300, actual.Price)
			assert.Equal(t, testUUID, actual.UserId)
			assert.True(t, actual.StartDate.Equal(*utils.ParseToDate("07-2023"))) 
		},
	},
	{
		name: "negative price",
		input: models.SubscriptionCreateReq{
			ServiceName: "Netflix", Price: -1,
			UserId: testUUID, StartDate: "07-2023",
		},
		setupMock: func(m *MockRepository) {},
		expectedStatus: http.StatusBadRequest,
		checkResponse: func(t *testing.T, body []byte) {
			var actual map[string]string
			json.Unmarshal(body, &actual)
			assert.Equal(t, "Invalid parameters", actual["message"])
		},
	},
	{
		name: "empty name",
		input: models.SubscriptionCreateReq{
			ServiceName: "", Price: 300,
			UserId: testUUID, StartDate: "07-2023",
		},
		setupMock: func(m *MockRepository) {},
		expectedStatus: http.StatusBadRequest,
		checkResponse: func(t *testing.T, body []byte) {
			var actual map[string]string
			json.Unmarshal(body, &actual)
			assert.Equal(t, "Invalid parameters", actual["message"])
		},
	},
	{
		name: "invalid uuid format",
		input: models.SubscriptionCreateReq{
			ServiceName: "Netflix", Price: 300,
			UserId: "not-uuid", StartDate: "07-2023",
		},
		setupMock: func(m *MockRepository) {},
		expectedStatus: http.StatusBadRequest,
		checkResponse: func(t *testing.T, body []byte) {
			var actual map[string]string
			json.Unmarshal(body, &actual)
			assert.Equal(t, "Invalid parameters", actual["message"])
		},
	},
	{
		name: "invalid date format",
		input: models.SubscriptionCreateReq{
			ServiceName: "Netflix", Price: 300,
			UserId: testUUID, StartDate: "0791832023",
		},
		setupMock: func(m *MockRepository) {},
		expectedStatus: http.StatusBadRequest,
		checkResponse: func(t *testing.T, body []byte) {
			var actual map[string]string
			json.Unmarshal(body, &actual)
			assert.Equal(t, "Invalid parameters", actual["message"])
		},
	},
	{
		name: "conflict",
		input: models.SubscriptionCreateReq{
			ServiceName: "Netflix", Price: 300,
			UserId: testUUID, StartDate: "07-2023",
		},
		setupMock: func(m *MockRepository) {
			m.On("Create", mock.MatchedBy(func(sub *models.Subscription) bool {
				return sub.ServiceName == "Netflix"
			})).Return(gorm.ErrDuplicatedKey)
		},
		expectedStatus: http.StatusConflict,
		checkResponse: func(t *testing.T, body []byte) {
			var actual map[string]string
			json.Unmarshal(body, &actual)
			assert.Equal(t, "This subscription already exist", actual["message"])
		},
	},
	{
		name: "internal server error",
		input: models.SubscriptionCreateReq{
			ServiceName: "Netflix", Price: 300,
			UserId: testUUID, StartDate: "07-2023",
		},
		setupMock: func(m *MockRepository) {
			m.On("Create", mock.Anything).Return(gorm.ErrInvalidDB)
		},
		expectedStatus: http.StatusInternalServerError,
		checkResponse: func(t *testing.T, body []byte) {
			var actual map[string]string
			json.Unmarshal(body, &actual)
			assert.Equal(t, "Error creating subscription record", actual["message"])
		},
	},
}