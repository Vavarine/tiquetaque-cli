package app

import (
	"context"
	"fmt"

	client "github.com/vavarine/ttq/internal/httpclient"
)

func GetPunchHistory(ctx context.Context, cli *client.Client) (*client.DayRecordsResponse, error) {
	// Loads the saved token, employeeID, and fullName from secure storage
	token, err := LoadToken()
	if err != nil || token == "" {
		return nil, fmt.Errorf("token não encontrado. Faça login primeiro")
	}

	employeeID, err := LoadEmployeeID()
	if err != nil || employeeID == "" {
		return nil, fmt.Errorf("employeeID não encontrado. Refazer login")
	}

	fullName, err := LoadEmployeeName()
	if err != nil || fullName == "" {
		return nil, fmt.Errorf("employeeName não encontrado. Refazer login")
	}

	resp, err := cli.GetDayRecords(ctx, token, employeeID)

	if err != nil {
		return nil, fmt.Errorf("error when fetching punch history: %w", err)
	}

	return resp, nil
}
