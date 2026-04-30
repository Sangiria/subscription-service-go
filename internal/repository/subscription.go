package repository

import (
	"subscription-service-go/internal/models"
	"subscription-service-go/internal/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
    Create(sub *models.Subscription) error
    Get(id string) (*models.Subscription, error)
    List(listReq models.ListParams) ([]models.Subscription, error)
    Delete(id string) error
    Update(id string, fields map[string]any) (*models.Subscription, error)
    Sum (sumReq models.SumSubscriptionPriceParams) (int, error)
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
    var sub models.Subscription

    err := r.db.Where("id = ?", id).First(&sub).Error
    if err != nil {
        return nil, err
    }

    return &sub, err
}

func (r *subscriptionRepo) List(listReq models.ListParams) ([]models.Subscription, error) {
    var (
        subs []models.Subscription
        limit = -1
        offset = -1
        query = r.db.Model(&models.Subscription{})
    )

    if listReq.Limit != nil {
        limit = *listReq.Limit
    }
    if listReq.Offset != nil {
        offset = *listReq.Offset
    }
    if listReq.UserID != nil {
        query = query.Where("user_id = ?", *listReq.UserID)
    }

    err := query.Limit(limit).Offset(offset).Find(&subs).Error
    if err != nil {
        return nil, err
    }

    return subs, nil
}

func (r *subscriptionRepo) Delete(id string) error {
    result := r.db.Delete(models.Subscription{}, "id = ?", id)
    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }

    return nil
}

func (r *subscriptionRepo) Update(id string, fields map[string]any) (*models.Subscription, error) {
    var sub models.Subscription
    result := r.db.Model(&sub).Clauses(clause.Returning{}).Where("id = ?", id).Updates(fields)
    if result.Error != nil {
        return &sub, result.Error
    }

    if result.RowsAffected == 0 {
        return &sub, gorm.ErrRecordNotFound
    }

    return &sub, nil
}

func (r *subscriptionRepo) Sum(sumReq models.SumSubscriptionPriceParams) (int, error) {
    var total int
    query := r.db.Model(&models.Subscription{}).Where("user_id = ?", sumReq.UserID)

    if sumReq.ServiceName != "" {
        query = query.Where("service_name = ?", sumReq.ServiceName)
    }

    startDateTime, endDateTime := utils.ParseToDate(sumReq.StartDate), utils.ParseToDate(sumReq.EndDate)

    if startDateTime != nil {
        query = query.Where("start_date >= ?", *startDateTime)
    }
    if endDateTime != nil {
        query = query.Where("start_date <= ?", *endDateTime)
    }

    if err := query.Select("SUM(price)").Scan(&total).Error; err != nil {
        return 0, err
    }

    return total, nil
}