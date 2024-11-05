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
)

type UseCases struct {
	aufUC  *usecases.AufUseCase
	DickUc *usecases.DickUseCase
}

func NewUseCases(aufUC *usecases.AufUseCase, dickUC *usecases.DickUseCase) *UseCases {
	return &UseCases{aufUC: aufUC, DickUc: dickUC}
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
