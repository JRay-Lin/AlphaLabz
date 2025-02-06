package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func verifyToken(header string, app *pocketbase.PocketBase) error {
	if header == "" {
		return fmt.Errorf("Unauthorized")
	}

	if strings.HasPrefix(header, "Bearer ") {
		requesterToken := strings.TrimSpace(header[7:]) // Trim spaces for safety
		if requesterToken == "" {
			return fmt.Errorf("Unauthorized")
		}

		requester, err := app.FindAuthRecordByToken(requesterToken, core.TokenTypeAuth)
		if err != nil {
			return fmt.Errorf("Unauthorized")
		}

		if requester.Collection().Name != "users" {
			return fmt.Errorf("Unauthorized")
		}
	}

	return nil
}

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		// Return user's permission
		se.Router.GET("/api/permissions/{token}", func(e *core.RequestEvent) error {
			// Hanlde requester permission
			if err := verifyToken(e.Request.Header.Get("Authorization"), app); err != nil {
				return e.String(http.StatusUnauthorized, "Unauthorized")
			}

			// Grant token from parameter
			token := e.Request.PathValue("token")
			user, err := app.FindAuthRecordByToken(token, core.TokenTypeAuth)
			if err != nil {
				return e.String(http.StatusUnauthorized, "Unauthorized")
			} else if user.Collection().Name != "users" {
				return e.String(http.StatusForbidden, "Unauthorized")
			}

			userRecord, err := app.FindRecordById("users", user.Id)
			if err != nil {
				return e.String(http.StatusInternalServerError, "Internal Server error")
			}

			fmt.Printf("User permission record: %v\n", userRecord.GetString("permissions"))

			userPermissions, err := app.FindRecordById("permissions", userRecord.GetString("permissions"))
			if err != nil {
				return e.String(http.StatusInternalServerError, "Internal Server error")
			}

			return e.String(200, userPermissions.GetString("permissions"))
		})

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
