package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"reg_service/model"
	"strconv"
	"time"
)

const (
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "postgres"
)

var db *gorm.DB
var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
	db = Connect()

	router := mux.NewRouter()
	router.HandleFunc("/registration", Registration).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/recover", RecoverPassword).Methods("POST")

	server := http.Server{
		Handler:           router,
		Addr:              "127.0.0.1:3000",
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}

	server.ListenAndServe()
}

func Connect() *gorm.DB {
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", 5432, dbUser, dbPassword, dbName)

	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{})

	if err != nil {
		logger.Println(err)
	}

	err = db.AutoMigrate(&model.User{})

	if err != nil {
		logger.Println(err)
	}

	return db
}

func Registration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	login := r.PostForm.Get("login")
	password := r.PostForm.Get("password")

	user := model.User{}
	err = (&user).GetByLogin(db, login)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.Login != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := (&model.User{Login: login, Password: password}).Create(db)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(`{"id": ` + strconv.Itoa(id) + `}`))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println(err)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	login, password, err := parseLoginAndPass(r)

	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := model.User{}
	err = (&user).GetByLogin(db, login)

	if err != nil || user.Login == "" {
		logger.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if password != user.Password {
		logger.Println("Wrong password for " + login)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(`{"id": ` + strconv.Itoa(user.Id) + `}`))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println(err)
	}
}

func RecoverPassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login := r.PostForm.Get("login")

	user := model.User{}
	err = (&user).GetByLogin(db, login)

	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	msg, err := json.Marshal(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(msg)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println(err)
	}
}

func parseLoginAndPass(r *http.Request) (string, string, error) {
	err := r.ParseForm()

	if err != nil {
		return "", "", err
	}

	login := r.PostForm.Get("login")
	password := r.PostForm.Get("password")

	return login, password, nil
}
