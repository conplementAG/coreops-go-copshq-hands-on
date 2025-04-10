package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/google/uuid"
)

// AzureTableStore implements NoteStore using Azure Table Storage.
type AzureTableStore struct {
	serviceClient *aztables.ServiceClient
	client        *aztables.Client
	tableName     string
}

const partitionKey = "notes"

// NewAzureTableStore creates a new AzureTableStore using a connection string.
func NewAzureTableStoreWithConnectionString(connectionString string) (*AzureTableStore, error) {
	if connectionString == "" {
		return nil, fmt.Errorf("azure connection string cannot be empty")
	}

	serviceClient, err := aztables.NewServiceClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Printf("Error creating Azure Table service client from connection string: %v\n", err)
		return nil, fmt.Errorf("failed to create Azure Table service client: %w", err)
	}

	store := newAzureTableStoreInternal(serviceClient)
	log.Printf("Successfully initialized Azure Table Store client for table '%s' using connection string", store.tableName)

	return store, nil
}

// NewAzureTableStoreWithDefaultCredentials creates a new AzureTableStore using DefaultAzureCredential.
func NewAzureTableStoreWithDefaultCredentials(storageAccount string) (*AzureTableStore, error) {
	if storageAccount == "" {
		return nil, fmt.Errorf("azure storage account name cannot be empty")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Printf("Error creating default Azure credential: %v\n", err)
		return nil, fmt.Errorf("failed to create default Azure credential: %w", err)
	}

	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net", storageAccount)
	serviceClient, err := aztables.NewServiceClient(serviceURL, cred, nil)
	if err != nil {
		log.Printf("Error creating Azure Table service client with default credential: %v\n", err)
		return nil, fmt.Errorf("failed to create Azure Table service client: %w", err)
	}

	store := newAzureTableStoreInternal(serviceClient)
	log.Printf("Successfully initialized Azure Table Store client for account '%s', table '%s' using DefaultAzureCredential", storageAccount, store.tableName)

	return store, nil
}

// newAzureTableStoreInternal handles the common initialization logic after serviceClient is created.
func newAzureTableStoreInternal(serviceClient *aztables.ServiceClient) *AzureTableStore {
	tableName := "Notes" // Centralized table name definition
	client := serviceClient.NewClient(tableName)

	return &AzureTableStore{
		serviceClient: serviceClient, // Corrected field name
		client:        client,
		tableName:     tableName,
	}
}

// noteEntity represents the structure of a note entity in Azure Table Storage.
type noteEntity struct {
	aztables.Entity
	Text string `json:"Text"`
}

// GetAllNotes fetches all notes from Azure Table Storage.
func (s *AzureTableStore) GetAllNotes() ([]Note, error) {
	log.Println("AzureStore: GetAllNotes called")
	notes := []Note{}
	pager := s.client.NewListEntitiesPager(nil) // No filter, list all entities

	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		if err != nil {
			log.Printf("Error fetching page: %v\n", err)
			return nil, fmt.Errorf("failed to list notes: %w", err)
		}

		for _, entityBytes := range resp.Entities {
			var entity noteEntity
			if err := json.Unmarshal(entityBytes, &entity); err != nil {
				log.Printf("Error unmarshalling entity: %v\n", err)
				continue // Skip problematic entities
			}

			// RowKey is now the string ID (UUID)
			notes = append(notes, Note{
				Id:   entity.RowKey,
				Text: entity.Text,
			})
		}
	}

	log.Printf("AzureStore: Found %d notes", len(notes))
	return notes, nil
}

// AddNote adds a note to Azure Table Storage using a UUID as the ID.
func (s *AzureTableStore) AddNote(text string) (string, error) {
	log.Printf("AzureStore: AddNote called for text: %s", text)

	// Generate a new UUID for the RowKey
	newID := uuid.New().String()

	newNoteEntity := noteEntity{
		Entity: aztables.Entity{
			PartitionKey: partitionKey,
			RowKey:       newID,
		},
		Text: text,
	}

	entityBytes, err := json.Marshal(newNoteEntity)
	if err != nil {
		log.Printf("Error marshalling note entity: %v\n", err)
		return "", fmt.Errorf("failed to marshal note: %w", err)
	}

	_, err = s.client.AddEntity(context.Background(), entityBytes, nil)
	if err != nil {
		log.Printf("Error adding note entity: %v\n", err)
		return "", fmt.Errorf("failed to add note to table: %w", err)
	}

	log.Printf("AzureStore: Successfully added note with ID: %s", newID)
	return newID, nil
}

// DeleteNote deletes a note from Azure Table Storage by its string ID (UUID).
func (s *AzureTableStore) DeleteNote(id string) error {
	log.Printf("AzureStore: DeleteNote called for ID: %s", id)

	// ID is already the string RowKey
	rowKey := id

	// ETag can be used for optimistic concurrency. Using nil for unconditional delete.
	_, err := s.client.DeleteEntity(context.Background(), partitionKey, rowKey, nil)

	if err != nil {
		// TODO: Check for specific errors like 'Not Found' if needed
		log.Printf("Error deleting note entity with ID (RowKey): %s: %v\n", rowKey, err)
		return fmt.Errorf("failed to delete note %s: %w", id, err)
	}

	log.Printf("AzureStore: Successfully deleted note with ID: %s", id)
	return nil
}
