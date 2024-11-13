package repo

import (
	"database/sql"
	"fmt"
)

type Repo struct {
	db *sql.DB
}
type Subject struct {
	Name string
}
type Class struct {
	Name string
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetUserHash(username string) (string, error) {
	var passwordHash string
	err := r.db.QueryRow(`SELECT password_hash FROM users WHERE username = $1`, username).Scan(&passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("Пользователь не найден: %w", err)
		}
		return "", fmt.Errorf("Ошибка получения хеша: %w", err)
	}
	return passwordHash, nil
}
func (r *Repo) UserExists(username string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("Ошибка проверки существования пользователя: %w", err)
	}
	return exists, nil
}

func (r *Repo) AddUser(user string, pswd string) error {
	fmt.Println("oj")
	_, err := r.db.Exec(`INSERT INTO users (username, password_hash) VALUES ($1, $2)`, user, pswd)
	if err != nil {
		return fmt.Errorf("Ошибка добавления пользователя: %w", err)
	}
	us := "Карлик"
	sub, _ := r.GetSubjects(us)
	class, _ := r.GetClasses(sub[0].Name)
	studentGrades, _ := r.GetStudentsAndGradesByClass(class[0].Name)
	for _, sg := range studentGrades {
		fmt.Printf("Имя: %s %s, Предмет: %s, Оценка: %d\n", sg.FirstName, sg.LastName, sg.SubjectName, sg.Grade)
	}
	fmt.Println(studentGrades, class[0].Name)
	return nil
}
func (r *Repo) GetSubjects(user string) ([]Subject, error) {
	var subjects []Subject
	rows, err := r.db.Query(`SELECT s.name AS subject_name
FROM users u
JOIN teachers t ON u.username = t.user_name
JOIN subjects s ON t.id = s.teacher_id
WHERE u.username = $1`, user)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении предметов: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var subject Subject
		if err := rows.Scan(&subject.Name); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}
		subjects = append(subjects, subject)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении строк: %w", err)
	}

	return subjects, nil

}

type StudentGrade struct {
	FirstName   string
	LastName    string
	SubjectName string
	Grade       int
	Date        string
}

func (r *Repo) GetClasses(subjectName string) ([]Class, error) {
	var classes []Class
	rows, err := r.db.Query(`SELECT DISTINCT sub.class_name
FROM sub
WHERE sub.subject_name = $1`, subjectName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении классов: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var class Class
		if err := rows.Scan(&class.Name); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}
		classes = append(classes, class)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении строк: %w", err)
	}

	return classes, nil
}
func (r *Repo) GetStudentsAndGradesByClass(className string) ([]StudentGrade, error) {
	var studentGrades []StudentGrade
	rows, err := r.db.Query(`
SELECT 
  s.first_name, 
  s.last_name, 
  g.grade
FROM 
  students s 
JOIN 
  classes c ON s.class_id = c.id 
LEFT JOIN 
  grades g ON s.id = g.student_id 
WHERE 
  c.name = $1;
`, className)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении учеников и оценок: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sg StudentGrade
		if err := rows.Scan(&sg.FirstName, &sg.LastName, &sg.Grade); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}
		studentGrades = append(studentGrades, sg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении строк: %w", err)
	}

	return studentGrades, nil
}
func (r *Repo) GetSdudents(className string) ([]StudentGrade, error) {
	var studentGrades []StudentGrade
	rows, err := r.db.Query(`SELECT st.first_name, st.last_name
FROM students st
JOIN classes c ON st.class_id = c.id
WHERE c.name = $1;`, className)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении учеников и оценок: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sg StudentGrade
		if err := rows.Scan(&sg.FirstName, &sg.LastName, &sg.SubjectName, &sg.Grade); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}
		studentGrades = append(studentGrades, sg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении строк: %w", err)
	}

	return studentGrades, nil
}

type Student struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (r *Repo) GetStudentsByClass(className string) ([]Student, error) {
	var students []Student
	rows, err := r.db.Query(`SELECT
	s.id,
		s.last_name,
		s.first_name
	FROM
	students s
	JOIN
	classes c ON s.class_id = c.id
	WHERE
	c.name = $1
	ORDER BY
	s.last_name;`, className)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var student Student
		if err := rows.Scan(&student.ID, &student.FirstName, &student.LastName); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return students, nil
}
