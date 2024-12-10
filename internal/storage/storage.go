package storage

import (
	"github.com/surajNirala/students_api/internal/types"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
	UpdateDataStudentById(id int64, name_new string, email_new string, age_new int) (types.Student, error)
	DeleteDataStudentById(id int64) (string, error)
}
