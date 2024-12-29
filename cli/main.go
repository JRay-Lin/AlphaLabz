package main

import (
	"alphalab-cli/pkg/installation"
	"fmt"

	"github.com/alecthomas/kong"
)

const cliVersion string = "0.0.1"

// GreetCmd represents the 'greet' command
// This is a command for testing, it will be remove after develpment is finished
type GreetCmd struct {
	Name string `arg:"" required:"" help:"Name of the person to greet."`
}

func (g *GreetCmd) Run() error {
	fmt.Printf("Hello, %s!\n", g.Name)
	return nil
}

// VersionCmd represents the 'version' command
type VersionCmd struct{}

// getAppVersion use API to retrive the version of app
func getAppVersion() (string, error) {
	// for testing
	return "v0.0.1", nil
}

func (v *VersionCmd) Run() error {
	// Print the CLI version
	fmt.Println("CLI version:", cliVersion)

	// Retrieve and print the application version
	appVersion, err := getAppVersion()
	if err != nil {
		return fmt.Errorf("failed to get app version: %w", err)
	}
	fmt.Println("Application version:", appVersion)
	return nil
}

// CLI struct defines the commands and flags for CLI
type CLI struct {
	Greet   GreetCmd                `cmd:"" help:"Print a greeting."`
	Version VersionCmd              `cmd:"" short:"v" help:"Print CLI & App version information."`
	Install installation.InstallCmd `cmd:"" help:"Install the app."`
}

func main() {
	// Parse and run the CLI commands
	var cli CLI
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
