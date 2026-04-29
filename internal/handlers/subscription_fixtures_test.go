package handlers

import (
	"net/http"
	"subscription-service-go/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type CreateTest struct {
	name           string
    input          models.SubscriptionCreateReq
    setupMock      func(m *MockRepository)
    expectedStatus int
    expectedBody   string
}

type UpdateTest struct {
	name           	string
	paramID			string
    input          	models.SubscriptionUpdateReq
    setupMock      	func(m *MockRepository)
    expectedStatus 	int
    expectedBody   	string
}

var testUUID = uuid.New().String()

var UpdateSubscriptionTests = []UpdateTest{
	{
        name: "success update",
        paramID: testUUID,
        input: models.SubscriptionUpdateReq{ServiceName: new("Netflix Premium"), Price: new(500)},
        setupMock: func(m *MockRepository) {
            m.On("Update", testUUID, mock.MatchedBy(func(fields map[string]any) bool {
                return fields["service_name"] == "Netflix Premium" && fields["price"] == 500
            })).Return(&models.Subscription{ServiceName: "Netflix Premium"}, nil)
        },
        expectedStatus: http.StatusOK,
        expectedBody: "Netflix Premium",
    },
	{
        name: "invalid uuid format",
        paramID: "not-a-uuid",
        input: models.SubscriptionUpdateReq{ServiceName: new("New")},
        setupMock: func(m *MockRepository) {},
        expectedStatus: http.StatusBadRequest,
        expectedBody: "Invalid UUID format",
    },
	{
        name: "empty request body",
        paramID: testUUID,
        input: models.SubscriptionUpdateReq{},
        setupMock: func(m *MockRepository) {},
        expectedStatus: http.StatusBadRequest,
        expectedBody: "Update failed",
    },
	{
        name: "record not found",
        paramID: testUUID,
        input: models.SubscriptionUpdateReq{ServiceName: new("Ghost")},
        setupMock: func(m *MockRepository) {
            m.On("Update", testUUID, mock.Anything).Return(nil, gorm.ErrRecordNotFound)
        },
        expectedStatus: http.StatusNotFound,
        expectedBody: "This subscription doesn't exist",
    },
	{
        name: "internal server error",
        paramID: testUUID,
        input: models.SubscriptionUpdateReq{ServiceName: new("Error")},
        setupMock: func(m *MockRepository) {
            m.On("Update", testUUID, mock.Anything).Return(nil, gorm.ErrInvalidDB)
        },
        expectedStatus: http.StatusInternalServerError,
        expectedBody: "Error updating subscription record",
    },
}

var CreateSubscriptionTests = []CreateTest{
	{
        name: "success",
        input: models.SubscriptionCreateReq{
            ServiceName: "Netflix", Price: 300, 
            UserId: testUUID, StartDate: "07-2023",
        },
        setupMock: func(m *MockRepository) {
            m.On("Create", mock.Anything).Return(nil)
        },
        expectedStatus: http.StatusOK,
        expectedBody: "Netflix",
    },
    {
        name: "validation failed - negative price",
        input: models.SubscriptionCreateReq{
            ServiceName: "Netflix", Price: -1, 
            UserId: testUUID, StartDate: "07-2023",
        },
        setupMock: func(m *MockRepository) {},
        expectedStatus: http.StatusBadRequest,
        expectedBody: "Validation failed",
    },
	{
        name: "validation failed - empty name",
        input: models.SubscriptionCreateReq{
            ServiceName: "", Price: 300, 
            UserId: testUUID, StartDate: "07-2023",
        },
        setupMock: func(m *MockRepository) {},
        expectedStatus: http.StatusBadRequest,
        expectedBody: "Validation failed",
    },
	{
        name: "validation failed - invalid uuid",
        input: models.SubscriptionCreateReq{
            ServiceName: "Netflix", Price: 300, 
            UserId: "string", StartDate: "07-2023",
        },
        setupMock: func(m *MockRepository) {},
        expectedStatus: http.StatusBadRequest,
        expectedBody: "Validation failed",
    },
	{
        name: "validation failed - invalid date",
        input: models.SubscriptionCreateReq{
            ServiceName: "Netflix", Price: 300, 
            UserId: testUUID, StartDate: "0791832023",
        },
        setupMock: func(m *MockRepository) {},
        expectedStatus: http.StatusBadRequest,
        expectedBody: "Validation failed",
    },
    {
        name: "conflict",
        input: models.SubscriptionCreateReq{
            ServiceName: "Netflix", Price: 300, 
            UserId: testUUID, StartDate: "07-2023",
        },
        setupMock: func(m *MockRepository) {
            m.On("Create", mock.Anything).Return(gorm.ErrDuplicatedKey)
        },
        expectedStatus: http.StatusConflict,
        expectedBody: "This subscription already exist",
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
        expectedBody: "Couldn't create subscription record",
    },
}