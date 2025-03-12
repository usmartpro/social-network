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
