package casbin

import (
	"bufio"
	"fmt"
	"os"
)

// SavePoliciesToFile writes Casbin policies to a file (for use in CSV format)
func SavePoliciesToFile(policies []string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, policy := range policies {
		_, err := writer.WriteString(policy + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	writer.Flush()
	fmt.Println("Casbin policies saved to", filename)
}
