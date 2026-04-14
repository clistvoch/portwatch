package alert

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDo_SucceedsFirstAttempt(t *testing.T) {
	calls := 0
	err := Do(context.Background(), DefaultRetryConfig(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}
	calls := 0
	sentinel := errors.New("fail")
	err := Do(context.Background(), cfg, func() error {
		calls++
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDo_SucceedsOnSecondAttempt(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		Multiplier:   2.0,
	}
	calls := 0
	err := Do(context.Background(), cfg, func() error {
		calls++
		if calls < 2 {
			return errors.New("not yet")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:  5,
		InitialDelay: 50 * time.Millisecond,
		MaxDelay:     200 * time.Millisecond,
		Multiplier:   2.0,
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := Do(ctx, cfg, func() error {
		return errors.New("always fail")
	})
	if err == nil {
		t.Fatal("expected an error due to cancelled context or exhausted attempts")
	}
}
