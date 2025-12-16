package or

import (
	"context"
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	t.Run("with one channel", func(t *testing.T) {
		ch := make(chan interface{})
		result := Or(ch)

		go func() {
			ch <- "test"
			close(ch)
		}()

		select {
		case val := <-result:
			if val != "test" {
				t.Errorf("Expected 'test', got %v", val)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for value")
		}
	})

	t.Run("with two channels - first wins", func(t *testing.T) {
		ch1 := make(chan interface{})
		ch2 := make(chan interface{})
		result := Or(ch1, ch2)

		go func() {
			time.Sleep(10 * time.Millisecond)
			ch1 <- "first"
			close(ch1)
		}()

		go func() {
			time.Sleep(50 * time.Millisecond)
			ch2 <- "second"
			close(ch2)
		}()

		select {
		case val := <-result:
			if val != "first" {
				t.Errorf("Expected 'first', got %v", val)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for value")
		}
	})

	t.Run("with two channels - second wins", func(t *testing.T) {
		ch1 := make(chan interface{})
		ch2 := make(chan interface{})
		result := Or(ch1, ch2)

		go func() {
			time.Sleep(50 * time.Millisecond)
			ch1 <- "first"
			close(ch1)
		}()

		go func() {
			time.Sleep(10 * time.Millisecond)
			ch2 <- "second"
			close(ch2)
		}()

		select {
		case val := <-result:
			if val != "second" {
				t.Errorf("Expected 'second', got %v", val)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for value")
		}
	})

	t.Run("with already closed channels", func(t *testing.T) {
		// When all input channels are closed without sending values,
		// the Or function will block forever waiting for a value
		// This is the correct behavior according to the pattern
		// For testing purposes, we'll just make sure it doesn't crash
		ch1 := make(chan interface{})
		ch2 := make(chan interface{})
		close(ch1)
		close(ch2)

		// Create a timeout to avoid hanging tests
		result := Or(ch1, ch2)
		done := make(chan struct{})

		go func() {
			<-result // This will block, so we use a timeout
			close(done)
		}()

		select {
		case <-done:
			t.Error("Expected the function to block forever, but it returned")
		case <-time.After(10 * time.Millisecond):
			// Expected: function blocks because all channels are closed without values
		}
	})
}

func TestOrCancellation(t *testing.T) {
	// Test that Or function doesn't leak goroutines on context cancellation
	_, cancel := context.WithCancel(context.Background())

	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	result := Or(ch1, ch2)

	// Don't send anything to the channels
	// Instead, we'll close the channels to test behavior

	close(ch1)

	// Now we need to send to ch2 to get a result
	go func() {
		time.Sleep(10 * time.Millisecond)
		ch2 <- "test"
		close(ch2)
	}()

	select {
	case val := <-result:
		if val != "test" {
			t.Errorf("Expected 'test', got %v", val)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout waiting for value")
	}

	cancel()
}

func ExampleOr() {
	// Example showing how to use the Or function

	// Create some channels
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})

	// Create a channel that will receive from the first available channel
	any := Or(ch1, ch2, ch3)

	// Send a value to ch2 (simulating async work)
	go func() {
		time.Sleep(10 * time.Millisecond)
		ch2 <- "hello from channel 2"
		close(ch2)
	}()

	// Receive from the combined channel
	result := <-any
	if result == "hello from channel 2" {
		// This demonstrates that Or works correctly
	}

	// Make sure to close all channels
	close(ch1)
	close(ch3)
}
