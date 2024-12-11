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

func ExecCommand(logOutput bool, name string, args ...string) (bool, string) {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Command failed: %s %v\nError: %v\nStderr: %s", name, args, err, stderr.String())
		return false, stderr.String() // Return false with stderr content
	}

	if logOutput {
		log.Printf("Command succeeded: %s %v\nStdout: %s", name, args, out.String())
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

		email := EmailInput()
		password := PasswordInput()

		if ConfirmInput(email, password) {
			saveEnv(email, password)
		}
		screen.Clear()
		screen.MoveTopLeft()

		elimtRunning := CheckElimtRunning()

		if elimtRunning {
			fmt.Println("Would you like to remove them?")

			option := []string{"Yes and Continue Installation", "No and Exit installation"}
			prompt := promptui.Select{
				Label: "Choose an option",
				Items: option,
			}
			index, _, err := prompt.Run()
			if err != nil {
				log.Fatalf("Prompt failed: %v\n", err)
			}

			if index == 0 {
				ClearAllDocker()
				RunDockerCompose()
			} else {
				log.Fatal("CLI terminated by user.") // Terminates the script
			}
		} else {
			RunDockerCompose()
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
