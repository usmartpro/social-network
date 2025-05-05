package app

type UserDB struct {
	ID         string `db:"id"`
	FirstName  string `db:"first_name"`
	SecondName string `db:"second_name"`
	BirthDate  string `db:"birthdate"`
	Biography  string `db:"biography"`
	City       string `db:"city"`
	Password   string `db:"password"`
}

type PostDB struct {
	ID           string `db:"id"`
	Text         string `db:"text"`
	AuthorUserID string `db:"author_user_id"`
}
