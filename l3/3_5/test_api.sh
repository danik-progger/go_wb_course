#!/bin/bash

echo "Testing Event Booking Service API..."

# Test 1: Create an event
echo "1. Creating an event..."
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Go Workshop",
    "date": "2025-12-15T18:00:00Z",
    "total_seats": 5
  }'

echo -e "\n2. Getting all events..."
curl http://localhost:8080/events

# Test 3: Book a seat
echo -e "\n3. Booking a seat..."
curl -X POST http://localhost:8080/events/1/book \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test_user"
  }'

echo -e "\n4. Getting event details..."
curl http://localhost:8080/events/1

# Test 5: Confirm the booking
echo -e "\n5. Confirming the booking..."
curl -X POST http://localhost:8080/events/-1/confirm \
  -H "Content-Type: application/json" \
  -d '{
    "booking_id": 1
  }'

echo -e "\n6. Getting event details after confirmation..."
curl http://localhost:8080/events/1

echo -e "\nAPI testing complete!"
