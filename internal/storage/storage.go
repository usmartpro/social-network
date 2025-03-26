package storage

import (
	"context"
	"fmt"
	"os"
	"social-network/internal/app"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	ctx  context.Context
	conn *pgxpool.Pool
	dsn  string
}

func New(ctx context.Context, dsn string) *Storage {
	return &Storage{
		ctx: ctx,
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) app.Storage {
	conn, err := pgxpool.Connect(ctx, s.dsn)
	if err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "Error connect to database: %v\n", err); err != nil {
			return nil
		}
		os.Exit(1)
	}
	s.conn = conn
	return s
}

func (s *Storage) Close() {
	s.conn.Close()
}

func (s *Storage) RegisterUser(firstName, secondName, birthDate, biography, city, password string) (id *string, err error) {
	sql := `INSERT INTO users (first_name, second_name, birthdate, biography, city, password, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, md5($6), NOW(), NOW()) RETURNING id`
	var rows pgx.Rows
	rows, err = s.conn.Query(s.ctx, sql, firstName, secondName, birthDate, biography, city, password)
	if err != nil && err != pgx.ErrNoRows {
		return nil, app.ErrRegisterUser
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&id)
	if err != nil && err != pgx.ErrNoRows {
		return nil, app.ErrRegisterUser
	}

	return id, err
}

func (s *Storage) GetUser(ID string) (userDB *app.UserDB, err error) {
	var uID, firstName, secondName, birthDate, biography, city string

	sql := `SELECT id, first_name, second_name, TO_CHAR(birthdate, 'YYYY-MM-DD') AS birthdate, biography, city
			FROM users WHERE id = $1 LIMIT 1`
	err = s.conn.QueryRow(s.ctx, sql, ID).Scan(
		&uID,
		&firstName,
		&secondName,
		&birthDate,
		&biography,
		&city,
	)

	if err != nil {
		return nil, err
	}

	userDB = &app.UserDB{
		ID:         uID,
		FirstName:  firstName,
		SecondName: secondName,
		BirthDate:  birthDate,
		Biography:  biography,
		City:       city,
	}
	return
}

func (s *Storage) UserSearch(firstName, lastName string) (usersDB []app.UserDB, err error) {
	var rows pgx.Rows
	sql := `SELECT id, first_name, second_name, TO_CHAR(birthdate, 'YYYY-MM-DD') AS birthdate, biography, city 
			FROM public.users WHERE first_name LIKE $1 AND second_name LIKE $2
			ORDER BY id`

	err = s.conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	rows, err = s.conn.Query(s.ctx, sql, firstName+"%", lastName+"%")
	if err != nil {
		return nil, app.ErrExecQuery
	}
	defer rows.Close()

	for rows.Next() {
		var user app.UserDB
		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.SecondName,
			&user.BirthDate,
			&user.Biography,
			&user.City,
		)
		if err != nil {
			return nil, err
		}
		usersDB = append(usersDB, user)
	}

	if len(usersDB) == 0 {
		return nil, app.ErrObjectNotFound
	}
	return
}
