package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/elliotxx/osp/pkg/log"
)

// AskForConfirmation asks the user for confirmation in command line interface.
// A user must type in "yes" or "no" and then press enter. It has fuzzy matching,
// so "y", "Y", "yes", "YES", and "Yes" all count as confirmations.
// If the input is not recognized, it will ask again.
// The function does not return until it gets a valid response from the user.
// Empty input (just pressing enter) is treated as "no".
func AskForConfirmation(message string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		log.P("?").C(log.ColorBlue).N().Log("%s [y/n]: ", message)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("failed to read user input: %w", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		// Handle empty input (just pressing enter)
		if response == "" {
			return false, nil
		}

		// Check for positive responses
		if response == "y" || response == "yes" {
			return true, nil
		}

		// Check for negative responses
		if response == "n" || response == "no" {
			return false, nil
		}
	}
}
