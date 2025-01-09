package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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
	ID        int
	LastName  string
	FirstName string
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
		if err := rows.Scan(&student.ID, &student.LastName, &student.FirstName); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %w", err)
	}

	return students, nil
}

type Quarter struct {
	ID        int       `json:"id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func (r *Repo) GetQuarterByID(quart int) (Quarter, error) {
	var quarter Quarter
	rows, err := r.db.Query(`
 SELECT st.start, st.ID, st.end
FROM quarter st
WHERE st.ID = $1;`, quart)
	if err != nil {
		return Quarter{}, fmt.Errorf("error executing query(GetQuarter): %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return Quarter{}, fmt.Errorf("quarter with ID %d not found", quart)
	}

	if err := rows.Scan(&quarter.StartDate, &quarter.ID, &quarter.EndDate); err != nil {
		return Quarter{}, fmt.Errorf("error scanning row(GetQuarter): %w", err)
	}
	fmt.Println(quarter)
	return quarter, nil
}

type Grade struct {
	Name     string
	LastName string
	Grade    int
	Date     time.Time
	Subject  string
}

func (r *Repo) GetGrades(dateStart, dateEnd time.Time, sub string, class string) ([]Grade, error) {
	var grades []Grade
	rows, err := r.db.Query(`
SELECT
  g.grade,
  s.first_name,
  s.last_name, 
  su.name AS subject_name,
  g.date AS grade_date -- Добавили столбец с датой оценки
FROM
  grades g
JOIN
  students s ON g.student_id = s.id
JOIN
  subjects su ON g.subject_id = su.id
JOIN
  classes c ON s.class_id = c.id
WHERE
  c.name = $1
  AND su.name = $2
  AND g.date BETWEEN $3 AND $4
ORDER BY
  s.last_name, s.first_name, g.date;`, class, sub, dateStart, dateEnd)
	fmt.Println(rows)
	if err != nil {
		return nil, fmt.Errorf("error executing query(GetGrades): %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var grade Grade
		if err := rows.Scan(&grade.Grade, &grade.Name, &grade.LastName, &grade.Subject, &grade.Date); err != nil {
			return nil, fmt.Errorf("error scanning row(GetGrades): %w", err)
		}
		grades = append(grades, grade)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows(GetGrades): %w", err)
	}

	return grades, nil

}

func (r *Repo) UpdateGrades(studentId int, subject string, date time.Time, grade int) error {
	ubId, err := r.db.Query(`SELECT id FROM subjects WHERE name = $1`, subject)
	if err != nil {
		return fmt.Errorf("error executing query(UpdateGrades): %w", err)
	}
	defer ubId.Close()

	var subId int
	if ubId.Next() {
		if err := ubId.Scan(&subId); err != nil {
			return fmt.Errorf("error scanning subject ID: %w", err)
		}
	} else {
		return fmt.Errorf("subject '%s' not found", subject)
	}

	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("начало транзакции: %w", err)
	}
	defer tx.Rollback()

	var count int
	err = tx.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM grades WHERE student_id = $1 AND subject_id = $2 AND date = $3`, studentId, subId, date.Format("2006-01-02")).Scan(&count) //Форматируем дату для сравнения
	if err != nil {
		return fmt.Errorf("проверка существования оценки: %w", err)
	}

	if count > 0 {
		_, err = tx.ExecContext(context.Background(), `UPDATE grades SET grade = $1, created_at = $2 WHERE student_id = $3 AND subject_id = $4 AND date = $5`, grade, time.Now(), studentId, subId, date.Format("2006-01-02")) //Форматируем дату для UPDATE
		if err != nil {
			return fmt.Errorf("обновление оценки: %w", err)
		}
		fmt.Println("Оценка обновлена")
	} else {
		_, err = tx.ExecContext(context.Background(), `INSERT INTO grades (student_id, subject_id, date, grade, created_at) VALUES ($1, $2, $3, $4, $5)`, studentId, subId, date.Format("2006-01-02"), grade, time.Now()) //Форматируем дату для INSERT
		if err != nil {
			return fmt.Errorf("добавление оценки: %w", err)
		}
		fmt.Println("Оценка добавлена")
	}

	return tx.Commit()
}

type StudentGradeAvarage struct {
	IdStudent    int     `json:"id_student"`
	Name         string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	AvarageGrade float64 `json:"avarage_grade"`
}

func (r *Repo) GetАverageGrades(class, subject string, start, end time.Time) ([]StudentGradeAvarage, error) {
	query := `
  SELECT
   s.id AS id_ученика,
   s.first_name AS имя,
   s.last_name AS фамилия,
   COALESCE(AVG(g.grade), 0) AS средний_балл
  FROM
   students s
  JOIN
   classes c ON s.class_id = c.id
  LEFT JOIN
   grades g ON s.id = g.student_id
    AND g.subject_id = (SELECT id FROM subjects WHERE name = $1)
    AND g.date >= $2
    AND g.date <= $3
    AND g.grade <= 5
  JOIN
   subjects sub ON g.subject_id = sub.id
  WHERE
   c.name = $4
  GROUP BY
   s.id,
   s.first_name,
   s.last_name
  ORDER BY
   s.id; --Added ORDER BY clause for consistent results`

	rows, err := r.db.QueryContext(context.Background(), query, subject, start, end, class)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var results []StudentGradeAvarage
	for rows.Next() {
		var result StudentGradeAvarage
		if err := rows.Scan(&result.IdStudent, &result.Name, &result.LastName, &result.AvarageGrade); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

type StudentGradeCount struct {
	FirstName string
	LastName  string
	CountSkip int
}

func (r *Repo) GetCallSkip(class, subject string, start, end time.Time) ([]StudentGradeCount, error) {
	query := `
  SELECT
   s.first_name,
   s.last_name,
   COUNT(g.id) AS count_of_sixes
  FROM
   students s
  JOIN
   grades g ON s.id = g.student_id
  JOIN
   subjects sub ON g.subject_id = sub.id
  JOIN
   classes c ON s.class_id = c.id
  WHERE
   sub.name = $1
   AND c.name = $2
   AND g.grade = 6
   AND g.date BETWEEN $3 AND $4
  GROUP BY
   s.id, s.first_name, s.last_name
  HAVING
   COUNT(g.id) > 0;
 `
	rows, err := r.db.Query(query, subject, class, start, end)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []StudentGradeCount
	for rows.Next() {
		var result StudentGradeCount
		err = rows.Scan(&result.FirstName, &result.LastName, &result.CountSkip)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return results, nil
}
func (r *Repo) UpdGradeQuart(subject string, student, quart, grade int) error {
	now := time.Now()
	ubId, err := r.db.Query(`SELECT id FROM subjects WHERE name = $1`, subject)
	if err != nil {
		return fmt.Errorf("error executing query(UpdateGrades): %w", err)
	}
	defer ubId.Close()
	var subId int
	if ubId.Next() {
		if err := ubId.Scan(&subId); err != nil {
			return fmt.Errorf("error scanning subject ID: %w", err)
		}
	} else {
		return fmt.Errorf("subject '%s' not found", subject)
	}
	//  Используем ON CONFLICT DO UPDATE для вставки или обновления

	_, err = r.db.Exec(`
        INSERT INTO quartgrade (student_id, subject_id, grade, created_at, quartnum) 
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (student_id, subject_id, quartnum) DO UPDATE 
        SET grade = $3, created_at = $4
    `, student, subId, grade, now, quart)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении/обновлении оценки: %w", err)
	}
	return nil
}
func (r *Repo) UpdGetGradeQuart(class, subject string, quart int) ([]struct {
	StudentID int
	FirstName string
	LastName  string
	Grade     int
}, error) {
	rows, err := r.db.Query(`
        SELECT
            s.id AS student_id,
            s.first_name,
            s.last_name,
            qg.grade
        FROM
            students s
        JOIN
            classes c ON s.class_id = c.id
        JOIN
            quartgrade qg ON s.id = qg.student_id
  JOIN
      subjects sub ON qg.subject_id = sub.id
        WHERE
            c.name = $1
            AND sub.name = $2
            AND qg.quartnum = $3;
    `, class, subject, quart)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var grades []struct {
		StudentID int
		FirstName string
		LastName  string
		Grade     int
	}

	for rows.Next() {
		var grade struct {
			StudentID int
			FirstName string
			LastName  string
			Grade     int
		}
		if err := rows.Scan(&grade.StudentID, &grade.FirstName, &grade.LastName, &grade.Grade); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		grades = append(grades, grade)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %w", err)
	}
	fmt.Println(1, grades)
	return grades, nil
}
