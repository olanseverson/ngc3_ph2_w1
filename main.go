package main

import (
	"ngc/config"
	"ngc/handler"
)

func main() {
	router, server := config.SetupServer()
	db := &handler.Handler{DB: config.Connect()}
	router.GET("/inventory", db.GetInventories)
	router.GET("/inventory/:id", db.GetInventoryByID)
	router.POST("/inventory", db.AddInventory)
	router.PUT("/inventory/:id", db.UpdateInventory)
	router.DELETE("/inventory/:id", db.DeleteInventory)

	panic(server.ListenAndServe())
}
