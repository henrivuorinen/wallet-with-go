package services

import (
	"fmt"
	"github.com/google/uuid"
)

// GenerateUUID creates a random UUID
func GenerateUUID() string {
	// Generate a new random UUID
	newUUID := uuid.New()

	// Convert the UUID to its string representation
	return newUUID.String()
}

func main() {
	// Example usage
	randomUUID := GenerateUUID()
	fmt.Println("Generated UUID:", randomUUID)
}
