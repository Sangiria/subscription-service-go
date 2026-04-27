package repository

import (
	"subscription-service-go/internal/models"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
    Create(sub *models.Subscription) error
}

type subscriptionRepo struct {
    db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
    return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(sub *models.Subscription) error {
    return r.db.Create(sub).Error
}