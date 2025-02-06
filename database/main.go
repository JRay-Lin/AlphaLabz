package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

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
