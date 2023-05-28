package repository



import (

	"a21hc3NpZ25tZW50/model"


	"gorm.io/gorm"

)



type ClassRepository interface {

	FetchAll() ([]model.Class, error)

}



type classRepoImpl struct {

	db *gorm.DB

}



func NewClassRepo(db *gorm.DB) *classRepoImpl {

	return &classRepoImpl{db}

}



func (s *classRepoImpl) FetchAll() ([]model.Class, error) {

	//beginanswer

	var classes []model.Class

	err := s.db.Find(&classes).Error

	if err != nil {

		return nil, err

	}



	return classes, nil

	//endanswer return []model.Class{}, nil

}
