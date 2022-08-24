package collection_repository

import (
	"time"

	repository_intf "github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) repository_intf.CollectionRepository {
	return &repository{db}
}

func (r *repository) GetMyCollections(user_id uuid.UUID) ([]*entity.Collection, error) {
	datas := []Collection{}
	err := r.db.
		Raw(`
			SELECT * FROM collection
			WHERE author_id = ?
			AND deleted_at = null
		`, user_id).
		Scan(&datas).
		Error
	if err != nil {
		return nil, err
	}
	resp := []*entity.Collection{}
	for _, data := range datas {
		resp = append(resp, data.ToEntity())
	}
	return resp, nil
}

func (r *repository) IsCollectionLikedByUser(id, userId uuid.UUID) (bool, error) {
	metrics := CollectionUserMetrics{}
	err := r.db.
		Table("collection_user_metrics").
		Where("id = ? AND user_id = ?", id, userId).
		Scan(&metrics).
		Error
	if err != nil {
		return false, err
	}
	return metrics.Liked, nil
}

func (r *repository) IsCollectionDislikedByUser(id, userId uuid.UUID) (bool, error) {
	metrics := CollectionUserMetrics{}
	err := r.db.
		Table("collection_user_metrics").
		Where("id = ? AND user_id = ?", id, userId).
		Scan(&metrics).
		Error
	if err != nil {
		return false, err
	}
	return metrics.Disliked, nil
}

func (r *repository) IsCollectionViewedByUser(id, userId uuid.UUID) (bool, error) {
	metrics := CollectionUserMetrics{}
	err := r.db.
		Table("collection_user_metrics").
		Where("id = ? AND user_id = ?", id, userId).
		Scan(&metrics).
		Error
	if err != nil {
		return false, err
	}
	return metrics.Viewed, nil
}

func (r *repository) CollectionLikeInteraction(id, userId uuid.UUID, isLiked bool) error {
	metrics := CollectionMetrics{}
	err := r.db.
		Table("collection_metrics").
		Where("id = ? AND user_id = ?", id, userId).
		Scan(&metrics).
		Error
	if err != nil {
		return err
	}
	if isLiked {
		err = r.db.Table("collection_metrics").Where("collection_id = ?", id).
			Updates(map[string]interface{}{
				"likes": metrics.Likes - 1,
			}).Error
		if err != nil {
			return err
		}
		err = r.db.
			Table("collection_user_metrics").
			Where("collection_id = ? AND user_id = ?", id, userId).
			Updates(map[string]interface{}{
				"liked": false,
			}).
			Error
		if err != nil {
			return err
		}
	} else {
		err = r.db.Table("collection_metrics").Where("collection_id = ?", id).
			Updates(map[string]interface{}{
				"likes": metrics.Likes + 1,
			}).Error
		err = r.db.
			Table("collection_user_metrics").
			Where("collection_id = ? AND user_id = ?", id, userId).
			Updates(map[string]interface{}{
				"liked": true,
			}).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) CollectionDislikeInteraction(id, userId uuid.UUID, isDisliked bool) error {
	metrics := CollectionMetrics{}
	err := r.db.
		Table("collection_metrics").
		Where("id = ? AND user_id = ?", id, userId).
		Scan(&metrics).
		Error
	if err != nil {
		return err
	}
	if isDisliked {
		err = r.db.Table("collection_metrics").Where("collection_id = ?", id).
			Updates(map[string]interface{}{
				"dislikes": metrics.Dislikes - 1,
			}).Error
		err = r.db.
			Table("collection_user_metrics").
			Where("collection_id = ? AND user_id = ?", id, userId).
			Updates(map[string]interface{}{
				"disliked": false,
			}).
			Error
		if err != nil {
			return err
		}
	} else {
		err = r.db.Table("collection_metrics").Where("collection_id = ?", id).
			Updates(map[string]interface{}{
				"dislikes": metrics.Dislikes + 1,
			}).Error
		err = r.db.
			Table("collection_user_metrics").
			Where("collection_id = ? AND user_id = ?", id, userId).
			Updates(map[string]interface{}{
				"disliked": true,
			}).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) ViewCollection(id, userId uuid.UUID) error {
	metrics := CollectionMetrics{}
	err := r.db.
		Table("collection_metrics").
		Where("id = ? AND user_id = ?", id, userId).
		Scan(&metrics).
		Error
	if err != nil {
		return err
	}
	err = r.db.Table("collection_metrics").Where("collection_id = ?", id).
		Updates(map[string]interface{}{
			"views": metrics.Views + 1,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SearchCollectionByName(name string) error {
	panic("Not implemented")
}
func (r *repository) CreateCollection(collection entity.Collection) (*entity.Collection, error) {
	collection_model := Collection{
		Name:      collection.Name,
		Topics:    collection.Topics,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := r.db.Table("collection").Create(collection_model).Error
	if err != nil {
		return nil, err
	}

	return collection_model.ToEntity(), nil
}
