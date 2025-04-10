package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"note-app/store"
	"os"
)

//go:embed templates
var templateFS embed.FS

// PageData holds the data passed to the HTML template.
type PageData struct {
	Notes []store.Note // Use store.Note
}

var tpl *template.Template

func main() {
	// --- Initialize Template ---
	var err error
	tpl, err = template.ParseFS(templateFS, "templates/index.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// --- Initialize Storage ---
	noteStore, err := store.CreateNotesStore()
	if err != nil {
		log.Fatalf("Error creating notes store: %v", err)
	}

	// --- HTTP Handlers
	http.HandleFunc("/", handleIndex(noteStore))
	http.HandleFunc("/add", handleAddNote(noteStore))
	http.HandleFunc("/delete", handleDeleteNote(noteStore))

	// --- Start Server ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Server starting on port %s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// handleIndex serves the main page, displaying notes.
func handleIndex(noteStore store.NoteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		notes, err := noteStore.GetAllNotes()
		if err != nil {
			log.Printf("Error getting notes from store: %v", err)
			http.Error(w, "Failed to retrieve notes", http.StatusInternalServerError)
			return
		}

		data := PageData{
			Notes: notes,
		}

		err = tpl.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

// handleAddNote adds a new note based on form submission.
func handleAddNote(noteStore store.NoteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		noteText := r.FormValue("noteText")
		if noteText == "" {
			log.Println("Attempted to add empty note")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		_, err = noteStore.AddNote(noteText)
		if err != nil {
			log.Printf("Error adding note to store: %v", err)
			http.Error(w, "Failed to add note", http.StatusInternalServerError)
			return
		}

		// Redirect back to the index page to show the updated list
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// handleDeleteNote removes a note based on its string ID.
func handleDeleteNote(noteStore store.NoteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		noteID := r.FormValue("noteID")
		if noteID == "" {
			log.Println("Attempted to delete note without ID")
			http.Error(w, "Bad request: Missing note ID", http.StatusBadRequest)
			return
		}

		err = noteStore.DeleteNote(noteID)
		if err != nil {
			log.Printf("Error deleting note from store: %v", err)
			// Check specific error if needed (could refine error handling)
			if err.Error() == fmt.Sprintf("note with ID %s not found", noteID) {
				http.Error(w, "Note not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to delete note", http.StatusInternalServerError)
			}
			return
		}

		// Redirect back to the index page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
