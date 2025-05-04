package internalhttp

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	"social-network/internal/app"
	v1 "social-network/internal/httperrors/v1"
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

type httpError interface {
	StatusCode() int
	HTTPError() v1.Error
}

// Login ...
func (s *ServerHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var t string
	regReq := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&regReq); err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	userDB, err := s.app.GetUser(regReq.ID)
	if err != nil || userDB == nil {
		responseError(w, http.StatusInternalServerError, errUserNotFound)
		return
	}

	// Ищем пользователя
	// Если пользователь найден, но у него другой пароль, возвращаем ошибку
	if userDB.Password != getMD5Hash(regReq.Password) {
		responseError(w, http.StatusBadRequest, errInvalidData)
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
		responseError(w, http.StatusInternalServerError, err)
	}

	result := LoginResponse{Token: t}
	responseData, err := json.Marshal(result)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseData)
}

// RegisterUser ...
func (s *ServerHandlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userExt UserExtended
	if err := json.NewDecoder(r.Body).Decode(&userExt); err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	id, err := s.app.RegisterUser(userExt.FirstName, userExt.SecondName, userExt.BirthDate, userExt.Biography, userExt.City, userExt.Password)

	if err != nil || id == nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	result := RegistrationResponse{UserID: *id}

	responseData, err := json.Marshal(result)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseData)
}

// GetUser ...
func (s *ServerHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Query().Get("id")

	userDB, err := s.app.GetUser(ID)

	if err != nil || userDB == nil {
		responseError(w, http.StatusInternalServerError, err)
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
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseData)
}

// UserSearch ...
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
		responseError(w, errCode, err)
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
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseData)
}

// PostFeed ...
func (s *ServerHandlers) PostFeed(w http.ResponseWriter, r *http.Request) {
	var (
		result                   []Post
		errCode, vLimit, vOffset int
		vID                      string
		err                      error
		postsDB                  []app.PostDB
	)
	ID := r.URL.Query().Get("id")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	vID, vLimit, vOffset, err = validateFeed(ID, limit, offset)

	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	postsDB, err = s.app.PostFeed(vID, vLimit, vOffset)

	if err != nil {
		errCode = http.StatusInternalServerError
		responseError(w, errCode, err)
		return
	}

	for _, postDB := range postsDB {
		result = append(result, Post{
			ID:           postDB.ID,
			Text:         postDB.Text,
			AuthorUserID: postDB.AuthorUserID,
		})
	}
	responseData, err := json.Marshal(result)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseData)
}

func responseError(w http.ResponseWriter, code int, err error) {
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

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func validateFeed(ID, limit, offset string) (vID string, vLimit, vOffset int, err error) {
	if ID == "" {
		err = validationError(errInvalidData, "Не задан один из обязательных параметров")
		return
	}
	vID = strings.Trim(ID, " ")
	if vLimit, err = strconv.Atoi(limit); err != nil {
		vLimit = 100
		err = nil
	}
	if vOffset, err = strconv.Atoi(offset); err != nil {
		vOffset = 10
		err = nil
	}

	return
}

func validationError(origin error, message string) (err error) {
	if _, ok := origin.(httpError); ok {
		err = origin
		return
	}
	err = v1.NewHTTPError(http.StatusBadRequest, v1.Error{
		Code:    v1.CodeValidationError,
		Message: message,
	})
	return
}
