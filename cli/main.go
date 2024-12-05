package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/inancgumus/screen"
	"github.com/manifoldco/promptui"
)

// checkCommand checks if a command exists in the system.
func ExecCommand(name string, args ...string) (bool, string) {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return false, stderr.String() // Return false with stderr content
	}
	return true, out.String() // Return true with stdout content
}

func saveEnv(email, password string) {
	file, err := os.Create(".env")
	if err != nil {
		log.Fatalf("Failed to create .env file: %v\n", err)
	}
	defer file.Close()

	content := fmt.Sprintf("ADMIN_EMAIL=%s\nADMIN_PASSWORD=%s\n", email, password)
	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalf("Failed to write to .env file: %v\n", err)
	}
}

func main() {
	screen.Clear()
	screen.MoveTopLeft()
	fmt.Println("Welcome to the ELIMT CLI!")
	fmt.Println(strings.Repeat("-", 40))

	dockerInstalled := CheckDockerInstalled()
	if !dockerInstalled {
		fmt.Println("Docker is not installed. Please install Docker and try again.")
	}

	funcOptions := []string{"Setup ELIMT", "Update ELIMT", "Exit"}
	prompt := promptui.Select{
		Label: "Choose an option",
		Items: funcOptions,
	}

	index, _, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}

	switch index {
	case 0:
		screen.Clear()
		screen.MoveTopLeft()
		for {
			email := EmailInput()
			password := PasswordInput()

			if ConfirmInput(email, password) {
				saveEnv(email, password)
				fmt.Println("Setup completed successfully!")
				break
			}
		}
	case 1:
		screen.Clear()
		screen.MoveTopLeft()
		fmt.Println("Update ELIMT")

	case 2: // Exit setup
		fmt.Println("Exiting setup. No changes made.")
		log.Fatal("CLI terminated by user.") // Terminates the script
	}

}
