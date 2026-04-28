package repository

import (
	"context"
	"subscription-service-go/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
    Create(sub *models.Subscription) error
    Get(id string) (*models.Subscription, error)
    List(limit int, offest int) []*models.Subscription
}

type subscriptionRepo struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(sub *models.Subscription) error {
    return gorm.G[models.Subscription](r.db).Create(context.Background(), sub)
}

func (r *subscriptionRepo) Get(id string) (*models.Subscription, error) {
    sub, err := gorm.G[models.Subscription](r.db).Where("id = ?", id).First(context.Background())
    if err != nil {
        return nil, err
    }

    return &sub, err
}

func (r *subscriptionRepo) List(limit int, offset int) []*models.Subscription {
    var subs []*models.Subscription
    r.db.Limit(limit).Offset(offset).Find(subs)

    return subs
}