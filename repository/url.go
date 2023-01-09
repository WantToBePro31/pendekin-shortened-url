package repository

import (
	"pendekin/dto"
	"pendekin/model"

	"gorm.io/gorm"
)

type UrlRepository interface {
	GetShortenedUrlsByUser(userId uint) ([]model.Url, error)
	GetByShortenedUrl(shortenedUrl string) (model.Url, error)
	GetById(urlId uint) (model.Url, error)
	InsertShortenedUrl(pendekin *dto.UrlInsert) error
	UpdateShortenedUrl(pendekin *dto.UrlUpdate, urlId uint) error
	DeleteShortenedUrl(urlId uint) error
}

type urlRepository struct {
	db *gorm.DB
}

func InitUrlRepository(db *gorm.DB) *urlRepository {
	return &urlRepository{db}
}

func (ur *urlRepository) GetShortenedUrlsByUser(userId uint) ([]model.Url, error) {
	var pendekins []model.Url
	if err := ur.db.Table("urls").Where("user_id = ?", userId).Find(&pendekins).Error; err != nil {
		return []model.Url{}, err
	}
	return pendekins, nil
}

func (ur *urlRepository) GetByShortenedUrl(shortenedUrl string) (model.Url, error) {
	var pendekin model.Url
	if err := ur.db.Table("urls").Where("shortened_url = ?", shortenedUrl).First(&pendekin).Error; err != nil {
		return model.Url{}, err
	}

	return pendekin, nil
}

func (ur *urlRepository) GetById(urlId uint) (model.Url, error) {
	var pendekin model.Url
	if err := ur.db.Table("urls").Where("id = ?", urlId).First(&pendekin).Error; err != nil {
		return model.Url{}, err
	}

	return pendekin, nil
}

func (ur *urlRepository) InsertShortenedUrl(pendekin *dto.UrlInsert) error {
	if err := ur.db.Table("urls").Create(&pendekin).Error; err != nil {
		return err
	}
	return nil
}

func (ur *urlRepository) UpdateShortenedUrl(pendekin *dto.UrlUpdate, urlId uint) error {
	if err := ur.db.Table("urls").Where("id = ?", urlId).Updates(&pendekin).Error; err != nil {
		return err
	}
	return nil
}

func (ur *urlRepository) DeleteShortenedUrl(urlId uint) error {
	if err := ur.db.Table("urls").Where("id = ?", urlId).Delete(&model.Url{}).Error; err != nil {
		return err
	}
	return nil
}