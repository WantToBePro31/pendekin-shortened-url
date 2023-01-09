package service

import (
	"pendekin/dto"
	"pendekin/model"
	"pendekin/repository"
)

type SessionService interface {
	GetSession(token string) (model.Session, error)
	StoreSession(session *dto.Session) error
	RemoveSession(sessionId uint) error
}

type sessionService struct {
	sessionRepository repository.SessionRepository
}

func InitSessionService(sessionRepository repository.SessionRepository) *sessionService {
	return &sessionService{sessionRepository}
}

func (ss *sessionService) GetSession(token string) (model.Session, error) {
	session, err := ss.sessionRepository.GetSessionByToken(token)
	if err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func (ss *sessionService) StoreSession(session *dto.Session) error {
	if err := ss.sessionRepository.InsertSession(session); err != nil {
		return err
	}
	return nil
}

func (ss *sessionService) RemoveSession(sessionId uint) error {
	if err := ss.sessionRepository.DeleteSession(sessionId); err != nil {
		return err
	}
	return nil
}
