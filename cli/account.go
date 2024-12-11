package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/inancgumus/screen"
	"github.com/manifoldco/promptui"
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func EmailInput() string {
	prompt := promptui.Prompt{
		Label:   "Enter Admin Email (leave blank to generate one)",
		Default: "",
		Validate: func(input string) error {
			if input == "" {
				return nil
			}
			emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
			if match := regexp.MustCompile(emailRegex).MatchString(input); !match {
				return fmt.Errorf("invalid email address")
			}
			return nil
		},
	}

	email, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}

	email = strings.TrimSpace(email)
	if email == "" {
		email = GenerateRandomString(6) + "@elimt.com"
		// fmt.Printf("Generated email: %s\n", email)
	}

	return email
}

func PasswordInput() string {
	prompt := promptui.Prompt{
		Label:   "Enter Admin Password (leave blank to generate one)",
		Mask:    '*',
		Default: "",
		Validate: func(input string) error {
			if len(input) == 0 {
				return nil
			}
			if len(input) < 8 {
				return fmt.Errorf("password must be at least 8 characters")
			}
			return nil
		},
	}

	password, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}

	password = strings.TrimSpace(password)
	if password == "" {
		password = GenerateRandomString(16)
		// fmt.Printf("Generated password: %s\n", password)
	}

	return password
}

func ConfirmInput(email, password string) bool {
	screen.Clear()
	screen.MoveTopLeft()

	options := []string{"Confirm and start installation", "Re-enter email and password", "Exit setup"}

	for {
		fmt.Println(strings.Repeat("=", 40))
		fmt.Printf("Email   : %s\n", email)
		fmt.Printf("Password: %s\n", password)
		fmt.Println(strings.Repeat("=", 40))

		prompt := promptui.Select{
			Label: "Choose an option",
			Items: options,
		}

		index, _, err := prompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed: %v\n", err)
		}

		switch index {
		case 0: // Confirm and start installation
			screen.Clear()
			return true
		case 1: // Re-enter email and password
			screen.Clear()
			screen.MoveTopLeft()
			return false
		case 2: // Exit setup
			fmt.Println("Exiting setup. No changes made.")
			log.Fatal("Setup terminated by user.") // Terminates the script
		}
	}
}
