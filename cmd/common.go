package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// GetToken parses the token from either the token flag or environment variable
func GetToken(cmd *cobra.Command) (string, error) {
	token, err := cmd.Flags().GetString("token")
	if err != nil {
		return "", err
	}

	if token != "" {
		return token, nil
	}

	token = os.Getenv("GIMME_GITHUB_TOKEN")
	if token != "" {
		return token, nil
	}

	return "", fmt.Errorf("cmd.GetToken: No token found in either --token flag or GIMME_GITHUB_TOKEN environment variable")
}
