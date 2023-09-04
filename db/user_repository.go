package db

import (
	"database/sql"
)

type User struct {
	ID       int64
	Username string
	Password string
}

type UserRepository interface {
	Get(username string) (*User, error)
	GetByID(id int64) (*User, error)
	Add(username string, password string) (id int64, err error)
}

type userRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) UserRepository {
	return userRepository{
		conn: conn,
	}
}

func (u userRepository) Get(username string) (*User, error) {
	row := u.conn.QueryRow("SELECT `id`, `username`, `password` FROM `users` WHERE `username`=?", username)

	user := User{}

	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u userRepository) GetByID(id int64) (*User, error) {
	row := u.conn.QueryRow("SELECT `id`, `username` FROM `users` WHERE `id`=?", id)

	user := User{}

	if err := row.Scan(&user.ID, &user.Username); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u userRepository) Add(username string, password string) (id int64, err error) {
	res, err := u.conn.Exec("INSERT INTO `users` (`username`, `password`) VALUES ($1, $2)", username, password)
	if err != nil {
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return
}
