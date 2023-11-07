package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	baseTemplate := "templates/base.html"
	tmpl = fmt.Sprintf("templates/%s", tmpl)
	t, err := template.ParseFiles(baseTemplate, tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	// Render index template
	renderTemplate(w, "index.html", nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	// Query the database for all fruits
	rows, err := DB.Query("SELECT * FROM fruits")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var fruits []Fruit
	for rows.Next() {
		var fruit Fruit
		if err := rows.Scan(&fruit.ID, &fruit.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fruits = append(fruits, fruit)
	}

	// Render list template with the list of fruits
	// Render the list.html template with the list of fruits
	renderTemplate(w, "list.html", fruits)
}

func fruitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the fruit ID from the URL
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Fruit ID not provided", http.StatusBadRequest)
		return
	}

	// Query the database to retrieve the fruit details by ID
	var fruit Fruit
	err := DB.QueryRow("SELECT id, name FROM fruits WHERE id = $1", id).Scan(&fruit.ID, &fruit.Name)
	if err != nil {
		http.Error(w, "Fruit not found", http.StatusNotFound)
		return
	}

	// Render the single.html template with the fruit details
	renderTemplate(w, "fruit.html", fruit)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Render the add.html template for GET requests
		renderTemplate(w, "add.html", nil)

	case http.MethodPost:

		name := r.FormValue("name")

		// Insert new fruit into the database
		_, err := DB.Query("INSERT INTO fruits (name) VALUES ($1)", name)
		if err != nil {
			http.Error(w, "Failed to insert into the database", http.StatusInternalServerError)
			log.Printf("Failed to insert into the database: %v", err)
			return
		}

		http.Redirect(w, r, "/list", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func editHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Render the edit.html template for GET requests
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Fruit ID not provided", http.StatusBadRequest)
			return
		}

		// Query the database to retrieve the fruit details by ID
		var fruit Fruit
		err := DB.QueryRow("SELECT id, name FROM fruits WHERE id = $1", id).Scan(&fruit.ID, &fruit.Name)
		if err != nil {
			http.Error(w, "Fruit not found", http.StatusNotFound)
			return
		}

		renderTemplate(w, "edit.html", fruit)
	case http.MethodPost:
		// Handle form submission and update the fruit for POST requests
		id := r.FormValue("id")
		name := r.FormValue("name")

		_, err := DB.Exec("UPDATE fruits SET name = $1 WHERE id = $2", name, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/list", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Render the delete.html template for GET requests
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Fruit ID not provided", http.StatusBadRequest)
			return
		}

		// Query the database to retrieve the fruit details by ID
		var fruit Fruit
		err := DB.QueryRow("SELECT id, name FROM fruits WHERE id = $1", id).Scan(&fruit.ID, &fruit.Name)
		if err != nil {
			http.Error(w, "Fruit not found", http.StatusNotFound)
			return
		}

		renderTemplate(w, "delete.html", fruit)
	case http.MethodPost:
		// Handle form submission and delete the fruit for POST requests
		id := r.FormValue("id")

		_, err := DB.Exec("DELETE FROM fruits WHERE id = $1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/list", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
