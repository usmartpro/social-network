package app

import (
	"errors"
	"strconv"
	"time"
)

type App struct {
	Logger  Logger
	Storage Storage
	Cache   Cache
}

type Logger interface {
	Error(format string, params ...interface{})
	Info(format string, params ...interface{})
}

type Storage interface {
	RegisterUser(firstName, secondName, birthDate, biography, city, password string) (id *string, err error)
	GetUser(ID string) (userDB *UserDB, err error)
	UserSearch(firstName, lastName string) (usersDB []UserDB, err error)
	PostFeed(ID string, limit, offset int) (postsDB []PostDB, err error)
}

type Cache interface {
	Get(key string) (result []PostDB, exists bool, err error)
	Set(key string, value []PostDB, expiration time.Duration) error
	Clear(key string) error
}

func New(logger Logger, storage Storage, cache Cache) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
		Cache:   cache,
	}
}

var (
	ErrExecQuery      = errors.New("ошибка выполнения запроса")
	ErrObjectNotFound = errors.New("анкета не найдена")
	ErrRegisterUser   = errors.New("get rows error")
)

// RegisterUser ...
func (a *App) RegisterUser(firstName, secondName, birthDate, biography, city, password string) (id *string, err error) {
	if id, err = a.Storage.RegisterUser(firstName, secondName, birthDate, biography, city, password); err != nil {
		a.Logger.Error("Error register user: %s", err)
		return nil, err
	}

	return id, nil
}

// GetUser ...
func (a *App) GetUser(ID string) (userDB *UserDB, err error) {
	userDB, err = a.Storage.GetUser(ID)
	return
}

// UserSearch ...
func (a *App) UserSearch(firstName, lastName string) (usersDB []UserDB, err error) {
	usersDB, err = a.Storage.UserSearch(firstName, lastName)
	return
}

// PostFeed ...
func (a *App) PostFeed(ID string, limit, offset int) (postsDB []PostDB, err error) {
	key := getFeedKey(ID, limit, offset)
	// попвтка получить данные из кеша
	if postsDB, existCache, _ := a.Cache.Get(key); existCache {
		// данные есть, возвращаем
		return postsDB, nil
	}

	// получаем из БД
	postsDB, err = a.Storage.PostFeed(ID, limit, offset)
	if err == nil {
		// кешируем ленту
		_ = a.Cache.Set(key, postsDB, 1*time.Minute)
	}
	return
}

func getFeedKey(ID string, limit, offset int) string {
	return ID + ";" + strconv.Itoa(limit) + ";" + strconv.Itoa(offset)
}
