package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/surajNirala/students_api/internal/config"
	"github.com/surajNirala/students_api/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	// tableQuery :=
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					name TEXT,
					email TEXT,
					age INTEGER
				)`)
	if err != nil {
		return nil, err
	}
	return &Sqlite{Db: db}, nil

}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("Student not found with id %s", id)

		}
		return types.Student{}, fmt.Errorf("Query Error: %w", err)
	}
	return student, nil
}

// Student List
func (s *Sqlite) GetAllStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return students, nil
}

func (s *Sqlite) UpdateDataStudentById(id int64, name string, email string, age int) (types.Student, error) {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return types.Student{}, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return types.Student{}, fmt.Errorf("error executing update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return types.Student{}, fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return types.Student{}, fmt.Errorf("no student found with id %d", id)
	}

	// Fetch and return the updated student
	var updatedStudent types.Student
	err = s.Db.QueryRow("SELECT id, name, email, age FROM students WHERE id = ?", id).
		Scan(&updatedStudent.Id, &updatedStudent.Name, &updatedStudent.Email, &updatedStudent.Age)
	if err != nil {
		return types.Student{}, fmt.Errorf("error fetching updated student: %w", err)
	}

	return updatedStudent, nil
}

func (s *Sqlite) DeleteDataStudentById(id int64) (string, error) {
	// Fetch the student data before deleting to return it later
	var deletedStudent types.Student
	err := s.Db.QueryRow("SELECT id, name, email, age FROM students WHERE id = ?", id).
		Scan(&deletedStudent.Id, &deletedStudent.Name, &deletedStudent.Email, &deletedStudent.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no student found with id %d", id)
		}
		return "", fmt.Errorf("error fetching student before delete: %w", err)
	}

	// Prepare and execute the DELETE statement
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return "", fmt.Errorf("error preparing delete statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return "", fmt.Errorf("error executing delete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return "", fmt.Errorf("no student found with id %d", id)
	}

	return "Student deleted Successfully.", nil
}
