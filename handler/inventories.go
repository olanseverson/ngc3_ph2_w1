package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ngc/entity"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	*sql.DB
}

func (h *Handler) GetInventories(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var inventories []entity.Inventory

	rows, err := h.Query("SELECT kode, nama, stock, description, status, hero_id FROM inventories")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var i entity.Inventory

		err = rows.Scan(&i.Kode, &i.Nama, &i.Stock, &i.Description, &i.Status, &i.HeroID)
		if err != nil {
			log.Fatal(err)
		}
		inventories = append(inventories, i)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventories)
}

func (h *Handler) GetInventoryByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	var i entity.Inventory
	err := h.QueryRow("SELECT kode, nama, stock, description, status, hero_id FROM inventories WHERE kode=?", id).
		Scan(&i.Kode, &i.Nama, &i.Stock, &i.Description, &i.Status, &i.HeroID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatal("Empty Table")
		} else {
			log.Fatal(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(i)
}

func (h *Handler) AddInventory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var inventory entity.Inventory
	// get the data for 'inventory'
	err := json.NewDecoder(r.Body).Decode(&inventory)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	fmt.Println(inventory)
	// ensure the data is valid
	if inventory.Nama == "" || inventory.Stock <= 0 || inventory.Description == "" || inventory.Status == "" || inventory.HeroID <= 0 {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO inventories (nama, stock, description, status, hero_id) VALUES (?, ?, ?, ?, ?)"
	result, err := h.Exec(query, inventory.Nama, inventory.Stock, inventory.Description, inventory.Status, inventory.HeroID)
	if err != nil {
		http.Error(w, "Error inserting inventory", http.StatusInternalServerError)
		return
	}

	// get the last id
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error retrieving last insert ID", http.StatusInternalServerError)
		return
	}

	inventory.Kode = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inventory)
}

func (h *Handler) UpdateInventory(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	kode, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid inventory id", http.StatusBadRequest)
		return
	}

	var i entity.Inventory

	err = json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	fmt.Println(i)
	if i.Stock <= 0 || i.Description == "" || i.Status == "" || i.HeroID <= 0 {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	query := `
	UPDATE inventories
	SET stock = ?,
		description = ?,
		status = ?,
		hero_id = ?
	WHERE kode = ?
	`
	result, err := h.Exec(query, i.Stock, i.Description, i.Status, i.HeroID, kode)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error updating inventory", http.StatusInternalServerError)
		return
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "error getting rows affected", http.StatusInternalServerError)
		return
	}

	if affectedRow == 0 {
		http.Error(w, "inventory not found", http.StatusInternalServerError)
		return
	}

	i.Kode = kode
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(i)
}

func (h *Handler) DeleteInventory(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	kode, err := strconv.Atoi(id)
	if err != nil || kode <= 0 {
		http.Error(w, "invalid inventory id", http.StatusBadRequest)
		return
	}
	result, err := h.Exec("DELETE FROM inventories WHERE kode=?", kode)
	if err != nil {
		http.Error(w, "error while deleting", http.StatusBadRequest)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "error while get rows affected", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "inventory not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
