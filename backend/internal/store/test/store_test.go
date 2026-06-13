package store_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/robin38n/reqviz/backend/internal/store"
)

func TestStore_SaveAndGet(t *testing.T) {
	s := store.New()
	id := s.Save(&store.StoredSpec{Title: "My API", Version: "1.0"})

	if id == uuid.Nil {
		t.Fatal("Save returned a nil UUID")
	}

	got, err := s.Get(id)
	if err != nil {
		t.Fatalf("Get(%s) returned error: %v", id, err)
	}
	if got.Title != "My API" {
		t.Errorf("Title = %q, want %q", got.Title, "My API")
	}
	if got.ID != id {
		t.Errorf("ID = %s, want %s", got.ID, id)
	}
	if got.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set by Save")
	}
}

func TestStore_GetUnknown(t *testing.T) {
	s := store.New()
	if _, err := s.Get(uuid.New()); err == nil {
		t.Error("expected error getting an unknown spec id")
	}
}

func TestStore_Approve(t *testing.T) {
	s := store.New()
	id := s.Save(&store.StoredSpec{Title: "My API"})

	hosts := []string{"api.example.com"}
	updated, err := s.Approve(id, hosts)
	if err != nil {
		t.Fatalf("Approve(%s) returned error: %v", id, err)
	}
	if !updated.Approved {
		t.Error("expected Approved = true after Approve")
	}
	if len(updated.AllowedHosts) != 1 || updated.AllowedHosts[0] != "api.example.com" {
		t.Errorf("AllowedHosts = %v, want [api.example.com]", updated.AllowedHosts)
	}

	// The change is persisted in the store.
	got, _ := s.Get(id)
	if !got.Approved {
		t.Error("expected persisted spec to be Approved")
	}
}

func TestStore_ApproveUnknown(t *testing.T) {
	s := store.New()
	if _, err := s.Approve(uuid.New(), nil); err == nil {
		t.Error("expected error approving an unknown spec id")
	}
}
