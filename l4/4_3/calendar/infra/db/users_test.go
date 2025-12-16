package db

import (
	"calendar/entities"
	"testing"
)

func TestUsersDb(t *testing.T) {
	db := InitUsersDb()

	user1 := entities.NewUser(1, "testuser1")
	user2 := entities.NewUser(2, "testuser2")

	t.Run("Add Users", func(t *testing.T) {
		err := db.AddUser(user1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		err = db.AddUser(user2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("Get User", func(t *testing.T) {
		gotUser, err := db.GetUser(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if gotUser.Id != user1.Id {
			t.Errorf("expected user id %v, got %v", user1.Id, gotUser.Id)
		}
	})

	t.Run("Has User", func(t *testing.T) {
		exists, err := db.HasUser(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !exists {
			t.Errorf("expected user to exist, but it doesn't")
		}

		exists, err = db.HasUser(3)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if exists {
			t.Errorf("expected user to not exist, but it does")
		}
	})

	t.Run("Get Non-existent User", func(t *testing.T) {
		user, err := db.GetUser(99)
		if err != nil {
			t.Fatalf("expected no error for non-existent user, got %v", err)
		}
		if user.Id != 0 {
			t.Errorf("expected user id 0 for non-existent user, got %v", user.Id)
		}
	})
}
