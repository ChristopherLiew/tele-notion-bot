package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type NotionUser struct {
	UserName  string
	Token     string
	Timestamp string
}

func InitNotionUserDB() {

	db := sqlx.MustConnect("sqlite3", "file::memory:?cache=shared")

	//TODO: Might want to salt the keys before inserting into the db
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			username text,
			token text,
			timestamp text
		);
	`
	db.MustExec(schema)
}

func AddNotionUser(teleHandle string, token string) {

	db := sqlx.MustConnect("sqlite3", "file::memory:?cache=shared")

	notionUser := `INSERT INTO users (username, token, timestamp) VALUES (?, ?, ?)`
	authTimeStamp := time.Now().Format("01-02-2023 15:15:15")

	db.MustExec(notionUser, teleHandle, token, authTimeStamp)

}

func GetNotionUser(teleHandle string) (user NotionUser) {

	db := sqlx.MustConnect("sqlite3", "file::memory:?cache=shared")

	user = NotionUser{}
	query := fmt.Sprintf(`
		SELECT 	username,
				token,
				timestamp
		FROM	users
		WHERE	username = '%s'
	`, teleHandle)
	db.QueryRowx(query).StructScan(&user)

	return user
}
