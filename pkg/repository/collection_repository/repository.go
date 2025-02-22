package collection_repository

import (
	"errors"
	"time"

	repositoryIntf "github.com/flash-cards-vocab/backend/app/repository"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) repositoryIntf.CollectionRepository {
	return &repository{db}
}

func (r *repository) GetMyCollections(userId uuid.UUID) ([]*entity.Collection, error) {
	datas := []Collection{}
	err := r.db.
		Table("collection").
		Where("author_id=? AND deleted_at IS NULL", userId).
		Find(&datas).
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

func (r *repository) GetUserCollectionsStatistics(userId uuid.UUID) (*entity.UserCollectionStatistics, error) {
	var collectionsCreated int64
	err := r.db.
		Table("collection").
		Where("author_id=? AND deleted_at IS null", userId).
		Count(&collectionsCreated).
		Error
	if err != nil {
		return nil, err
	}

	return &entity.UserCollectionStatistics{
		CollectionsCreated: uint32(collectionsCreated),
	}, nil
}

func (r *repository) GetTotalCardsInCollection(collectionId uuid.UUID) (int, error) {
	var total int64
	err := r.db.
		Table("card").
		Select("COUNT(1)").
		Joins("INNER JOIN collection_cards AS cc ON cc.card_id = card.id").
		Joins("INNER JOIN collection AS c ON cc.collection_id = c.id").
		Where("c.id = ?", collectionId).
		Where("c.deleted_at IS NULL").
		Where("cc.deleted_at IS NULL").
		Where("card.deleted_at IS NULL").
		Count(&total).
		Error
	if err != nil {
		return 0, err
	}
	return int(total), nil
}

func (r *repository) GetRecommendedCollectionsPreview(userId uuid.UUID, limit, offset int) ([]*entity.Collection, error) {
	datas := []Collection{}
	err := r.db.
		Table("collection").
		Where("author_id <> ? AND deleted_at IS null", userId).
		Limit(limit).
		Offset(offset).
		Find(&datas).
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

func (r *repository) GetLikedCollectionsPreview(userId uuid.UUID) ([]*entity.Collection, error) {
	result := []Collection{}
	err := r.db.
		Table("collection coll").
		Joins("INNER JOIN collection_user_metrics coll_um ON coll_um.collection_id = coll.id").
		Where("author_id <> ? AND coll_um.liked = TRUE AND deleted_at IS null", userId).
		Find(&result).
		Error
	if err != nil {
		return nil, err
	}
	resp := []*entity.Collection{}
	for _, data := range result {
		resp = append(resp, data.ToEntity())
	}
	return resp, nil
}

func (r *repository) GetStarredCollectionsPreview(userId uuid.UUID) ([]*entity.Collection, error) {
	datas := []Collection{}
	err := r.db.
		Table("collection coll").
		Joins("INNER JOIN collection_user_metrics coll_um ON coll_um.collection_id = coll.id").
		Where("author_id <> ? AND coll_um.starred=TRUE AND deleted_at IS null", userId).
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

func (r *repository) StarCollectionById(id, userId uuid.UUID) error {
	metrics := CollectionUserMetrics{}
	err := r.db.
		Table("collection_user_metrics").
		Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
		First(&metrics).
		Error
	if err != nil {
		return err
	}

	if metrics.Starred {
		err = r.db.
			Table("collection_user_metrics").
			Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
			Updates(map[string]interface{}{
				"starred":    false,
				"updated_at": time.Now(),
			}).
			Error
		if err != nil {
			return err
		}
	} else {
		err = r.db.
			Table("collection_user_metrics").
			Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
			Updates(map[string]interface{}{
				"starred":    true,
				"updated_at": time.Now(),
			}).
			Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) CreateCollectionUserMetrics(id, userId uuid.UUID) error {
	collUserMetrics := CollectionUserMetrics{
		Id:           uuid.New(),
		UserId:       userId,
		CollectionId: id,
		Liked:        false,
		Disliked:     false,
		Viewed:       true,
		Starred:      false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := r.db.
		Table("collection_user_metrics").
		Create(collUserMetrics).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CreateCollectionUserProgress(id, userId uuid.UUID) error {
	collUserProgress := CollectionUserProgress{
		Id:           uuid.New(),
		CollectionId: id,
		UserId:       userId,
		Mastered:     0,
		Reviewing:    0,
		Learning:     0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := r.db.
		Table("collection_user_progress").
		Create(collUserProgress).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) IsCollectionLikedOrDislikedByUser(id, userId uuid.UUID) (bool, bool, error) {
	metrics := CollectionUserMetrics{}
	err := r.db.
		Table("collection_user_metrics").
		Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
		First(&metrics).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, false, repositoryIntf.ErrCollectionUserMetricsNotFound
		}
		return false, false, err
	}
	return metrics.Liked, metrics.Disliked, nil
}

func (r *repository) IsCollectionLikedByUser(id, userId uuid.UUID) (bool, error) {
	metrics := CollectionUserMetrics{}
	err := r.db.
		Table("collection_user_metrics").
		Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
		First(&metrics).
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
		Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
		First(&metrics).
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
		Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
		First(&metrics).
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
		Where("collection_id = ? AND deleted_at IS null", id).
		First(&metrics).
		Error
	if err != nil {
		return err
	}
	if isLiked {
		err = r.db.Table("collection_metrics").Where("collection_id = ? AND deleted_at IS null", id).
			Updates(map[string]interface{}{
				"likes":      metrics.Likes - 1,
				"updated_at": time.Now(),
			}).Error
		if err != nil {
			return err
		}
		err = r.db.
			Table("collection_user_metrics").
			Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
			Updates(map[string]interface{}{
				"liked":      false,
				"updated_at": time.Now(),
			}).
			Error
		if err != nil {
			return err
		}
	} else {
		err = r.db.Table("collection_metrics").Where("collection_id = ? AND deleted_at IS null", id).
			Updates(map[string]interface{}{
				"likes":      metrics.Likes + 1,
				"updated_at": time.Now(),
			}).Error
		err = r.db.
			Table("collection_user_metrics").
			Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
			Updates(map[string]interface{}{
				"liked":      true,
				"disliked":   false,
				"updated_at": time.Now(),
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
		Where("collection_id = ? AND deleted_at IS null", id).
		First(&metrics).
		Error
	if err != nil {
		return err
	}
	if isDisliked {
		err = r.db.Table("collection_metrics").Where("collection_id = ? AND deleted_at IS null", id).
			Updates(map[string]interface{}{
				"dislikes":   metrics.Dislikes - 1,
				"updated_at": time.Now(),
			}).Error
		err = r.db.
			Table("collection_user_metrics").
			Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
			Updates(map[string]interface{}{
				"disliked":   false,
				"updated_at": time.Now(),
			}).
			Error
		if err != nil {
			return err
		}
	} else {
		err = r.db.Table("collection_metrics").Where("collection_id = ? AND deleted_at IS null", id).
			Updates(map[string]interface{}{
				"dislikes":   metrics.Dislikes + 1,
				"updated_at": time.Now(),
			}).Error
		err = r.db.
			Table("collection_user_metrics").
			Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
			Updates(map[string]interface{}{
				"disliked":   true,
				"liked":      false,
				"updated_at": time.Now(),
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
		Where("collection_id = ? AND deleted_at IS null", id).
		First(&metrics).
		Error
	if err != nil {
		return err
	}
	err = r.db.Table("collection_metrics").Where("collection_id = ? AND deleted_at IS null", id).
		Updates(map[string]interface{}{
			"views":      metrics.Views + 1,
			"updated_at": time.Now(),
		}).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SearchCollectionByName(search string, userId uuid.UUID) ([]*entity.Collection, error) {

	results := []*Collection{}
	err := r.db.
		Table("collection coll").
		Joins("INNER JOIN collection_metrics cm on coll.id = cm.collection_id").
		Where("lower(coll.name) like lower(?) AND coll.author_id <> ? AND coll.deleted_at IS null", search, userId).
		Order("cm.likes").
		Limit(20).
		Scan(&results).
		Error
	if err != nil {
		return nil, err
	}
	resp := []*entity.Collection{}
	for _, data := range results {
		resp = append(resp, data.ToEntity())
	}
	return resp, nil

}

func (r *repository) CreateCollectionWithCards(collection entity.Collection, cards []*entity.Card) (*entity.Collection, error) {
	tx := r.db.Begin()
	collectionModel := Collection{
		Id:        uuid.New(),
		Name:      collection.Name,
		Topics:    collection.Topics,
		AuthorId:  collection.AuthorId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := r.db.Table("collection").Create(collectionModel).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	userProgress := CollectionUserProgress{
		Id:           uuid.New(),
		CollectionId: collectionModel.Id,
		UserId:       collectionModel.AuthorId,
		Mastered:     0,
		Reviewing:    0,
		Learning:     0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = r.db.Table("collection_user_progress").Create(userProgress).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	collUserMetrics := CollectionUserMetrics{
		Id:           uuid.New(),
		UserId:       collectionModel.AuthorId,
		CollectionId: collectionModel.Id,
		Liked:        false,
		Disliked:     false,
		Viewed:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = r.db.Table("collection_user_metrics").Create(collUserMetrics).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	collectionMetrics := CollectionMetrics{
		Id:           uuid.New(),
		CollectionId: collectionModel.Id,
		Likes:        0,
		Dislikes:     0,
		Views:        1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = r.db.Table("collection_metrics").Create(collectionMetrics).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	cardsModels := []*Card{}
	for _, card := range cards {
		cardsModels = append(cardsModels, &Card{
			Id:         uuid.New(),
			Word:       card.Word,
			ImageUrl:   card.ImageUrl,
			Definition: card.Definition,
			Sentence:   card.Sentence,
			Antonyms:   card.Antonyms,
			AuthorId:   collectionModel.AuthorId,
			Synonyms:   card.Synonyms,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}

	err = r.db.Table("card").Create(cardsModels).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	cardUserProgress := []*CardUserProgress{}
	for _, card := range cardsModels {
		cardUserProgress = append(cardUserProgress, &CardUserProgress{
			Id:            uuid.New(),
			CardId:        card.Id,
			UserId:        collectionModel.AuthorId,
			Status:        entity.CardUserProgressType_None,
			LearningCount: 0,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}
	err = r.db.Table("card_user_progress").Create(cardUserProgress).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	cardMetrics := []*CardMetrics{}
	for _, card := range cardsModels {
		cardMetrics = append(cardMetrics, &CardMetrics{
			Id:        uuid.New(),
			CardId:    card.Id,
			Likes:     0,
			Dislikes:  0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	err = r.db.Table("card_metrics").Create(cardMetrics).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	collectionCards := []*CollectionCards{}
	for _, card := range cardsModels {
		collectionCards = append(collectionCards, &CollectionCards{
			Id:           uuid.New(),
			CardId:       card.Id,
			CollectionId: collectionModel.Id,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}
	err = r.db.Table("collection_cards").Create(collectionCards).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return collectionModel.ToEntity(), nil
}

func (r *repository) GetCollectionMetrics(id uuid.UUID) (*entity.CollectionMetrics, error) {
	metrics := CollectionMetrics{}
	err := r.db.
		Table("collection_metrics").
		Where("collection_id = ? AND deleted_at IS null", id).
		First(&metrics).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entity.CollectionMetrics{
				CollectionId: id,
				Likes:        0,
				Dislikes:     0,
				Views:        0,
			}, repositoryIntf.ErrCollectionMetricsNotFound
		}
		return nil, err
	}
	return metrics.ToEntity(), nil
}

func (r *repository) GetCollectionUserProgress(id, userId uuid.UUID) (*entity.CollectionUserProgress, error) {
	metrics := CollectionUserProgress{}
	err := r.db.
		Table("collection_user_progress").
		Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
		First(&metrics).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entity.CollectionUserProgress{ //returning null struct because user didn't start learning with the collection yet and should see all zeros
				CollectionId: id,
				UserId:       userId,
				Mastered:     0,
				Reviewing:    0,
				Learning:     0,
			}, repositoryIntf.ErrCollectionUserProgressNotFound
		}
		return nil, err
	}
	return metrics.ToEntity(), nil
}

func (r *repository) GetCollectionUserMetrics(id, userId uuid.UUID) (*entity.CollectionUserMetrics, error) {
	metrics := CollectionUserMetrics{}
	err := r.db.
		Table("collection_user_metrics").
		Where("collection_id = ? AND user_id = ? AND deleted_at IS null", id, userId).
		First(&metrics).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &entity.CollectionUserMetrics{
				CollectionId: id,
				Liked:        false,
				Disliked:     false,
				Viewed:       false,
				Starred:      false,
			}, repositoryIntf.ErrCollectionUserMetricsNotFound
		}
		return nil, err
	}
	return metrics.ToEntity(), nil

}

func (r *repository) GetCollection(id uuid.UUID) (*entity.Collection, error) {
	data := &Collection{}
	err := r.db.
		Table("collection").
		Where("id = ? AND deleted_at IS null", id).
		First(&data).
		Error
	if err != nil {
		return nil, err
	}
	return data.ToEntity(), nil
}

func (r *repository) GetCollectionCards(collectionId, userId uuid.UUID, limit, offset int) (*entity.CardForUserPagination, error) {
	cards := []*CardForUser{}
	err := r.db.
		Table("card").
		Select("card.*").
		Joins("INNER JOIN collection_cards ON collection_cards.card_id = card.id").
		Joins("INNER JOIN collection ON collection_cards.collection_id = collection.id").
		Where("collection.id = ? AND card.deleted_at IS NULL AND collection_cards.deleted_at IS NULL AND collection.deleted_at IS NULL", collectionId).
		Limit(limit).
		Offset(offset).
		Find(&cards).
		Error
	if err != nil {
		return nil, err
	}

	var total int64
	err = r.db.
		Table("card").
		Joins("INNER JOIN collection_cards ON collection_cards.card_id = card.id").
		Joins("INNER JOIN collection ON collection_cards.collection_id = collection.id").
		Where("collection.id = ?", collectionId).
		Where("card.deleted_at IS NULL").
		Where("collection_cards.deleted_at IS NULL").
		Where("collection.deleted_at IS NULL").
		Count(&total).
		Error

	res := []*entity.CardForUser{}
	for _, card := range cards {
		res = append(res, &entity.CardForUser{
			Id:         card.Id,
			Word:       card.Word,
			ImageUrl:   card.ImageUrl,
			Definition: card.Definition,
			Sentence:   card.Sentence,
			Antonyms:   card.Antonyms,
			Synonyms:   card.Synonyms,
			Status:     card.Status,
			AuthorId:   card.AuthorId,
		})
	}

	data := &entity.CardForUserPagination{
		CardForUser: CardForUser{}.ToArrayEntity(cards),
		Page:        limit,
		Size:        offset,
		Total:       int(total),
	}
	return data, nil

}

func (r *repository) UpdateCollection(collection entity.Collection) error {
	collectionToUpd := Collection{
		Name:      collection.Name,
		Topics:    collection.Topics,
		UpdatedAt: time.Now(),
	}
	return r.db.
		Table("collection").
		Where("id = ? AND deleted_at is NULL", collection.Id).
		Updates(collectionToUpd).
		Error
}

func (r *repository) SearchCollectionByNameForUnregistered(search string) ([]*entity.Collection, error) {
	datas := []*Collection{}
	err := r.db.
		Table("collection").
		Select("collection.*").
		Joins("INNER JOIN collection_metrics ON collection.id = collection_metrics.collection_id").
		Where("lower(collection.name) LIKE lower(?) AND collection.deleted_at IS NULL", "%"+search+"%").
		Order("collection_metrics.likes").
		Limit(10).
		Find(&datas).
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

func (r *repository) GetCollectionCardsForUnregistered(collectionId uuid.UUID, limit int, offset int) (*entity.CardForUserPagination, error) {
	cards := []*CardForUser{}
	err := r.db.
		Table("card").
		Select("card.*, card_user_progress.status").
		Joins("INNER JOIN collection_cards ON collection_cards.card_id = card.id").
		Joins("INNER JOIN collection ON collection_cards.collection_id = collection.id").
		Joins("INNER JOIN card_user_progress ON card_user_progress.card_id = card.id").
		Where(`collection.id = ? 
			AND card.deleted_at IS null 
			AND collection_cards.deleted_at IS null 
			AND collection.deleted_at IS null 
			AND card_user_progress.deleted_at IS null`, collectionId).
		Limit(limit).
		Offset(offset).
		Find(&cards).
		Error
	if err != nil {
		return nil, err
	}

	var total int64
	err = r.db.
		Table("card").
		Joins("INNER JOIN collection_cards ON collection_cards.card_id = card.id").
		Joins("INNER JOIN collection ON collection_cards.collection_id = collection.id").
		Where("collection.id = ? AND card.deleted_at IS null AND collection_cards.deleted_at IS null AND collection.deleted_at IS null", collectionId).
		Count(&total).
		Error
	if err != nil {
		return nil, err
	}

	res := []*entity.CardForUser{}
	for _, card := range cards {
		res = append(res, &entity.CardForUser{
			Id:         card.Id,
			Word:       card.Word,
			ImageUrl:   card.ImageUrl,
			Definition: card.Definition,
			Sentence:   card.Sentence,
			Antonyms:   card.Antonyms,
			Synonyms:   card.Synonyms,
			Status:     card.Status,
			AuthorId:   card.AuthorId,
		})
	}

	data := &entity.CardForUserPagination{
		CardForUser: CardForUser{}.ToArrayEntity(cards),
		Page:        limit,
		Size:        offset,
		Total:       int(total),
	}
	return data, nil
}

func (r *repository) GetRecommendedCollectionsPreviewForUnregistered(limit, offset int) ([]*entity.Collection, error) {
	datas := []Collection{}
	err := r.db.
		Table("collection").
		Where("deleted_at IS null").
		Limit(limit).
		Offset(offset).
		Find(&datas).
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
