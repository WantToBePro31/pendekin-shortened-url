package service

import (
	"pendekin/dto"
	"pendekin/model"
	"pendekin/repository"
)

type UrlService interface {
	GetAllShortenedUrls(userId uint) ([]model.Url, error)
	GetShortenedUrl(shortenedUrl string) (model.Url, error)
	GetShortenedUrlById(urlId uint) (model.Url, error)
	StoreShortenedUrl(pendekin *dto.UrlInsert) error
	EditShortenedUrl(pendekin *dto.UrlUpdate, urlId uint) error
	DeleteShortenedUrl(urlId uint) error
}

type urlService struct {
	urlRepository repository.UrlRepository
}

func InitUrlService(urlRepository repository.UrlRepository) *urlService {
	return &urlService{urlRepository}
}

func (us *urlService) GetAllShortenedUrls(userId uint) ([]model.Url, error) {
	urls, err := us.urlRepository.GetShortenedUrlsByUser(userId)
	if err != nil {
		return []model.Url{}, err
	}
	return urls, nil
}

func (us *urlService) GetShortenedUrl(shortenedUrl string) (model.Url, error) {
	url, err := us.urlRepository.GetByShortenedUrl(shortenedUrl)
	if err != nil {
		return model.Url{}, err
	}
	return url, nil
}

func (us *urlService) GetShortenedUrlById(urlId uint) (model.Url, error) {
	url, err := us.urlRepository.GetById(urlId)
	if err != nil {
		return model.Url{}, err
	}
	return url, nil
}

func (us *urlService) StoreShortenedUrl(pendekin *dto.UrlInsert) error {
	if err := us.urlRepository.InsertShortenedUrl(pendekin); err != nil {
		return err
	}
	return nil
}

func (us *urlService) EditShortenedUrl(pendekin *dto.UrlUpdate, urlId uint) error {
	if err := us.urlRepository.UpdateShortenedUrl(pendekin, urlId); err != nil {
		return err
	}
	return nil
}

func (us *urlService) DeleteShortenedUrl(urlId uint) error {
	if err := us.urlRepository.DeleteShortenedUrl(urlId); err != nil {
		return err
	}
	return nil
}