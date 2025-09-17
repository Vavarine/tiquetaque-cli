package app

import (
	"context"
	"fmt"
	"time"

	client "github.com/vavarine/ttq/internal/httpclient"
)

// Aqui você pode chamar SaveToken/SaveEmployeeID (keyring/fallback)
func DoLogin(ctx context.Context, cli *client.Client, email, code string) error {
	// poderia usar contexto com timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	token, employeeID, err := cli.Login(email, code)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	if token != "" {
		if err := SaveToken(token); err != nil {
			return fmt.Errorf("error saving token: %w", err)
		}

		fullName, err := cli.GetFullName(token)

		if err != nil {
			return fmt.Errorf("error getting employee info: %w", err)
		}

		if fullName != "" {
			if err := SaveEmployeeName(fullName); err != nil {
				return fmt.Errorf("error saving employeeName: %w", err)
			}

			fmt.Println("Bem-vindo,", fullName)
		}
	}

	if employeeID != "" {
		if err := SaveEmployeeID(employeeID); err != nil {
			return fmt.Errorf("error saving employeeID: %w", err)
		}

	}

	return nil
}
