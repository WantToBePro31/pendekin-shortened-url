package repository

import (
	"pendekin/dto"
	"pendekin/model"

	"gorm.io/gorm"
)

type SessionRepository interface {
	GetSessionByToken(token string) (model.Session, error)
	InsertSession(session *dto.Session) error
	DeleteSession(sessionId uint) error
}

type sessionRepository struct {
	db *gorm.DB
}

func InitSessionRepository(db *gorm.DB) *sessionRepository {
	return &sessionRepository{db}
}

func (sr *sessionRepository) GetSessionByToken(token string) (model.Session, error) {
	var session model.Session
	if err := sr.db.Table("sessions").Where("jwt = ?", token).First(&session).Error; err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func (sr *sessionRepository) InsertSession(session *dto.Session) error {
	if err := sr.db.Table("sessions").Create(&session).Error; err != nil {
		return err
	}
	return nil
}

func (sr *sessionRepository) DeleteSession(sessionId uint) error {
	if err := sr.db.Table("sessions").Where("id = ?", sessionId).Delete(&model.Session{}).Error; err != nil {
		return err
	}
	return nil
}
