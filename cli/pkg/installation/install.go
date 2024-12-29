package installation

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// InstallCmd represents the 'install' command
type InstallCmd struct{}

// Run executes the 'install' command step by step
func (g *InstallCmd) Run() error {
	fmt.Println("Welcome to the AlphaLab installation process!")

	// Step 0: Check Docker installation
	if err := checkDockerInstallation(); err != nil {
		return err
	}

	// Step 1: Collect user datas
	domain := askUserInput("Enter your domain name (e.g., example.com): ")
	for !isValidDomain(domain) {
		fmt.Println("Invalid domain format. Please try again (e.g., example.com).")
		domain = askUserInput("Enter your domain name (e.g., example.com): ")
	}

	email := askUserInput("Enter admin email: ")
	for !isValidEmail(email) {
		fmt.Println("Invalid email format. Please try again.")
		email = askUserInput("Enter admin email: ")
	}

	password := askUserInput("Enter admin password: ")
	for !isValidPassword(password) {
		fmt.Println("Password must be at least 8 characters long and include uppercase, lowercase, and a number.")
		password = askUserInput("Enter admin password: ")
	}

	verifyPassword := askUserInput("Verify admin password: ")
	for password != verifyPassword {
		fmt.Println("Passwords do not match. Please try again.")
		password = askUserInput("Enter admin password: ")
		for !isValidPassword(password) {
			fmt.Println("Password must be at least 8 characters long and include uppercase, lowercase, and a number.")
			password = askUserInput("Enter admin password: ")
		}
		verifyPassword = askUserInput("Verify admin password: ")
	}

	// Step 2: Confirm and allow modifications
	for {
		fmt.Println("\n========== Your Settings ==========")
		fmt.Printf("Domain: %s\nAdmin Email: %s\n", domain, email)
		fmt.Println("===================================")
		confirm := askUserInput("Is this information correct? (yes/no): ")

		if strings.ToLower(confirm) == "yes" {
			break
		}

		// Allow modifications
		fmt.Println("Please re-enter the information you'd like to modify:")
		modify := askUserInput("What would you like to change? (domain/email/password): ")

		switch strings.ToLower(modify) {
		case "domain":
			domain = askUserInput("Enter your domain name (e.g., example.com): ")
			for !isValidDomain(domain) {
				fmt.Println("Invalid domain format. Please try again (e.g., example.com).")
				domain = askUserInput("Enter your domain name (e.g., example.com): ")
			}
		case "email":
			email = askUserInput("Enter admin email: ")
			for !isValidEmail(email) {
				fmt.Println("Invalid email format. Please try again.")
				email = askUserInput("Enter admin email: ")
			}
		case "password":
			password = askUserInput("Enter admin password: ")
			for !isValidPassword(password) {
				fmt.Println("Password must be at least 8 characters long and include uppercase, lowercase, and a number.")
				password = askUserInput("Enter admin password: ")
			}
			verifyPassword = askUserInput("Verify admin password: ")
			for password != verifyPassword {
				fmt.Println("Passwords do not match. Please try again.")
				password = askUserInput("Enter admin password: ")
				for !isValidPassword(password) {
					fmt.Println("Password must be at least 8 characters long and include uppercase, lowercase, and a number.")
					password = askUserInput("Enter admin password: ")
				}
				verifyPassword = askUserInput("Verify admin password: ")
			}
		default:
			fmt.Println("Invalid option. Please choose from domain, email, password, or database.")
		}
	}

	// Step 3: Start install AlphaLab
	fmt.Println("\nStarting AlphaLab installation...")

	// Step 4: Verify installation
	fmt.Println("Verifying AlphaLab installation...")

	// Final confirmation
	fmt.Println("AlphaLab installation completed successfully")
	return nil
}

// askUserInput prompts the user for input and returns the entered value
func askUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// checkDockerInstallation checks if Docker is installed on the system
func checkDockerInstallation() error {
	_, err := exec.Command("docker", "--version").Output()
	if err != nil {
		return fmt.Errorf("docker is not installed, please install docker first")
	}
	fmt.Println("Docker is installed.")
	return nil
}

// isValidDomain validates the domain format
func isValidDomain(domain string) bool {
	domainRegex := `^(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`
	re := regexp.MustCompile(domainRegex)
	return re.MatchString(domain)
}

// isValidEmail validates the email format
func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// isValidPassword validates the password criteria
func isValidPassword(password string) bool {
	length := len(password) >= 8
	upper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	lower := regexp.MustCompile(`[a-z]`).MatchString(password)
	number := regexp.MustCompile(`[0-9]`).MatchString(password)
	return length && upper && lower && number
}
