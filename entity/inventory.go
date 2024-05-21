package entity

type Inventory struct {
	Kode        int    `json:"kode"`
	Nama        string `json:"nama"`
	Stock       int    `json:"stock"`
	Description string `json:"description"`
	Status      string `json:"status"`
	HeroID      int    `json:"hero_id"`
}
