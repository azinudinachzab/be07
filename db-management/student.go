package repository



import (

	"a21hc3NpZ25tZW50/model"


	"gorm.io/gorm"

)



type StudentRepository interface {

	FetchAll() ([]model.Student, error)

	FetchByID(id int) (*model.Student, error)

	Store(s *model.Student) error

	Update(id int, s *model.Student) error

	Delete(id int) error

	FetchWithClass() (*[]model.StudentClass, error)

}



type studentRepoImpl struct {

	db *gorm.DB

}



func NewStudentRepo(db *gorm.DB) *studentRepoImpl {

	return &studentRepoImpl{db}

}



func (s *studentRepoImpl) FetchAll() ([]model.Student, error) {

	//beginanswer

	var students []model.Student

	err := s.db.Find(&students).Error

	if err != nil {

		return nil, err

	}



	return students, nil

	//endanswer return []model.Student{}, nil

}



func (s *studentRepoImpl) Store(student *model.Student) error {

	//beginanswer

	err := s.db.Create(student).Error

	if err != nil {

		return err

	}



	return nil

	//endanswer return nil

}



func (s *studentRepoImpl) Update(id int, student *model.Student) error {

	//beginanswer

	err := s.db.Model(&model.Student{}).Where("id = ?", id).Updates(map[string]interface{}{

		"name":     student.Name,

		"address":  student.Address,

		"class_id": student.ClassId,

	}).Error

	if err != nil {

		return err

	}



	return nil

	//endanswer return nil

}



func (s *studentRepoImpl) Delete(id int) error {

	//beginanswer

	err := s.db.Delete(&model.Student{}, id).Error

	if err != nil {

		return err

	}



	return nil

	//endanswer return fmt.Errorf("not implement")

}



func (s *studentRepoImpl) FetchByID(id int) (*model.Student, error) {

	//beginanswer

	var student model.Student

	err := s.db.Where("id = ?", id).First(&student).Error

	if err != nil {

		return nil, err

	}



	return &student, nil

	//endanswer return nil, nil

}



func (s *studentRepoImpl) FetchWithClass() (*[]model.StudentClass, error) {

	//beginanswer

	results := []model.StudentClass{}

	s.db.Table("students").

		Select("students.name as name, students.address as address, classes.name as class_name, classes.professor as professor, classes.room_number as room_number").

		Joins("left join classes on classes.id = students.class_id").Scan(&results)

	return &results, nil

	//endanswer return nil, nil

}
