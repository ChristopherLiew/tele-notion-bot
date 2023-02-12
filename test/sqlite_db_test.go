package test

import (
	"tele-notion-bot/internal/database"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestInitialiseDB(t *testing.T) {
	database.InitNotionUserDB()
}

func TestGetValidNotionUser(t *testing.T) {
	database.AddNotionUser("ChrisLiew", "abc123")
	user := database.GetNotionUser("ChrisLiew")
	if user.UserName == "" {
		t.Error("Valid user not found")
	}
}

func TestGetNonExistUser(t *testing.T) {
	user := database.GetNotionUser("Hoon Cheep")
	if user.UserName != "" {
		t.Error("User should be invalid")
	}
}
