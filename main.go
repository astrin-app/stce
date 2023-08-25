package main

import (
	"astrin/main/db"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type PDF struct {
	ID      string
	SUBJECT string
	META    string
	URL     string
	DESC    string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := os.Getenv("DBURL")
	DB, err := db.CreateDB(dbUrl)
	if err != nil {
		fmt.Print(err)
	}
	engine := mustache.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		pdfs := []PDF{}
		rows, err := DB.Query("SELECT * FROM pdfs;")
		if err != nil {
			log.Fatalf("Query: %v", err)
		}
		for rows.Next() {
			var pdf PDF
			err = rows.Scan(&pdf.ID, &pdf.META, &pdf.SUBJECT, &pdf.URL, &pdf.DESC)
			if err != nil {
				log.Fatalf("Scan: %v", err)
			}
			pdfs = append(pdfs, pdf)
		}
		fmt.Print(pdfs)
		return c.Render("index", fiber.Map{
			"PDFS": pdfs,
		})
	})
	app.Get("/create", func(c *fiber.Ctx) error {
		return c.Render("create", &fiber.Map{})
	})
	app.Post("/create", func(c *fiber.Ctx) error {
		newPDF := PDF{
			ID:      uuid.NewString(),
			SUBJECT: c.FormValue("subject"),
			META:    c.FormValue("meta"),
			URL:     c.FormValue("url"),
			DESC:    c.FormValue("desc"),
		}
		query := fmt.Sprintf("INSERT INTO PDFS (ID, SUBJECT, META,URL,DESC) VALUES ('%s', '%s', '%s', '%s','%s');", newPDF.ID, newPDF.META, newPDF.SUBJECT, newPDF.URL, newPDF.DESC)
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
		c.Set("HX-Redirect", "/")
		return c.Send([]byte("Hello World"))
	})
	app.Listen(":6969")
}
