package server

import (
	"dnevnik/usecases"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type UseCases struct {
	aufUC   *usecases.AufUseCase
	GradeUc *usecases.GradeUseCase
}

func NewUseCases(aufUC *usecases.AufUseCase, dickUC *usecases.GradeUseCase) *UseCases {
	return &UseCases{aufUC: aufUC, GradeUc: dickUC}
}

type Handlers struct {
	useCases *UseCases
}

func NewHandlers(cases *UseCases) *Handlers {
	h := &Handlers{cases}
	return h
}

type Route struct {
	Name          string
	Method        string
	Pattern       string
	HandlerFunc   http.HandlerFunc
	MiddlewareAuf func(handler http.Handler) http.Handler
}

func NewRouter(h *Handlers) *mux.Router {
	var (
		routes = Routes{
			Route{Name: "Home", Method: http.MethodGet, Pattern: "/", HandlerFunc: h.HomePage, MiddlewareAuf: usecases.AuthMiddleware},
			Route{Name: "login", Method: http.MethodPost, Pattern: "/login", HandlerFunc: h.HandleLogin},
			Route{Name: "reg", Method: http.MethodPost, Pattern: "/reg", HandlerFunc: h.HandleRegistration},
			Route{Name: "Homeauf", Method: http.MethodGet, Pattern: "/auf", HandlerFunc: h.Homeauf},
			Route{Name: "t/classes", Method: http.MethodGet, Pattern: "/t/classes", HandlerFunc: h.handleClasses, MiddlewareAuf: usecases.AuthMiddleware},
			Route{Name: "t/sub", Method: http.MethodGet, Pattern: "/t/sub", HandlerFunc: h.handleSubjects, MiddlewareAuf: usecases.AuthMiddleware},
			Route{Name: "Gr", Method: http.MethodGet, Pattern: "/Gr", HandlerFunc: h.HomeGr, MiddlewareAuf: usecases.AuthMiddleware},
			Route{Name: "t/stu", Method: http.MethodGet, Pattern: "/t/stu", HandlerFunc: h.handleStudents, MiddlewareAuf: usecases.AuthMiddleware},
			Route{Name: "t/quart", Method: http.MethodGet, Pattern: "/t/quart", HandlerFunc: h.handleQuarte},
			Route{Name: "t/HomeTest", Method: http.MethodGet, Pattern: "/t/s", HandlerFunc: h.TableTest},
			Route{Name: "t/grades", Method: http.MethodGet, Pattern: "/t/grades", HandlerFunc: h.handleGetGrades},
			Route{Name: "t/gradesTable", Method: http.MethodGet, Pattern: "/t/gradesTable", HandlerFunc: h.handleGetGradesTable},
			Route{Name: "/t/updateGrade", Method: http.MethodPost, Pattern: "/t/updateGrade", HandlerFunc: h.handleUpdateGrade},
			Route{Name: "/t/getAverage", Method: http.MethodGet, Pattern: "/t/getAverage", HandlerFunc: h.handleAverageGrade},
			Route{Name: "Analitic", Method: http.MethodGet, Pattern: "/Analitic", HandlerFunc: h.HomeAnalitic, MiddlewareAuf: usecases.AuthMiddleware},
			Route{Name: "GetSkip", Method: http.MethodGet, Pattern: "/GetSkip", HandlerFunc: h.handleGetSkip},
			Route{Name: "handleUpdGradeQuart", Method: http.MethodPost, Pattern: "/t/updGradeQuart", HandlerFunc: h.handleUpdGradeQuart},
			Route{Name: "handleGetGradeQuart", Method: http.MethodPost, Pattern: "/t/GetGradeQuart", HandlerFunc: h.handleGetGradeQuart},
		}
	)
	router := mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
	//	auf.AuthMiddleware(http.HandlerFunc(auf.HandleProtected)).ServeHTTP(w, r)
	//}).Methods("GET")

	for _, route := range routes {

		if route.MiddlewareAuf != nil {
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(route.MiddlewareAuf(route.HandlerFunc))
		} else {
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(route.HandlerFunc)
		}

	}
	return router
}
func (h *Handlers) Homeauf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	htmlFile := "templates/auf.html"
	html, err := ioutil.ReadFile(htmlFile)
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	w.Write([]byte(html))
}
func (h *Handlers) HomePage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	htmlFile := "templates/homePage.html"
	html, err := ioutil.ReadFile(htmlFile)
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	w.Write([]byte(html))
}

type Routes []Route

func (h *Handlers) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Декодируем JSON
	var user usecases.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.useCases.aufUC.Registr(user.Username, user.Password)
	if err != nil {
		fmt.Println("ошибка при регистрации", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Fprintf(w, "Регистрация успешна!")

}

func (h *Handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user usecases.User
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := h.useCases.aufUC.Login(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	ipAddress := r.RemoteAddr
	hashedIP, err := bcrypt.GenerateFromPassword([]byte(ipAddress), 2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Создаем куку с хешированным IP-адресом
	cookie := http.Cookie{
		Name:   "hashed_ip",
		Value:  string(hashedIP),
		Path:   "/",
		MaxAge: 365 * 24 * 60 * 60, // 1 год
	}

	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, token)
}

func (h *Handlers) handleClasses(w http.ResponseWriter, r *http.Request) {
	subjectName := r.URL.Query().Get("subject")
	if subjectName == "" {
		http.Error(w, "Название предмета обязательно", http.StatusBadRequest)
		return
	}

	data, err := h.useCases.GradeUc.GetClasses(subjectName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения данных: %v", err), http.StatusInternalServerError)
		return
	}

	dataJs, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка js маршал: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(dataJs) //tmpl, err := template.ParseFiles("templates/subject.html")

}

func (h *Handlers) handleSubjects(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("username")
	if userName == "" {
		http.Error(w, "Название предмета обязательно", http.StatusBadRequest)
		return
	}

	data, err := h.useCases.GradeUc.GetSubjects(userName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения данных: %v", err), http.StatusInternalServerError)
		return
	}

	dataJs, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка js маршал: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(dataJs)
}

func (h *Handlers) HomeGr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	htmlFile := "templates/subject.html"
	html, err := ioutil.ReadFile(htmlFile)
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	w.Write([]byte(html))
}
func (h *Handlers) handleStudents(w http.ResponseWriter, r *http.Request) {

	className := r.URL.Query().Get("class") //TODO надо порешать хуй знает какой ключ
	if className == "" {
		http.Error(w, "название класса обязательно", http.StatusBadRequest)
		return
	}

	data, err := h.useCases.GradeUc.GetStudents(className)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения данных: %v", err), http.StatusInternalServerError)
		return
	}

	dataJs, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка js маршал: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(dataJs) //tmpl, err := template.ParseFiles("templates/subject.html")
}
func (h *Handlers) handleQuarte(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	id = "1"
	if id == "" {
		http.Error(w, "номер четверти обязателен", http.StatusBadRequest)
		return
	}
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка конвертации в инт: %v", err), http.StatusBadRequest)
	}
	data, err := h.useCases.GradeUc.GetQuarter(intId)
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка полечения данных(GetQuarter): %v", err), http.StatusBadRequest)
	}
	dataJs, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка js маршал(handleQuarte): %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(dataJs)
}
func (h *Handlers) handleGetGrades(w http.ResponseWriter, r *http.Request) {
	dateStart := "2023-06-01"
	dateEnd := "2024-07-01"
	sub := "Математика"
	class := "10А"
	dtStart, err := time.Parse("2006-01-02", dateStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	dtEnd, err := time.Parse("2006-01-02", dateEnd)
	data, err := h.useCases.GradeUc.GetGrades(dtStart, dtEnd, sub, class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	dataJs, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(dataJs)
}

//	func (h *Handlers) HomeTest(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "text/html")
//		htmlFile := "templates/quart.html"
//		html, err := ioutil.ReadFile(htmlFile)
//		if err != nil {
//			log.Fatalf("Ошибка чтения файла: %v", err)
//		}
//
//		w.Write([]byte(html))
//	}

func (h *Handlers) TableTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	htmlFile := "templates/testTableGrades.html"
	html, err := ioutil.ReadFile(htmlFile)
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	w.Write([]byte(html))
}

func (h *Handlers) handleGetGradesTable(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)

	if err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	subject := query.Get("subject")
	class := query.Get("class")
	quarterStr := query.Get("quarter")

	// проверка параметров
	if subject == "" || class == "" || quarterStr == "" {
		http.Error(w, "Не все параметры указаны", http.StatusBadRequest)
		return
	}

	quarter, err := parseInt(quarterStr)
	if err != nil {
		http.Error(w, "Неверный формат четверти", http.StatusBadRequest)
		return
	}
	data, err := h.useCases.GradeUc.GetGradesTable(subject, class, quarter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Получение оценок из базы данных
	//grades, err := getGradesFromDatabase(subject, class, quarter)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Ошибка при получении оценок: %v", err), http.StatusInternalServerError)
	//	return
	//}
	//
	//// Отправка JSON-ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

// вспомогательная функция для парсинга строки в int
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscan(s, &i)
	return i, err
}

func (h *Handlers) handleUpdateGrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	//todo не передается subject
	var grade struct {
		StudentID int    `json:"studentId"`
		Date      string `json:"date"`
		Grade     int    `json:"grade"`
		Subject   string `json:"subject"`
		Time      string `json:"time"`
	}

	err := json.NewDecoder(r.Body).Decode(&grade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	parsedDate, err := time.Parse("2006-01-02", grade.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Преобразование строки времени в time.Time
	fmt.Println(grade.StudentID, grade.Grade, grade.Subject, parsedDate)
	err = h.useCases.GradeUc.UpdateGradesBd(grade.StudentID, grade.Grade, grade.Subject, parsedDate)
	if err != nil {
		fmt.Println(1, err, 1)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// TODO: Обновление оценки в базе данных
	// Пример (замените на ваш код):
	//db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/database")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//

	// ...Ваш код обновления оценки в базе данных...
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Оценка обновлена"))
}
func (h *Handlers) handleAverageGrade(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	class := query.Get("class")
	subject := query.Get("subject")
	quarterStr := query.Get("quarter")
	quarter, err := parseInt(quarterStr)
	data, err := h.useCases.GradeUc.GetAvarage(class, subject, quarter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dataJs, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(dataJs)
}
func (h *Handlers) HomeAnalitic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	htmlFile := "templates/anali.html"
	html, err := ioutil.ReadFile(htmlFile)
	if err != nil {
		log.Fatalf("Ошибка чтения файла: %v", err)
	}

	w.Write([]byte(html))
}
func (h *Handlers) handleGetSkip(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	class := query.Get("class")
	subject := query.Get("subject")
	quarterStr := query.Get("quarter")

	quarter, err := strconv.Atoi(quarterStr)
	if err != nil {
		http.Error(w, "Неверный формат четверти", http.StatusBadRequest)
		return
	}

	if class == "" || subject == "" || quarterStr == "" {
		http.Error(w, "Необходимо указать класс, предмет и четверть", http.StatusBadRequest)
		return
	}

	data, err := h.useCases.GradeUc.GetSkip(class, subject, quarter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	dataJs, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(dataJs)
}
func (h *Handlers) handleUpdGradeQuart(w http.ResponseWriter, r *http.Request) {
	var gradeData struct {
		StudentID int    `json:"studentId"`
		SubjectId string `json:"subjectId"`
		Grade     int    `json:"grade"`
		Quarter   string `json:"quarter"`
	}

	// Декодируем JSON из тела запроса в структуру gradeData
	err := json.NewDecoder(r.Body).Decode(&gradeData)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Выводим полученные данные для отладки
	fmt.Println("Полученные данные:", gradeData)
	// Декодируем JSON из тела запроса в структуру gradeData
	quart, err := strconv.Atoi(gradeData.Quarter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = h.useCases.GradeUc.UpdGradeQuart(gradeData.SubjectId, gradeData.StudentID, quart, gradeData.Grade)
	if err != nil {
		fmt.Println(1, err, 1)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)

}
func (h *Handlers) handleGetGradeQuart(w http.ResponseWriter, r *http.Request) {
	var gradeData struct {
		Class     string `json:"class"`
		SubjectId string `json:"subject"`
		Quarter   string `json:"quarter"`
	}

	// Декодируем JSON из тела запроса в структуру gradeData
	err := json.NewDecoder(r.Body).Decode(&gradeData)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Выводим полученные данные для отладки
	fmt.Println("Полученные данные:", gradeData)
	// Декодируем JSON из тела запроса в структуру gradeData
	quart, err := strconv.Atoi(gradeData.Quarter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	data, err := h.useCases.GradeUc.GetGradeQuart(gradeData.Class, gradeData.SubjectId, quart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dataJs, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(dataJs)

}
