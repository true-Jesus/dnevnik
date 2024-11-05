package server

//
//import (
//	"crypto/md5"
//	"database/sql"
//	"encoding/hex"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"time"
//	_ "time"
//	"github.com/dgrijalva/jwt-go"
//	_ "github.com/lib/pq"
//)
//
//type User struct {
//	Username string `json:"Username"`
//	Password string `json:"Password"`
//}
//
//func getUserHash(db *sql.DB, username string) (string, error) {
//	var passwordHash string
//	err := db.QueryRow(`SELECT password_hash FROM users WHERE username = $1`, username).Scan(&passwordHash)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return "", fmt.Errorf("Пользователь не найден: %w", err)
//		}
//		return "", fmt.Errorf("Ошибка получения хеша: %w", err)
//	}
//	return passwordHash, nil
//}
//func userExists(db *sql.DB, username string) (bool, error) {
//	var exists bool
//	err := db.QueryRow(`SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`, username).Scan(&exists)
//	if err != nil {
//		return false, fmt.Errorf("Ошибка проверки существования пользователя: %w", err)
//	}
//	return exists, nil
//}
//
//func addUser(db *sql.DB, user string, pswd string) error {
//	fmt.Println("oj")
//	_, err := db.Exec(`INSERT INTO users (username, password_hash) VALUES ($1, $2)`, user, pswd)
//	if err != nil {
//		return fmt.Errorf("Ошибка добавления пользователя: %w", err)
//	}
//	return nil
//}
//
//func Registr(w http.ResponseWriter, r *http.Request) {
//
//	db, err := connectToDB()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	// Декодируем JSON
//	var user User
//	err = json.Unmarshal(body, &user)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	exit, err := userExists(db, user.Username)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//	}
//	if exit {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	r.ParseForm()
//	passwordHash := user.Password
//	hash := md5.Sum([]byte(passwordHash))
//	hashString := hex.EncodeToString(hash[:])
//	//if hashString == "8fc4c79ffac466caf56a9375e2ba0dfa" {
//	//	quantityHouse := make([]any, len(h.Bd))
//	//	for v, i := range h.Bd {
//	//		quantityHouse = append(quantityHouse, v, len(i))
//	//	}
//	//	w.Write([]byte(fmt.Sprintf("Привет хозяин, %s", quantityHouse)))
//	//	return
//	//}
//	err = addUser(db, user.Username, hashString)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//	}
//	w.Write([]byte("успешная регистрация"))
//
//}
//var mySigningKey = []byte("your_secret_key")
//
//func Login(w http.ResponseWriter, r *http.Request) {
//
//	db, err := connectToDB()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	defer db.Close()
//
//	var user User
//	body, err := ioutil.ReadAll(r.Body)
//	err = json.Unmarshal(body, &user)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//	}
//	passwordHash := user.Password
//	hash := md5.Sum([]byte(passwordHash))
//	hashString := hex.EncodeToString(hash[:])
//
//	passwordHash, err = getUserHash(db, user.Username)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//	}
//	if passwordHash == hashString {
//		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
//			Issuer:  "your_app_name", // Имя приложения
//			IssuedAt: time.Now().Unix(), // Время выдачи
//			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Время истечения
//			Subject:  user.Username, // Имя пользователя
//		})
//		// Подписываем токен
//		tokenString, err := token.SignedString(mySigningKey)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		fmt.Fprintf(w, "Успешный вход! Токен: %s", tokenString)
//	}else{
//		w.Write([]byte("неверный пароль"))
//
//	}
//
//}
//
//
//
//
//func authMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Получаем токен из заголовка Authorization
//		tokenString := r.Header.Get("Authorization")
//		if tokenString == "" {
//			http.Error(w, "Необходимо войти в систему", http.StatusUnauthorized)
//			return
//		}
//
//		// Проверяем токен
//		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, fmt.Errorf("Неверный метод подписи")
//			}
//			return mySigningKey, nil
//		})
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusUnauthorized)
//			return
//		}
//		if !token.Valid {
//			http.Error(w, "Неверный токен", http.StatusUnauthorized)
//			return
//		}
//
//		// Доступ разрешен
//		next.ServeHTTP(w, r)
//	})
//}
//
//
//
//
//
//
//
//
//
//// Проверяем имя пользователя и пароль
//if user.Username == "testuser" && user.Password == "testpassword" {
//// Успешная аутентификация
//// Генерируем JWT
//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
//Issuer:  "your_app_name", // Имя приложения
//IssuedAt: time.Now().Unix(), // Время выдачи
//ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Время истечения
//Subject:  user.Username, // Имя пользователя
//})
//
//// Подписываем токен
//tokenString, err := token.SignedString(mySigningKey)
//if err != nil {
//http.Error(w, err.Error(), http.StatusInternalServerError)
//return
//}
//
//// Отправляем токен клиенту
//fmt.Fprintf(w, "Успешный вход! Токен: %s", tokenString)
//} else {
//// ... (обработка ошибки)
//}
//}
//func authMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Получаем токен из заголовка Authorization
//		tokenString := r.Header.Get("Authorization")
//		if tokenString == "" {
//			http.Error(w, "Необходимо войти в систему", http.StatusUnauthorized)
//			return
//		}
//
//		// Проверяем токен
//		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, fmt.Errorf("Неверный метод подписи")
//			}
//			return mySigningKey, nil
//		})
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusUnauthorized)
//			return
//		}
//		if !token.Valid {
//			http.Error(w, "Неверный токен", http.StatusUnauthorized)
//			return
//		}
//
//		// Доступ разрешен
//		next.ServeHTTP(w, r)
//	})
//}
