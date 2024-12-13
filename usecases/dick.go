package usecases

import (
	"dnevnik/repo"
	"time"
)

func NewGradeUseCase(repo *repo.Repo) *GradeUseCase {
	return &GradeUseCase{repo: repo}
}

type GradeUseCase struct {
	repo *repo.Repo
}
type PageData struct {
	SubjectName string
	Classes     []repo.Class
	Students    []repo.StudentGrade
}

type Class struct {
	Name string `json:"name"`
}

func (g *GradeUseCase) GetClasses(subjectName string) ([]Class, error) {
	classes, err := g.repo.GetClasses(subjectName)
	if err != nil {
		return nil, err
	}
	result := make([]Class, 0, len(classes))
	for _, class := range classes {
		result = append(result, Class{class.Name})
	}
	return result, nil
}

type Subject struct {
	Name string `json:"name"`
}

func (g *GradeUseCase) GetSubjects(user string) ([]Subject, error) {
	su, err := g.repo.GetSubjects(user)
	if err != nil {
		return nil, err
	}
	result := make([]Subject, 0, len(su))
	for _, suu := range su {
		result = append(result, Subject{suu.Name})
	}
	return result, nil
}
func (g *GradeUseCase) GetStudents(classname string) ([]repo.Student, error) {
	stu, err := g.repo.GetStudentsByClass(classname)
	if err != nil {
		return nil, err
	}
	return stu, nil
}
func (g *GradeUseCase) GetQuarter(id int) (repo.Quarter, error) {
	quar, err := g.repo.GetQuarterByID(id)
	if err != nil {
		return repo.Quarter{}, err
	}
	return quar, nil
}
func (g *GradeUseCase) GetGrades(dateStart, dateEnd time.Time, sub string, class string) ([]repo.Grade, error) {
	grades, err := g.repo.GetGrades(dateStart, dateEnd, class, sub)
	if err != nil {
		return nil, err
	}
	return grades, nil
}

type DateRange struct {
	StartDate time.Time
	EndDate   time.Time
}
type Data struct {
	Students []repo.Student `json:"students"`
	Grades   []repo.Grade   `json:"grades"`
	Dates    DateRange      `json:"dates"`
}

func (g *GradeUseCase) GetGradesTable(subject, className string, quarter int) (*Data, error) {
	students, err := g.repo.GetStudentsByClass(className)
	if err != nil {
		return nil, err
	}
	date, err := g.repo.GetQuarterByID(quarter)
	if err != nil {
		return nil, err
	}
	dates := DateRange{date.StartDate, date.EndDate}
	grades, err := g.repo.GetGrades(dates.StartDate, dates.EndDate, subject, className)
	if err != nil {
		return nil, err
	}

	data := Data{students, grades, dates}

	return &data, nil

}

type Grade struct {
	StudentID int       `json:"StudentID"`
	Date      time.Time `json:"Date"`
	Grade     int       `json:"Grade"`
	Subject   string    `json:"Subject"`
	Time      time.Time `json:"Time"`
}

func (g *GradeUseCase) UpdateGradesBd(studentId, grade int, subject string, date time.Time) error {
	err := g.repo.UpdateGrades(studentId, subject, date, grade)
	if err != nil {
		return err
	}
	return nil
}
