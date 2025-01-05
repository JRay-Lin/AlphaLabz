package installation

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/manifoldco/promptui"
)

type InstallCmd struct{}

const repoUrl string = "https://github.com/JRay-Lin/AlphaLab.git"

func (g *InstallCmd) Run() error {
	fmt.Println("Welcome to the AlphaLab installation process!")

	if err := checkRequirements(); err != nil {
		return fmt.Errorf("failed to check requirements: %v", err)
	}

	email, err := promptEmail()
	if err != nil {
		return fmt.Errorf("email input error: %v", err)
	}

	password, err := promptPassword()
	if err != nil {
		return fmt.Errorf("password input error: %v", err)
	}

	if _, err := promptVerifyPassword(password); err != nil {
		return fmt.Errorf("password verification error: %v", err)
	}

	for {
		fmt.Println("\n========== Your Settings ==========")
		fmt.Printf("Admin Email: %s\n", email)
		fmt.Printf("Admin Password: %s\n", password)
		fmt.Println("===================================")

		confirm := promptui.Select{
			Label: "Please confirm your settings",
			Items: []string{"Confirm and continue", "Modify settings"},
		}

		idx, _, err := confirm.Run()
		if err != nil {
			return fmt.Errorf("confirmation selection error: %v", err)
		}

		if idx == 0 {
			// User confirmed, break the loop
			break
		}

		// User wants to modify settings
		choice := promptui.Select{
			Label: "What would you like to change",
			Items: []string{"Email", "Password", "Back to confirmation"},
		}

		idx, _, err = choice.Run()
		if err != nil {
			return fmt.Errorf("modification selection error: %v", err)
		}

		switch idx {
		case 0: // Email
			email, err = promptEmail()
			if err != nil {
				return fmt.Errorf("email modification error: %v", err)
			}
		case 1: // Password
			password, err = promptPassword()
			if err != nil {
				return fmt.Errorf("password modification error: %v", err)
			}
			if _, err := promptVerifyPassword(password); err != nil {
				return fmt.Errorf("password verification error: %v", err)
			}
		case 2: // Back to confirmation
			continue
		}
	}

	Clear()

	fmt.Println("\nStarting AlphaLab installation...")
	if err := installAlphalab(email, password); err != nil {
		return fmt.Errorf("failed to install AlphaLab: %v", err)
	}

	fmt.Println("Verifying AlphaLab installation...")
	if err := verifyAlphalabStatus(); err != nil {
		return fmt.Errorf("failed to verify AlphaLab installation: %v", err)
	}

	fmt.Println("AlphaLab installation completed successfully")
	return nil
}

func promptEmail() (string, error) {
	validate := func(input string) error {
		if !isValidEmail(input) {
			return fmt.Errorf("invalid email format")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Admin Email",
		Validate: validate,
	}

	return prompt.Run()
}

func promptPassword() (string, error) {
	validate := func(input string) error {
		if !isValidPassword(input) {
			return fmt.Errorf("password must be at least 8 characters long and include uppercase, lowercase, and a number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Admin Password",
		Mask:     '*',
		Validate: validate,
	}

	return prompt.Run()
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func isValidPassword(password string) bool {
	length := len(password) >= 8
	upper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	lower := regexp.MustCompile(`[a-z]`).MatchString(password)
	number := regexp.MustCompile(`[0-9]`).MatchString(password)
	return length && upper && lower && number
}

func promptVerifyPassword(originalPassword string) (string, error) {
	validate := func(input string) error {
		if input != originalPassword {
			return fmt.Errorf("passwords do not match")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Verify Password",
		Mask:     '*',
		Validate: validate,
	}

	return prompt.Run()
}

func checkRequirements() error {
	_, err := exec.Command("docker", "--version").Output()
	if err != nil {
		return fmt.Errorf("docker is not installed, please install docker first")
	}

	_, err = exec.Command("docker", "compose", "version").Output()
	if err != nil {
		return fmt.Errorf("docker compose is not installed, please install docker compose first")
	}

	_, err = exec.Command("git", "--version").Output()
	if err != nil {
		return fmt.Errorf("git is not installed, please install git first")
	}

	return nil
}

func installAlphalab(ADMIN_EMAIL string, ADMIN_PASSWORD string) error {
	fmt.Println("Cloning AlphaLab from GitHub...")
	// Check Alphalab repo exits
	_, err := exec.Command("git", "clone", repoUrl).Output()
	if err != nil {
		return fmt.Errorf("failed to clone AlphaLab: %v", err)
	}
	fmt.Println("AlphaLab cloned successfully")

	// Enter folder
	err = os.Chdir("AlphaLab")
	if err != nil {
		return fmt.Errorf("failed to enter AlphaLab folder: %v", err)
	}

	// Run docker compose
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("ADMIN_EMAIL=%s", ADMIN_EMAIL),
		fmt.Sprintf("ADMIN_PASSWORD=%s", ADMIN_PASSWORD),
	)

	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run docker compose: %v", err)
	} else {
		fmt.Println("AlphaLab is installing successfully")
	}

	return nil
}

func verifyAlphalabStatus() error {
	// Placeholder for verification logic
	return nil
}

func Clear() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Print("\033[H\033[2J")
	}
}
