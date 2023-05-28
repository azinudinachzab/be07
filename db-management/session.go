package repository



import (

	"a21hc3NpZ25tZW50/model"
	"fmt"



	"gorm.io/gorm"

)



type SessionsRepository interface {

	AddSessions(session model.Session) error

	DeleteSession(token string) error

	UpdateSessions(session model.Session) error

	SessionAvailName(name string) error

	SessionAvailToken(token string) (model.Session, error)

}



type sessionsRepoImpl struct {

	db *gorm.DB

}



func NewSessionRepo(db *gorm.DB) *sessionsRepoImpl {

	return &sessionsRepoImpl{db}

}



func (s *sessionsRepoImpl) AddSessions(session model.Session) error {

	//beginanswer

	if result := s.db.Create(&session); result.Error != nil {

		return gorm.ErrInvalidData

	}



	return nil

	//endanswer return nil

}



func (s *sessionsRepoImpl) DeleteSession(token string) error {

	//beginanswer

	session := model.Session{}

	if result := s.db.Where("token = ?", token).Delete(&session); result.Error != nil {

		return gorm.ErrInvalidData

	}



	return nil

	//endanswer return nil

}



func (s *sessionsRepoImpl) UpdateSessions(session model.Session) error {

	//beginanswer

	if result := s.db.Table("sessions").Where("username = ?", session.Username).Updates(session); result.Error != nil {

		return gorm.ErrInvalidData

	}



	return nil

	//endanswer return nil

}



func (s *sessionsRepoImpl) SessionAvailName(name string) error {

	//beginanswer

	var session model.Session

	if err := s.db.Where("username = ?", name).First(&session).Error; err != nil {

		return fmt.Errorf("Session Not Availability!")

	}



	return nil

	//endanswer return model.Session{}, nil

}



func (s *sessionsRepoImpl) SessionAvailToken(token string) (model.Session, error) {

	//beginanswer

	var session model.Session

	if err := s.db.Where("token = ?", token).First(&session).Error; err != nil {

		return model.Session{}, fmt.Errorf("Session Not Availability!")

	}



	return session, nil

	//endanswer return model.Session{}, nil

}
