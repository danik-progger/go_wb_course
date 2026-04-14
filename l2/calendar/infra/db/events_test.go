package db

import (
	"calendar/entities"
	"testing"
	"time"
)

func TestEventsDb_CRUD(t *testing.T) {
	db, err := InitEventsDb()
	if err != nil {
		t.Fatal(err)
	}

	userID := entities.Id(1)

	// Add
	date := time.Now()
	id, err := db.AddEvent(userID, entities.Event{
		StartingAt:  date,
		EndingAt:    date,
		Description: "Test Event",
	})
	if err != nil {
		t.Fatalf("AddEvent failed: %v", err)
	}
	if id != 1 {
		t.Errorf("expected id 1, got %v", id)
	}

	// Get
	event, err := db.GetEvent(userID, id)
	if err != nil {
		t.Fatalf("GetEvent failed: %v", err)
	}
	if event.Description != "Test Event" {
		t.Errorf("expected 'Test Event', got %v", event.Description)
	}

	// Update
	date2 := date.Add(24 * time.Hour)
	err = db.UpdateEvent(userID, id, entities.Event{
		Id:          id,
		StartingAt:  date2,
		EndingAt:    date2,
		Description: "Updated Event",
	})
	if err != nil {
		t.Fatalf("UpdateEvent failed: %v", err)
	}

	event, err = db.GetEvent(userID, id)
	if err != nil {
		t.Fatalf("GetEvent after update failed: %v", err)
	}
	if event.Description != "Updated Event" {
		t.Errorf("expected 'Updated Event', got %v", event.Description)
	}

	// GetEventsInRange
	events, err := db.GetEventsInRange(userID, date2.Add(-1*time.Hour), date2.Add(1*time.Hour))
	if err != nil {
		t.Fatalf("GetEventsInRange failed: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}

	// Delete
	err = db.DeleteEvent(userID, id)
	if err != nil {
		t.Fatalf("DeleteEvent failed: %v", err)
	}

	events, err = db.GetEventsInRange(userID, date.Add(-1*time.Hour), date.Add(48*time.Hour))
	if err != nil {
		t.Fatalf("GetEventsInRange after delete failed: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events after delete, got %d", len(events))
	}
}
