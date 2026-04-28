package repository

import (
	"subscription-service-go/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
    Create(sub *models.Subscription) error
}

type subscriptionRepo struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(sub *models.Subscription) error {
    return r.db.Create(sub).Error
}
func (r *subscriptionRepo) Get(id string, sub *models.Subscription) error {
    return nil
}