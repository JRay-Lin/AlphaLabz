package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

// checkDockerInstalled verifies if Docker and Docker Compose are installed.
func CheckDockerInstalled() bool {
	fmt.Println("Checking if Docker is installed...")

	dvSucess, dockerVersion := ExecCommand("docker", "--version")
	dcvSucess, dockerComposeVersion := ExecCommand("docker-compose", "--version")

	if !dvSucess || !dcvSucess {
		fmt.Println("Docker or Docker Compose is not installed. Please install Docker and Docker Compose.")
		return false
	}

	fmt.Printf("Docker and Docker Compose are installed!\n")
	fmt.Printf("Docker version: %s", dockerVersion)
	fmt.Printf("Docker Compose version: %s\n", dockerComposeVersion)
	return true
}

// clearAllDocker stops and removes all running containers.
func ClearAllDocker() {
	fmt.Println("Stopping and removing all Docker containers...")

	// Stop all containers
	cmd := exec.Command("docker", "stop", "$(docker ps -q)")
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to stop containers: %v", err)
	}

	// Remove all containers
	cmd = exec.Command("docker", "rm", "-f", "$(docker ps -aq)")
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to remove containers: %v", err)
	} else {
		fmt.Println("All containers have been removed.")
	}
}

func CheckElimtRunning() bool {
	success, feedback := ExecCommand("docker", "ps")
	if !success {
		fmt.Println("Failed to execute docker ps")
		return false
	}

	// Required prefix for container names
	requiredPrefix := "elimt-eln"

	// Count the number of matching containers
	count := 0
	for _, line := range strings.Split(feedback, "\n") {
		if strings.Contains(line, requiredPrefix) {
			count++
		}
	}

	if count == 3 {
		return true
	}

	if count > 0 {
		fmt.Println("There are some container running but not all require container is running. This is not expected.")
		return true
	}

	return false
}

// runDockerCompose starts Docker Compose.
func RunDockerCompose() {
	fmt.Println("Starting Docker Compose...")

	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to start Docker Compose: %v", err)
	}

	fmt.Println("Docker Compose started successfully.")
}

func InstallDocker() {
	osType := runtime.GOOS

	switch osType {
	case "linux":
		fmt.Println("Detected Linux. Proceeding with Docker installation for Linux...")
		installDockerLinux()
	case "darwin":
		fmt.Println("Detected macOS. Please install Docker Desktop manually from https://www.docker.com/products/docker-desktop.")
	case "windows":
		fmt.Println("Detected Windows. Please install Docker Desktop manually from https://www.docker.com/products/docker-desktop.")
	default:
		fmt.Printf("Unsupported operating system: %s\n", osType)
		log.Fatal("Docker installation is not supported on this OS.")
	}
}

func installDockerLinux() {
	// Update the package index
	fmt.Println("Updating package index...")
	cmd := exec.Command("sudo", "apt-get", "update")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to update package index: %v", err)
	}

	// Download the Docker installation script
	fmt.Println("Downloading Docker installation script...")
	cmd = exec.Command("curl", "-fsSL", "https://get.docker.com", "-o", "get-docker.sh")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to download Docker installation script: %v", err)
	}

	// Run the installation script
	fmt.Println("Running Docker installation script...")
	cmd = exec.Command("sudo", "sh", "get-docker.sh")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to install Docker: %v", err)
	}

	// Install Docker Compose
	fmt.Println("Installing Docker Compose...")
	composeURL := "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)"
	cmd = exec.Command("sudo", "curl", "-L", composeURL, "-o", "/usr/local/bin/docker-compose")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to download Docker Compose: %v", err)
	}

	cmd = exec.Command("sudo", "chmod", "+x", "/usr/local/bin/docker-compose")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to make Docker Compose executable: %v", err)
	}

	// Verify Installation
	sucess := CheckDockerInstalled()
	if !sucess {
		log.Fatalf("Failed to verify Docker installation.")
	}

	fmt.Println("Docker and Docker Compose installed successfully!")
}
