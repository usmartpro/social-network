package internalhttp

type User struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	BirthDate  string `json:"birthdate"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
}

type UserExtended struct {
	User
	Password string `json:"password"`
}

type Post struct {
	ID           string `json:"id"`
	Text         string `json:"text"`
	AuthorUserID string `json:"author_user_id"`
}

type LoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegistrationResponse struct {
	UserID string `json:"user_id"`
}

type Error struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
