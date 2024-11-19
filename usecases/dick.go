package usecases

import "dnevnik/repo"

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
