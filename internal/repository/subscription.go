package repository

import (
	"subscription-service-go/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
    Create(sub *models.Subscription) error
    Get(id string) (*models.Subscription, error)
    List(limit int, offest int) ([]models.Subscription, error)
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

func (r *subscriptionRepo) Get(id string) (*models.Subscription, error) {
    var sub *models.Subscription

    err := r.db.Where("id = ?", id).First(sub).Error
    if err != nil {
        return nil, err
    }

    return sub, err
}

func (r *subscriptionRepo) List(limit int, offset int) ([]models.Subscription, error) {
    var subs []models.Subscription

    err := r.db.Limit(limit).Offset(offset).Find(&subs).Error
    if err != nil {
        return nil, err
    }

    return subs, nil
}