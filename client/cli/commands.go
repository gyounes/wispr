package cli

import "fmt"

// ExecuteCommand handles CLI commands
func ExecuteCommand(cmd string) bool {
	switch cmd {
	case "/quit":
		fmt.Println("Goodbye!")
		return true
	case "/list":
		fmt.Println("Connected users: Alice, Bob") // placeholder for now
		return false
	default:
		fmt.Println("Unknown command")
		return false
	}
}
