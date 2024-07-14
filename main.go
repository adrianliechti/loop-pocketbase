package main

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		if err := initSettings(app); err != nil {
			return err
		}

		if err := initAdmin(app); err != nil {
			return err
		}

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func initSettings(app *pocketbase.PocketBase) error {
	settings := app.Settings()

	if val := os.Getenv("APP_URL"); val != "" {
		settings.Meta.AppUrl = val
	}

	if val := os.Getenv("APP_NAME"); val != "" {
		settings.Meta.AppName = val
	}

	return app.Dao().SaveSettings(settings)
}

func initAdmin(app *pocketbase.PocketBase) error {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	if username == "" {
		return nil
	}

	admin, err := app.Dao().FindAdminByEmail(username)

	if err != nil {
		err = nil

		admin = &models.Admin{
			Email: username,
		}

		if password == "" {
			password = uuid.New().String()
			log.Printf("ADMIN_PASSWORD not set. Initial password: %s", password)
		}
	}

	if password != "" {
		admin.SetPassword(password)
	}

	return app.Dao().SaveAdmin(admin)
}
