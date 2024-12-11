package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

// checkDockerInstalled verifies if Docker and Docker Compose are installed.
func CheckDockerInstalled() bool {
	fmt.Println("Checking if Docker is installed...")

	dvSucess, dockerVersion := ExecCommand(false, "docker", "-v")
	dcvSucess, dockerComposeVersion := ExecCommand(false, "docker", "compose", "version")

	if !dvSucess || !dcvSucess {
		InstallDocker()
		return false
	}

	fmt.Printf("Docker and Docker Compose are installed!\n")
	fmt.Printf("%s", dockerVersion)
	fmt.Printf("%s\n", dockerComposeVersion)
	return true
}

// clearAllDocker stops and removes all running containers.
func ClearAllDocker() {
	containers := []string{"elimt-eln-pocketbase-1", "elimt-eln-backend-1", "elimt-eln-frontend-1"}

	// Check and kill only running containers
	for _, container := range containers {
		isRunning := isContainerRunning(container)
		if isRunning {
			ExecCommand(true, "docker", "container", "kill", container)
		}
	}

	// Prune all stopped containers
	ExecCommand(true, "docker", "container", "prune", "-f")

	// Prune all unused images
	ExecCommand(true, "docker", "image", "prune", "-a", "-f")

	// Prune unused build cache with all build history
	ExecCommand(true, "docker", "builder", "prune", "--all", "-f")
}
func CheckElimtRunning() bool {
	success, feedback := ExecCommand(false, "docker", "container", "ls")
	if !success {
		fmt.Println("Failed to execute docker container ls")
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
		fmt.Println("ELIMT containers already running")
		fmt.Println("This is not expected...")
		return true
	}

	return false
}

// runDockerCompose starts Docker Compose.
func RunDockerCompose() error {
	fmt.Println("Starting Docker Compose...")

	// Execute the docker-compose command
	success, output := ExecCommand(true, "docker", "compose", "up", "--detach")
	if !success {
		return fmt.Errorf("failed to start Docker Compose: %s", output)
	}

	fmt.Println("Docker Compose started successfully.")
	return nil
}

func InstallDocker() {
	osType := runtime.GOOS

	switch osType {
	case "linux":
		fmt.Println("Detected Linux. Proceeding with Docker installation for Linux...")

		// Run docker-install.sh
		if success, output := ExecCommand(false, "sudo ./docker-install.sh"); !success {
			log.Fatalf("Docker installation failed: %s", output)
		}

	case "darwin":
		fmt.Println("Detected macOS. Please install Docker Desktop manually from https://www.docker.com/products/docker-desktop.")
	case "windows":
		fmt.Println("Detected Windows. Please install Docker Desktop manually from https://www.docker.com/products/docker-desktop.")
	default:
		fmt.Printf("Unsupported operating system: %s\n", osType)
		log.Fatal("Docker installation is not supported on this OS.")
	}
}

// IsContainerRunning checks if a specific container is running
func isContainerRunning(container string) bool {
	success, output := ExecCommand(false, "docker", "inspect", "-f", "{{.State.Running}}", container)
	if !success {
		// If inspection fails, assume the container is not running
		return false
	}
	return strings.TrimSpace(output) == "true"
}
