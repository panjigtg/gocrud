package main

import "crudprojectgo/config"

// @title CRUD Project Go API
// @version 1.0
// @description API documentation for CRUD Project Go
// @host 127.0.0.1:3000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.RunApp()
}
