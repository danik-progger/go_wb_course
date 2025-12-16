package db

import (
	"calendar/entities"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func TestEventsDb(t *testing.T) {
	originalConnectToUser := connectToUser
	defer func() { connectToUser = originalConnectToUser }()

	// Use in-memory sqlite database for testing
	cal, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to connect to in-memory db: %v", err)
	}
	connectToUser = func(u_id entities.Id) (*sqlx.DB, error) {
		return cal, nil
	}

	db, err := InitEventsDb()
	if err != nil {
		t.Fatalf("failed to initialize events db: %v", err)
	}

	userID := entities.Id(1)

	// Create the table
	err = db.AddUser(userID)
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}

	var eventID entities.Id
	t.Run("AddEvent", func(t *testing.T) {
		event := entities.NewEvent(0, time.Now(), time.Now().Add(time.Hour), "Test Event")
		id, err := db.AddEvent(userID, event)
		if err != nil {
			t.Fatalf("AddEvent failed: %v", err)
		}
		if id == 0 {
			t.Fatal("expected a non-zero event id")
		}
		eventID = id
	})

	t.Run("GetEvent", func(t *testing.T) {
		event, err := db.GetEvent(userID, eventID)
		if err != nil {
			t.Fatalf("GetEvent failed: %v", err)
		}
		if event.Id != eventID {
			t.Errorf("expected event id %v, got %v", eventID, event.Id)
		}
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		updatedEvent := entities.NewEvent(eventID, time.Now(), time.Now().Add(time.Hour*2), "Updated Test Event")
		err := db.UpdateEvent(userID, eventID, updatedEvent)
		if err != nil {
			t.Fatalf("UpdateEvent failed: %v", err)
		}

		event, err := db.GetEvent(userID, eventID)
		if err != nil {
			t.Fatalf("GetEvent after update failed: %v", err)
		}
		if event.Description != "Updated Test Event" {
			t.Errorf("expected updated description, got %s", event.Description)
		}
	})

	t.Run("GetEventsInRange", func(t *testing.T) {
		from := time.Now().Add(-time.Minute)
		to := time.Now().Add(time.Hour * 3)
		events, err := db.GetEventsInRange(userID, from, to)
		if err != nil {
			t.Fatalf("GetEventsInRange failed: %v", err)
		}
		if len(events) != 1 {
			t.Fatalf("expected 1 event in range, got %d", len(events))
		}
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		err := db.DeleteEvent(userID, eventID)
		if err != nil {
			t.Fatalf("DeleteEvent failed: %v", err)
		}

		_, err = db.GetEvent(userID, eventID)
		if err == nil {
			t.Fatal("expected an error when getting a deleted event, but got nil")
		}
	})
}
