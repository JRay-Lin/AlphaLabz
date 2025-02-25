package main

import (
	"alphalabz/pkg/casbin"
	"alphalabz/pkg/pocketbase"
	"alphalabz/pkg/routes/labbook"
	"alphalabz/pkg/routes/login"
	"alphalabz/pkg/routes/role"
	"alphalabz/pkg/routes/user"
	"alphalabz/pkg/settings"
	"alphalabz/pkg/smtp"
	"alphalabz/pkg/tools"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

var pbClient *pocketbase.PocketBaseClient
var casbinEnforcer *casbin.CasbinEnforcer
var SMTPClient *smtp.SMTPClient

// JWTAuthMiddleware will check the JWT token and validate it. If valid, it will pass the request to the next handler. Otherwise, it will return a 401 Unauthorized response.
func JWTExpirationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jwtSkipPaths = map[string]bool{
			"/health":        true,
			"/login/account": true,
			"/login/oauth":   true,
			"/login/sso":     true,
			"/users/signup":  true,
		}

		// Check if the path is in the skip list. If it is, then skip JWT validation and pass the request to the next handler.
		if _, ok := jwtSkipPaths[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		rawToken, err := tools.TokenExtractor(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "token missing or invalid", http.StatusUnauthorized)
			return
		}

		// Verify JWT expiration. If the token has expired or is invalid, return a 401 Unauthorized response.
		valid, err := tools.VerifyJWTExpiration(rawToken)
		if err != nil || !valid {
			http.Error(w, "token expired or invalid", http.StatusUnauthorized)
			return
		}

		// Token is valid. Pass the request to the next handler.
		next.ServeHTTP(w, r)
	})
}

func setupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Use(JWTExpirationMiddleware)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Server is healthy")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Corrected JSON encoding
		response := map[string]string{"message": "server is healthy"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	// Login to system
	r.Route("/login", func(r chi.Router) {
		r.Post("/account", func(w http.ResponseWriter, r *http.Request) {
			login.HandleAccountLogin(w, r, pbClient)
		})

		// r.Post("/oauth", func(w http.ResponseWriter, r *http.Request) {
		// 	// routes.HandleOAuth(w, r, pbClient)
		// })

		// r.Post("/sso", func(w http.ResponseWriter, r *http.Request) {
		// 	// routes.HandleSSO(w, r, pbClient)
		// })
	})

	// Users route
	r.Route("/users", func(r chi.Router) {
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			user.HandleUserList(w, r, pbClient, casbinEnforcer)
		})

		// !!! Deprecated !!!
		// r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		// 	user.HandleRegister(w, r, pbClient)
		// })

		r.Post("/invite", func(w http.ResponseWriter, r *http.Request) {
			user.HandleInviteNewUser(w, r, pbClient, casbinEnforcer, SMTPClient)
		})

		r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
			user.HandleSignUp(w, r, pbClient, casbinEnforcer)
		})

		r.Delete("/remove", func(w http.ResponseWriter, r *http.Request) {
			user.HandleUserRemove(w, r, pbClient, casbinEnforcer)
		})

		r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleUserUpdate(w, r, pbClient)
		})
	})

	// Lab_book route
	r.Route("/labbook", func(r chi.Router) {
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleLabBookList(w, r, pbClient)
		})

		r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
			labbook.HandleLabBookUpload(w, r, pbClient, casbinEnforcer)
		})

		r.Delete("/remove", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleLabBookRemove(w, r, pbClient)
		})

		r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleLabBookUpdate(w, r, pbClient)
		})

		r.Post("/verify", func(w http.ResponseWriter, r *http.Request) {
			labbook.HandleLabBookVerify(w, r, pbClient, casbinEnforcer)
		})
	})

	// Schedule routes
	r.Route("/schedule", func(r chi.Router) {
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleList(w, r, pbClient)
		})
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleCreate(w, r, pbClient)
		})
		r.Delete("/remove", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleRemove(w, r, pbClient)
		})
		r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleUpdate(w, r, pbClient)
		})
	})

	// Link routes
	r.Route("/link", func(r chi.Router) {
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleList(w, r, pbClient)
		})
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleCreate(w, r, pbClient)
		})
		r.Delete("/remove", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleRemove(w, r, pbClient)
		})
		r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleUpdate(w, r, pbClient)
		})
	})

	r.Route("/resources", func(r chi.Router) {
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleList(w, r, pbClient)
		})
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleCreate(w, r, pbClient)
		})
		r.Delete("/remove", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleRemove(w, r, pbClient)
		})
		r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleUpdate(w, r, pbClient)
		})
	})

	r.Route("/roles", func(r chi.Router) {
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			role.HandleRoleList(w, r, pbClient, casbinEnforcer)
		})
	})

	return r
}

func main() {
	settings, err := settings.LoadSettings("settings.yml")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Settings loaded successfully")
	}

	pbHost := os.Getenv("POCKETBASE_URL")
	if pbHost == "" {
		pbHost = "http://127.0.0.1:8090"
	}

	// Get admin password from env
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminEmail == "" || adminPassword == "" {
		log.Fatal("Missing required environment variables: ADMIN_EMAIL and ADMIN_PASSWORD")
	}

	pbClient, err = pocketbase.NewPocketBase(pbHost, adminEmail, adminPassword, 10, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to initialize PocketBase client: %v", err)
	}
	pbClient.StartSuperTokenAutoRenew(adminEmail, adminPassword)

	policies, err := casbin.FetchPermissions(pbClient)
	if err != nil {
		log.Fatalf("Failed to fetch policies: %v", err)
	}

	casbinEnforcer, err = casbin.InitializeCasbin(policies)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin: %v", err)
	}

	casbinEnforcer.StartPolicyAutoReload(pbClient, 60*time.Minute)

	SMTPClient = smtp.NewSMTPClient(
		settings.Mailer.Port,
		settings.Mailer.Host,
		settings.Mailer.Username,
		settings.Mailer.Password,
		settings.Mailer.FromAddress,
		settings.Mailer.FromName)

	if err = tools.CreateUploadsDir(); err != nil {
		log.Fatal("Failed to create uploads directory")
	}

	tools.StartAutoCleanUploads(24 * time.Hour)

	// Setup and start server
	r := setupRouter()
	port := settings.Server.Port
	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
