package store

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
)

// InMemoryStore implements NoteStore using an in-memory slice.
type InMemoryStore struct {
	notes []Note
}

// NewInMemoryStore creates a new InMemoryStore with initial data.
func NewInMemoryStore() (*InMemoryStore, error) {
	store := &InMemoryStore{
		notes: []Note{
			{Id: uuid.New().String(), Text: "My first demo note"},
			{Id: uuid.New().String(), Text: "Remember to integrate Azure"},
		},
	}
	return store, nil
}

func (s *InMemoryStore) GetAllNotes() ([]Note, error) {
	notesCopy := make([]Note, len(s.notes))
	copy(notesCopy, s.notes)
	slices.Reverse(notesCopy)
	return notesCopy, nil
}

func (s *InMemoryStore) AddNote(text string) (string, error) {
	newID := uuid.New().String()
	s.notes = append(s.notes, Note{Id: newID, Text: text})
	return newID, nil
}

func (s *InMemoryStore) DeleteNote(id string) error {
	// find index for note id
	index := -1
	for i, note := range s.notes {
		if note.Id == id {
			index = i
			break
		}
	}

	if index != -1 {
		s.notes = slices.Delete(s.notes, index, index+1)
		return nil
	}

	return fmt.Errorf("note with ID %s not found", id)
}
