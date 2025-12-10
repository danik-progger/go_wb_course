package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Event represents an event in the system
type Event struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Date       time.Time `json:"date"`
	TotalSeats int       `json:"total_seats"`
	Available  int       `json:"available"`
	CreatedAt  time.Time `json:"created_at"`
}

// Booking represents a booking for an event
type Booking struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	UserID    string    `json:"user_id"`
	Confirmed bool      `json:"confirmed"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// EventService handles event and booking operations
type EventService struct {
	db *sql.DB
	// Channel to track bookings that need to be expired
	expirationQueue chan int
}

// NewEventService creates a new event service
func NewEventService(dbPath string) (*EventService, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			date DATETIME NOT NULL,
			total_seats INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bookings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_id INTEGER NOT NULL,
			user_id TEXT NOT NULL,
			confirmed BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY(event_id) REFERENCES events(id)
		)
	`)
	if err != nil {
		return nil, err
	}

	// Create an index on expires_at to optimize querying expired bookings
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_bookings_expires_at ON bookings(expires_at)")
	if err != nil {
		return nil, err
	}

	service := &EventService{
		db: db,
		// Channel to track bookings that need to be expired
		expirationQueue: make(chan int, 1000), // Buffered channel for performance
	}

	// Start the background expiration processor
	go service.runExpirationProcessor()

	return service, nil
}

// runExpirationProcessor runs in the background to process expired bookings
func (s *EventService) runExpirationProcessor() {
	// Check for expired bookings every minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.processExpiredBookings()
		case bookingID := <-s.expirationQueue:
			// Process specific booking expiration (if needed)
			s.expireBooking(bookingID)
		}
	}
}

// processExpiredBookings finds and removes all expired bookings
func (s *EventService) processExpiredBookings() {
	rows, err := s.db.Query("SELECT id FROM bookings WHERE confirmed = 0 AND expires_at < ?", time.Now())
	if err != nil {
		log.Printf("Error querying expired bookings: %v", err)
		return
	}
	defer rows.Close()

	var expiredBookingIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("Error scanning expired booking ID: %v", err)
			continue
		}
		expiredBookingIDs = append(expiredBookingIDs, id)
	}

	// Expire each booking
	for _, id := range expiredBookingIDs {
		s.expireBooking(id)
	}
}

// expireBooking removes a specific booking from the database
func (s *EventService) expireBooking(bookingID int) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Error starting transaction for expiration: %v", err)
		return
	}
	defer tx.Rollback()

	// First get the event ID to update available seats
	var eventID int
	err = tx.QueryRow("SELECT event_id FROM bookings WHERE id = ?", bookingID).Scan(&eventID)
	if err != nil {
		log.Printf("Error getting event ID for booking %d: %v", bookingID, err)
		return
	}

	// Delete the expired booking
	_, err = tx.Exec("DELETE FROM bookings WHERE id = ?", bookingID)
	if err != nil {
		log.Printf("Error deleting expired booking %d: %v", bookingID, err)
		return
	}

	// Update available seats for the event
	_, err = tx.Exec("UPDATE events SET available = available + 1 WHERE id = ?", eventID)
	if err != nil {
		log.Printf("Error updating available seats for event %d: %v", eventID, err)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction for booking expiration: %v", err)
		return
	}

	log.Printf("Expired booking %d", bookingID)
}

// CreateEvent creates a new event
func (s *EventService) CreateEvent(name string, date time.Time, totalSeats int) (*Event, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO events (name, date, total_seats) VALUES (?, ?, ?)", name, date, totalSeats)
	if err != nil {
		return nil, err
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	event := &Event{
		ID:         int(eventID),
		Name:       name,
		Date:       date,
		TotalSeats: totalSeats,
		Available:  totalSeats,
		CreatedAt:  time.Now(),
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return event, nil
}

// GetEvent returns an event by ID with current available seats
func (s *EventService) GetEvent(eventID int) (*Event, error) {
	row := s.db.QueryRow(`
		SELECT id, name, date, total_seats,
		       (total_seats - COALESCE(booking_count, 0)) as available
		FROM events
		LEFT JOIN (
			SELECT event_id, COUNT(*) as booking_count
			FROM bookings
			WHERE confirmed = 1
			GROUP BY event_id
		) b ON events.id = b.event_id
		WHERE events.id = ?
	`, eventID)

	var event Event
	var totalSeats int
	err := row.Scan(&event.ID, &event.Name, &event.Date, &totalSeats, &event.Available)
	if err != nil {
		return nil, err
	}

	event.TotalSeats = totalSeats

	return &event, nil
}

// BookSeat creates a new booking for an event
func (s *EventService) BookSeat(eventID int, userID string) (*Booking, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check if event has available seats (not including unconfirmed bookings)
	var available int
	err = tx.QueryRow(`
		SELECT total_seats - COALESCE(confirmed_count, 0) - COALESCE(unconfirmed_count, 0) as available
		FROM events
		LEFT JOIN (
			SELECT event_id, COUNT(*) as confirmed_count
			FROM bookings
			WHERE confirmed = 1
			GROUP BY event_id
		) c ON events.id = c.event_id
		LEFT JOIN (
			SELECT event_id, COUNT(*) as unconfirmed_count
			FROM bookings
			WHERE confirmed = 0
			GROUP BY event_id
		) uc ON events.id = uc.event_id
		WHERE events.id = ?
	`, eventID).Scan(&available)

	if err != nil {
		return nil, err
	}

	if available <= 0 {
		return nil, fmt.Errorf("no available seats")
	}

	// Create booking with expiration time (30 minutes from now)
	expiresAt := time.Now().Add(30 * time.Minute)
	result, err := tx.Exec("INSERT INTO bookings (event_id, user_id, expires_at) VALUES (?, ?, ?)", eventID, userID, expiresAt)
	if err != nil {
		return nil, err
	}

	bookingID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	booking := &Booking{
		ID:        int(bookingID),
		EventID:   eventID,
		UserID:    userID,
		Confirmed: false,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return booking, nil
}

// ConfirmBooking confirms a booking, making it permanent
func (s *EventService) ConfirmBooking(bookingID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update booking to confirmed
	result, err := tx.Exec("UPDATE bookings SET confirmed = 1 WHERE id = ?", bookingID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking not found")
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetAllEvents returns all events with their available seats
func (s *EventService) GetAllEvents() ([]Event, error) {
	rows, err := s.db.Query(`
		SELECT id, name, date, total_seats,
		       (total_seats - COALESCE(booking_count, 0)) as available
		FROM events
		LEFT JOIN (
			SELECT event_id, COUNT(*) as booking_count
			FROM bookings
			WHERE confirmed = 1
			GROUP BY event_id
		) b ON events.id = b.event_id
		ORDER BY date
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		var totalSeats int
		err := rows.Scan(&event.ID, &event.Name, &event.Date, &totalSeats, &event.Available)
		if err != nil {
			return nil, err
		}
		event.TotalSeats = totalSeats
		events = append(events, event)
	}

	return events, nil
}

// HTTP Handlers
func (s *EventService) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestBody struct {
		Name       string    `json:"name"`
		DateString string    `json:"date"`
		Date       time.Time `json:"-"` // We'll parse from DateString
		TotalSeats int       `json:"total_seats"`
	}

	var req RequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parse date from string
	parsedDate, err := time.Parse(time.RFC3339, req.DateString)
	if err != nil {
		http.Error(w, "Invalid date format, use RFC3339", http.StatusBadRequest)
		return
	}

	req.Date = parsedDate

	event, err := s.CreateEvent(req.Name, req.Date, req.TotalSeats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (s *EventService) handleGetEvent(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.URL.Path[len("/events/"):]
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	event, err := s.GetEvent(eventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (s *EventService) handleBookEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract event ID from URL path
	path := r.URL.Path
	// Strip the leading "/" and split the path into parts
	path = strings.TrimPrefix(path, "/")
	pathParts := strings.Split(path, "/")

	// Find the event ID (it should be after 'events' in the path)
	eventID := -1
	for idx, part := range pathParts {
		if part == "events" && idx+1 < len(pathParts) {
			id, err := strconv.Atoi(pathParts[idx+1])
			if err != nil {
				http.Error(w, "Invalid event ID", http.StatusBadRequest)
				return
			}
			eventID = id
			break
		}
	}

	if eventID == -1 {
		http.Error(w, "Event ID not found in path", http.StatusBadRequest)
		return
	}

	type RequestBody struct {
		UserID string `json:"user_id"`
	}

	var req RequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	booking, err := s.BookSeat(eventID, req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

func (s *EventService) handleConfirmBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestBody struct {
		BookingID int `json:"booking_id"`
	}

	var req RequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := s.ConfirmBooking(req.BookingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Booking confirmed")
}

// Templates for web interface
const adminTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>Event Booking Admin</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; }
        input, button { padding: 8px; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .event-form { border: 1px solid #ccc; padding: 15px; margin-bottom: 20px; }
    </style>
</head>
<body>
    <h1>Event Booking Admin Panel</h1>

    <div class="event-form">
        <h2>Create New Event</h2>
        <form id="createEventForm">
            <div class="form-group">
                <label for="name">Event Name:</label>
                <input type="text" id="name" name="name" required>
            </div>
            <div class="form-group">
                <label for="date">Event Date:</label>
                <input type="datetime-local" id="date" name="date" required>
            </div>
            <div class="form-group">
                <label for="total_seats">Total Seats:</label>
                <input type="number" id="total_seats" name="total_seats" min="1" required>
            </div>
            <button type="submit">Create Event</button>
        </form>
    </div>

    <h2>Events List</h2>
    <table id="eventsTable">
        <thead>
            <tr>
                <th>ID</th>
                <th>Name</th>
                <th>Date</th>
                <th>Total Seats</th>
                <th>Available</th>
                <th>Actions</th>
            </tr>
        </thead>
        <tbody id="eventsBody">
        </tbody>
    </table>

    <script>
        document.getElementById('createEventForm').addEventListener('submit', async function(e) {
            e.preventDefault();

            const formData = new FormData(e.target);
            const eventData = {
                name: formData.get('name'),
                date: formData.get('date') + ':00Z',
                total_seats: parseInt(formData.get('total_seats'))
            };

            try {
                const response = await fetch('/events', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(eventData)
                });

                if (response.ok) {
                    loadEvents();
                    e.target.reset();
                } else {
                    alert('Error creating event: ' + await response.text());
                }
            } catch (error) {
                alert('Error creating event: ' + error.message);
            }
        });

        function loadEvents() {
            fetch('/events')
                .then(response => response.json())
                .then(events => {
                    const tbody = document.getElementById('eventsBody');
                    tbody.innerHTML = '';

                    events.forEach(event => {
                        const row = tbody.insertRow();
                        row.innerHTML = '<td>' + event.id + '</td><td>' + event.name + '</td><td>' + new Date(event.date).toLocaleString() + '</td><td>' + event.total_seats + '</td><td>' + event.available + '</td><td><a href="/admin/event/' + event.id + '">View Details</a></td>';
                    });
                })
                .catch(error => console.error('Error loading events:', error));
        }

        // Load events on page load
        loadEvents();
    </script>
</body>
</html>`

const userTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>Event Booking User</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .event-item { border: 1px solid #ccc; padding: 15px; margin-bottom: 15px; }
        .form-group { margin-bottom: 10px; }
        label { display: block; margin-bottom: 5px; }
        input, button { padding: 8px; }
        button { cursor: pointer; }
    </style>
</head>
<body>
    <h1>Event Booking System</h1>

    <h2>Available Events</h2>
    <div id="eventsContainer"></div>

    <script>
        function loadEvents() {
            fetch('/events')
                .then(response => response.json())
                .then(events => {
                    const container = document.getElementById('eventsContainer');
                    container.innerHTML = '';

                    events.forEach(event => {
                        if (event.available > 0) {
                            const eventDiv = document.createElement('div');
                            eventDiv.className = 'event-item';
                            eventDiv.innerHTML = '<h3>' + event.name + '</h3><p><strong>Date:</strong> ' + new Date(event.date).toLocaleString() + '</p><p><strong>Available Seats:</strong> ' + event.available + '/' + event.total_seats + '</p><div class="form-group"><label>Your Name:</label><input type="text" id="user-' + event.id + '" placeholder="Enter your name"></div><button onclick="bookSeat(' + event.id + ', \'' + event.name + '\', document.getElementById(\'user-' + event.id + '\').value)">Book Seat</button>';
                            container.appendChild(eventDiv);
                        }
                    });
                })
                .catch(error => console.error('Error loading events:', error));
        }

        async function bookSeat(eventId, eventName, userName) {
            if (!userName) {
                alert('Please enter your name');
                return;
            }

            try {
                const response = await fetch('/events/' + eventId + '/book', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ user_id: userName })
                });

                if (response.ok) {
                    const booking = await response.json();
                    alert('Seat booked successfully! Your booking ID: ' + booking.id + '. Please confirm payment within 30 minutes or the booking will be cancelled.');

                    // Reload events to show updated availability
                    loadEvents();
                } else {
                    alert('Error booking seat: ' + await response.text());
                }
            } catch (error) {
                alert('Error booking seat: ' + error.message);
            }
        }

        async function confirmBooking(bookingId) {
            if (!confirm('Are you sure you want to confirm booking ID: ' + bookingId + '?')) {
                return;
            }

            try {
                const response = await fetch('/events/-1/confirm', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ booking_id: bookingId })
                });

                if (response.ok) {
                    alert('Booking confirmed successfully!');
                } else {
                    alert('Error confirming booking: ' + await response.text());
                }
            } catch (error) {
                alert('Error confirming booking: ' + error.message);
            }
        }

        // Load events on page load
        loadEvents();
    </script>
</body>
</html>`

func (s *EventService) handleAdminPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("admin").Parse(adminTemplate))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *EventService) handleUserPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("user").Parse(userTemplate))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *EventService) handleGetAllEvents(w http.ResponseWriter, r *http.Request) {
	events, err := s.GetAllEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (s *EventService) setupRoutes(mux *http.ServeMux) {
	// API endpoints and web interface - using traditional mux without pattern matching
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleGetAllEvents(w, r)
		case http.MethodPost:
			s.handleCreateEvent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Handle event-specific routes with a generic handler that checks URL path
	mux.HandleFunc("/events/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the URL to extract the event ID and action
		path := r.URL.Path
		parts := strings.Split(path, "/")

		if len(parts) >= 3 {
			eventID, err := strconv.Atoi(parts[2])
			if err != nil {
				http.Error(w, "Invalid event ID", http.StatusBadRequest)
				return
			}

			if len(parts) == 3 { // /events/{id}
				switch r.Method {
				case http.MethodGet:
					// Create a new request context with the event ID in the URL
					originalURL := r.URL.Path
					r.URL.Path = "/events/" + strconv.Itoa(eventID)
					s.handleGetEvent(w, r)
					r.URL.Path = originalURL // restore original path
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
				return
			}

			if len(parts) >= 4 {
				switch parts[3] {
				case "book":
					if r.Method == http.MethodPost {
						s.handleBookEvent(w, r)
					} else {
						http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					}
				case "confirm":
					if r.Method == http.MethodPost {
						s.handleConfirmBooking(w, r)
					} else {
						http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					}
				default:
					http.Error(w, "Not found", http.StatusNotFound)
				}
			}
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})

	// Web interface
	mux.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleAdminPage(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/admin/event/", func(w http.ResponseWriter, r *http.Request) {
		s.handleGetEvent(w, r)
	})

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleUserPage(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleUserPage(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func main() {
	// Use environment variable for database path or default to "event_booking.db"
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "event_booking.db"
	}

	service, err := NewEventService(dbPath)
	if err != nil {
		log.Fatal("Failed to create event service:", err)
	}
	defer service.db.Close()

	mux := http.NewServeMux()
	service.setupRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	fmt.Printf("Server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
