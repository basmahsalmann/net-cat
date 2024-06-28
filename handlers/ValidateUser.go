package handlers

import (
	"fmt"
)

func IsUsernameTaken(username string) bool {
    ClientsMutex.Lock()
    defer ClientsMutex.Unlock()
    for _, existingName := range Clients {
        if existingName == username {
            return true
        }
    }
    return false
}

func ValidateUsername(name string) error {
    if len(name) > MaxLength {
        return fmt.Errorf("Name cannot exceed %d characters", MaxLength)
    }
    for _, char := range name {
        if (char < 'A' || char > 'Z') && (char < 'a' || char > 'z') {
            return fmt.Errorf("Name cannot contain special characters or numbers")
        }
    }
    return nil
}

func ValidateName(name string) error {
	if len(name) > MaxLength {
		return fmt.Errorf("Name cannot exceed %d characters", MaxLength)
	}
	for _, char := range name {
		if (char < 'A' || char > 'Z') && (char < 'a' || char > 'z') {
			return fmt.Errorf("Name cannot contain special characters or numbers")
		}
	}
	return nil
}
