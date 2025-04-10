package store

import "os"

// Note represents a single note item.
type Note struct {
	Id   string
	Text string
}

// NoteStore defines the interface for note persistence operations.
type NoteStore interface {
	GetAllNotes() ([]Note, error)
	AddNote(text string) (string, error)
	DeleteNote(id string) error
}

func CreateNotesStore() (NoteStore, error) {
	// CONNECTION_STRING with account name and access key
	connectionString := os.Getenv("CONNECTION_STRING")
	// storage account for usage with default azure credentials
	storageAccount := os.Getenv("STORAGE_ACCOUNT")

	if connectionString != "" {
		return NewAzureTableStoreWithConnectionString(connectionString)
	} else if storageAccount != "" {
		return NewAzureTableStoreWithDefaultCredentials(storageAccount)
	} else {
		return NewInMemoryStore()
	}
}
