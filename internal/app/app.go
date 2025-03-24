package app

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Error(format string, params ...interface{})
	Info(format string, params ...interface{})
}

type Storage interface {
	RegisterUser(firstName, secondName, birthDate, biography, city, password string) (id *string, err error)
	GetUser(ID string) (userDB *UserDB, err error)
	UserSearch(firstName, lastName string) (usersDB []UserDB, err error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) RegisterUser(firstName, secondName, birthDate, biography, city, password string) (id *string, err error) {
	if id, err = a.Storage.RegisterUser(firstName, secondName, birthDate, biography, city, password); err != nil {
		a.Logger.Error("Error register user: %s", err)
		return nil, err
	}

	return id, nil
}

func (a *App) GetUser(ID string) (userDB *UserDB, err error) {
	userDB, err = a.Storage.GetUser(ID)
	return
}

func (a *App) UserSearch(firstName, lastName string) (usersDB []UserDB, err error) {
	usersDB, err = a.Storage.UserSearch(firstName, lastName)
	return
}
