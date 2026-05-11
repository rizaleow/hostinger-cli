package poll

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWaitTerminalSuccess(t *testing.T) {
	calls := 0
	fetch := func(_ context.Context, _ int, _ int) (Action, error) {
		calls++
		if calls < 3 {
			return Action{ID: 1, State: "created"}, nil
		}
		return Action{ID: 1, State: StateSuccess}, nil
	}
	ctx := context.Background()
	a, err := Wait(ctx, 100, 1, fetch, Options{Interval: 1 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	if a.State != StateSuccess {
		t.Errorf("state = %q", a.State)
	}
	if calls != 3 {
		t.Errorf("calls = %d", calls)
	}
}

func TestWaitTerminalError(t *testing.T) {
	fetch := func(context.Context, int, int) (Action, error) {
		return Action{ID: 2, Name: "purchase", State: StateError}, nil
	}
	_, err := Wait(context.Background(), 1, 2, fetch, Options{Interval: time.Millisecond})
	if err == nil {
		t.Fatal("expected error for terminal failure state")
	}
}

func TestWaitFetchErrorPropagates(t *testing.T) {
	want := errors.New("network")
	fetch := func(context.Context, int, int) (Action, error) {
		return Action{}, want
	}
	_, err := Wait(context.Background(), 1, 2, fetch, Options{Interval: time.Millisecond})
	if !errors.Is(err, want) {
		t.Errorf("err = %v", err)
	}
}

func TestWaitTimeout(t *testing.T) {
	fetch := func(context.Context, int, int) (Action, error) {
		return Action{State: "created"}, nil
	}
	_, err := Wait(context.Background(), 1, 2, fetch,
		Options{Interval: 10 * time.Millisecond, Timeout: 30 * time.Millisecond})
	if err == nil {
		t.Fatal("expected deadline error")
	}
}
