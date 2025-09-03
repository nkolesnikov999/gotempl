package main

import (
	"log"
	"nkpro/gotempl/config"
	"nkpro/gotempl/internal/pages"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Init()
	dbConf := config.NewDatabaseConfig()
	log.Println(dbConf)
	app := fiber.New()

	pages.NewHandler(app)

	app.Listen(":3003")
}
