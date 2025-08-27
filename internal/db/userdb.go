package db

import (
    "database/sql"
    "errors"
    "log"
    "os"

    "github.com/lib/pq"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
    ErrLoginUsed   = errors.New("login already used")
    ErrUnknownUser = errors.New("unknown login")
)

type UserDB struct {
    conn *sql.DB
}

func NewUserDB() (*UserDB, error) {
    db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        return nil, err
    }

    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return nil, err
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file:///app/internal/db/migrations",
        "postgres", driver,
    )
    if err != nil {
        return nil, err
    }

    err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Migration failed: %v", err)
    }

    return &UserDB{conn: db}, nil
}

func (udb *UserDB) Insert(login string, password string) error {
	_, err := udb.conn.Exec(`INSERT INTO users (login, password) VALUES ($1, $2)`, login, password)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return ErrLoginUsed
		}
		return err
	}
	return nil
}

func (udb *UserDB) Get(login string) (string, error) {
	var password string
	err := udb.conn.QueryRow(`SELECT password FROM users WHERE login = $1`, login).Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUnknownUser
		}
		return "", err
	}
	return password, nil
}

func (udb *UserDB) Delete(login string) error {
	res, err := udb.conn.Exec(`DELETE FROM users WHERE login = $1`, login)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUnknownUser
	}

	return nil
}