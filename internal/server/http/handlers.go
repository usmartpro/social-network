package internalhttp

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"social-network/internal/app"
)

type ServerHandlers struct {
	app *app.App
}

var jwtSecretKey = []byte("very-secret-key")
var (
	errUserNotFound = errors.New("пользователь не найден")
	errInvalidData  = errors.New("невалидные данные")
)

func NewServerHandlers(a *app.App) *ServerHandlers {
	return &ServerHandlers{app: a}
}

func (s *ServerHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var t string
	regReq := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&regReq); err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	userDB, err := s.app.GetUser(regReq.ID)
	if err != nil || userDB == nil {
		ResponseError(w, http.StatusInternalServerError, errUserNotFound)
		return
	}

	// Ищем пользователя
	// Если пользователь найден, но у него другой пароль, возвращаем ошибку
	if userDB.Password != GetMD5Hash(regReq.Password) {
		ResponseError(w, http.StatusBadRequest, errInvalidData)
		return
	}

	// Генерируем JWT-токен для пользователя,
	// который он будет использовать в будущих HTTP-запросах

	// Генерируем полезные данные, которые будут храниться в токене
	payload := jwt.MapClaims{
		"sub": userDB.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	// Создаем новый JWT-токен и подписываем его по алгоритму HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err = token.SignedString(jwtSecretKey)
	if err != nil {
		logrus.WithError(err).Error("JWT token signing")
		ResponseError(w, http.StatusInternalServerError, err)
	}

	result := LoginResponse{Token: t}
	responseData, err := json.Marshal(result)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseData)
}

func (s *ServerHandlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userExt UserExtended
	if err := json.NewDecoder(r.Body).Decode(&userExt); err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	id, err := s.app.RegisterUser(userExt.FirstName, userExt.SecondName, userExt.BirthDate, userExt.Biography, userExt.City, userExt.Password)

	if err != nil || id == nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	result := RegistrationResponse{UserID: *id}

	responseData, err := json.Marshal(result)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseData)
}

func (s *ServerHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")

	userDB, err := s.app.GetUser(ID)

	if err != nil || userDB == nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	result := User{
		ID:         userDB.ID,
		FirstName:  userDB.FirstName,
		SecondName: userDB.SecondName,
		BirthDate:  userDB.BirthDate,
		Biography:  userDB.Biography,
		City:       userDB.City,
	}

	responseData, err := json.Marshal(result)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseData)
}

func (s *ServerHandlers) UserSearch(w http.ResponseWriter, r *http.Request) {
	var (
		result  []User
		errCode int
	)
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	usersDB, err := s.app.UserSearch(firstName, lastName)

	if err != nil || usersDB == nil {
		errCode = http.StatusInternalServerError
		if err == app.ErrObjectNotFound {
			errCode = http.StatusBadRequest
		}
		ResponseError(w, errCode, err)
		return
	}

	for _, userDB := range usersDB {
		result = append(result, User{
			ID:         userDB.ID,
			FirstName:  userDB.FirstName,
			SecondName: userDB.SecondName,
			BirthDate:  userDB.BirthDate,
			Biography:  userDB.Biography,
			City:       userDB.City,
		})
	}

	responseData, err := json.Marshal(result)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseData)
}

func ResponseError(w http.ResponseWriter, code int, err error) {
	data, err := json.Marshal(Error{
		false,
		err.Error(),
	})
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Failed to marshall error"))
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(data)
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
